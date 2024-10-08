---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pgvecto-rs-cloud Provider"
subcategory: ""
description: |-
  You can use the this Terraform provider to manage resources supported
  by PGVecto.rs Cloud https://cloud.pgvecto.rs. You can refer to get started guide https://docs.pgvecto.rs/cloud/manage/terraform.html for more information.
---

# pgvecto-rs-cloud Provider

You can use the this Terraform provider to manage resources supported 
by [PGVecto.rs Cloud](https://cloud.pgvecto.rs). You can refer to [get started guide](https://docs.pgvecto.rs/cloud/manage/terraform.html) for more information.

## Example Usage

```terraform
terraform {
  required_providers {
    pgvecto-rs-cloud = {
      source = "tensorchord/pgvecto-rs-cloud"
    }
  }
}

provider "pgvecto-rs-cloud" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_key` (String, Sensitive) PGVecto.rs Cloud API Key. Can be configured by setting PGVECTORS_CLOUD_API_KEY environment variable.

### Optional

- `api_url` (String) The URL of the PGVecto.rs Cloud API.
