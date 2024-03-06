package docs

func GenerateImportDocs(idFormat string, resourceName string) string {
	return `## Import

In Terraform v1.5.0 and later, use an [import block](https://developer.hashicorp.com/terraform/language/import) with ` + "`" + idFormat + "`" + `. For example:

` + "```" + `terraform
import {
  to = ` + resourceName + `.example
  id = "` + idFormat + `"
}
` + "```" + `

Otherwise, use ` + "`terraform import`" + ` with ` + "`" + idFormat + "`" + `. For example:

` + "```" + `console
terraform import ` + resourceName + `.example ` + idFormat + `
` + "```"
}
