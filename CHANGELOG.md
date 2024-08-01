## 1.0.2 (Aug 1, 2024)
Fixes a bug where imported `segment_reverse_etl_model` resources could have invalid configuration parameters.

## 1.0.1 (April 18, 2024)
Fixes a bug where the `segment_profiles_warehouse` resource was dereferencing a nil pointer upon invalid import.

## 1.0.0 (April 8, 2024)
General Availability release. No changes from the previous release.

## 0.10.3 (March 6, 2024)
Correctly taints resources upon failed create to ensure resources are not duplicated.

## 0.10.2 (March 6, 2024)
Fixes a rendering bug for `segment_profiles_warehouse` resource docs.

## 0.10.1 (March 6, 2024)
Adds import documentation for each resource. Also clean up documentation.

## 0.10.0 (January 30, 2024)
Fixes bug where segment_function IDs for insert destinations could not be passed directly into a segment_insert_function_instance resource.

## 0.9.0 (January 8, 2024)
Adds support for Transformation resource.

## 0.8.0 (November 16, 2023)
**BREAKING CHANGE**

Move Source schema settings configuration to `segment_source_tracking_plan_connection` resource. Also adds schema settings to Source data source.

## 0.7.0 (November 1, 2023)
Support Insert Function instance resource.

## 0.6.0 (October 31, 2023)
Adds support for Source schema settings in the Source resource. Also removes the description field from Labels fields in various resources since it is not populated and not likely to be used for configuration.

## 0.5.1 (October 27, 2023)
Fixes an issue where Function setting type could not be set to `TEXT_MAP``.

## 0.5.0 (October 26, 2023)
Adds support for Reverse ETL model resource.
(Previously misreleased as 0.0.5)

## 0.4.0 (October 19, 2023)
Adds support for Destination subscriptions and Source - Tracking Plan connection resources.

## 0.3.0 (October 16, 2023)
Adds support for Profiles Warehouse and fixes nil pointer access issues with Destination metadata.

## 0.2.0 (October 13, 2023)
Adds support for destination filters and functions.

## 0.1.0 (October 6, 2023)
Adds support for IAM resource (users, user groups)

## 0.0.2 (September 22, 2023)
Initial release
