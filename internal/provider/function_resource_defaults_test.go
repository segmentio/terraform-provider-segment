package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const functionDefaultLogo = "https://cdn-devcenter.segment.com/2e87e186-3bca-4d55-b93b-97705deb2a73.svg"

// functionFakeAPI is a stateful stand-in for the Functions API. It mirrors the
// backend defaults (logo assigned when absent, settings always returned) and,
// like public-api's `IsString().uri().optional()`, rejects a logoUrl that is
// present but empty.
type functionFakeAPI struct {
	mu     sync.Mutex
	fns    map[string]map[string]any
	posts  []map[string]any
	nextID int
}

func newFunctionFakeAPI() (*functionFakeAPI, *httptest.Server) {
	f := &functionFakeAPI{fns: map[string]map[string]any{}}

	return f, httptest.NewServer(http.HandlerFunc(f.serve))
}

func (f *functionFakeAPI) serve(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")
	f.mu.Lock()
	defer f.mu.Unlock()

	reply := func(status int, fn map[string]any) {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"function": fn}})
	}

	if req.URL.Path == "/functions" && req.Method == http.MethodPost {
		var in map[string]any
		raw, _ := io.ReadAll(req.Body)
		_ = json.Unmarshal(raw, &in)
		f.posts = append(f.posts, in)

		if v, present := in["logoUrl"]; present && v == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"errors":[{"type":"validation","message":"logoUrl must be a valid uri"}]}`))

			return
		}

		f.nextID++
		id := "fn-" + string(rune('0'+f.nextID))
		fn := map[string]any{
			"id": id, "workspaceId": "ws", "code": in["code"], "displayName": in["displayName"],
			"resourceType": in["resourceType"], "catalogId": "cat-" + id, "previewWebhookUrl": "",
			"logoUrl": functionDefaultLogo, "settings": []any{},
		}
		if v, ok := in["logoUrl"]; ok {
			fn["logoUrl"] = v
		}
		if v, ok := in["settings"]; ok && v != nil {
			fn["settings"] = v
		}
		if v, ok := in["description"]; ok {
			fn["description"] = v
		}
		f.fns[id] = fn
		reply(http.StatusOK, fn)

		return
	}

	id := strings.TrimPrefix(req.URL.Path, "/functions/")
	fn, ok := f.fns[id]
	switch {
	case !ok:
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors":[{"type":"not-found"}]}`))
	case req.Method == http.MethodDelete:
		delete(f.fns, id)
		_, _ = w.Write([]byte(`{"data":{"status":"SUCCESS"}}`))
	case req.Method == http.MethodPatch:
		var in map[string]any
		raw, _ := io.ReadAll(req.Body)
		_ = json.Unmarshal(raw, &in)
		for k, v := range in {
			if v != nil && v != "" {
				fn[k] = v
			}
		}
		reply(http.StatusOK, fn)
	default:
		reply(http.StatusOK, fn)
	}
}

func functionProviderConfig(url string) string {
	return `
		provider "segment" {
			url   = "` + url + `"
			token = "abc123"
		}
	`
}

func functionResourceConfig(extra string) string {
	return `
		resource "segment_function" "test" {
			code          = "// code"
			display_name  = "Defaults function"
			resource_type = "SOURCE"
			` + extra + `
		}
	`
}

// TestAccFunctionResource_OnlyOneDefaultedAttributeSet sets logo_url and settings
// one at a time. Setting only one must not be enough to trip the other.
func TestAccFunctionResource_OnlyOneDefaultedAttributeSet(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		extra      string
		wantLogo   string
		wantNSet   string
		wantLogoOn bool // logoUrl expected on the wire
	}{
		"only logo_url": {
			extra:      `logo_url = "https://segment.com/mine.png"`,
			wantLogo:   "https://segment.com/mine.png",
			wantNSet:   "0",
			wantLogoOn: true,
		},
		"only settings": {
			extra: `settings = [{
				name = "k", label = "K", description = "d", type = "STRING", required = true, sensitive = false
			}]`,
			wantLogo:   functionDefaultLogo,
			wantNSet:   "1",
			wantLogoOn: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			api, srv := newFunctionFakeAPI()
			defer srv.Close()

			cfg := functionProviderConfig(srv.URL) + functionResourceConfig(tc.extra)
			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: cfg,
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("segment_function.test", "logo_url", tc.wantLogo),
							resource.TestCheckResourceAttr("segment_function.test", "settings.#", tc.wantNSet),
						),
					},
					{Config: cfg, PlanOnly: true},
				},
			})

			if len(api.posts) == 0 {
				t.Fatal("no create request reached the API")
			}
			if _, present := api.posts[0]["logoUrl"]; present != tc.wantLogoOn {
				t.Fatalf("logoUrl on the wire = %v, want %v (body %v)", present, tc.wantLogoOn, api.posts[0])
			}
		})
	}
}

// TestAccFunctionResource_SettingsOmittedEmptyOrNull checks that omitting settings,
// setting it to `[]` and setting it to null are interchangeable.
func TestAccFunctionResource_SettingsOmittedEmptyOrNull(t *testing.T) {
	t.Parallel()

	_, srv := newFunctionFakeAPI()
	defer srv.Close()

	omitted := functionProviderConfig(srv.URL) + functionResourceConfig(``)
	empty := functionProviderConfig(srv.URL) + functionResourceConfig(`settings = []`)
	null := functionProviderConfig(srv.URL) + functionResourceConfig(`settings = null`)

	zero := resource.TestCheckResourceAttr("segment_function.test", "settings.#", "0")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: omitted, Check: zero},
			{Config: omitted, PlanOnly: true},
			// omitted -> [] must be a no-op: both mean "no settings".
			{Config: empty, PlanOnly: true},
			{Config: empty, Check: zero},
			// explicit null behaves like omitted for a computed attribute.
			{Config: null, PlanOnly: true},
			{Config: null, Check: zero},
		},
	})
}

// TestAccFunctionResource_ImportWithAPIDefaults imports a function created with
// API-populated defaults and expects no drift afterwards.
func TestAccFunctionResource_ImportWithAPIDefaults(t *testing.T) {
	t.Parallel()

	_, srv := newFunctionFakeAPI()
	defer srv.Close()

	cfg := functionProviderConfig(srv.URL) + functionResourceConfig(``)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "segment_function.test",
				Config:            cfg,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// No drift after the import round-trip.
			{Config: cfg, PlanOnly: true},
		},
	})
}
