# Eventline terraform provider

[Read the full documentation here](https://registry.terraform.io/providers/adyxax/eventline/latest/docs)

The [Eventline](https://www.exograd.com/products/eventline/) provider is used to interact with the resources supported by Eventline. The provider needs to be configured with the proper credentials before it can be used. It requires terraform 1.0 or later.

## Example Usage

```terraform
terraform {
  required_providers {
    eventline = {
      source = "adyxax/eventline"
    }
  }
}

provider "eventline" {
  api_key  = var.eventline_api_key
  endpoint = "http://localhost:8085/"
}

data "eventline_project" "main" {
  name = "main"
}

resource "eventline_identity" "example" {
  name       = "example"
  project_id = data.eventline_project.main.id

  connector = "eventline"
  data      = jsonencode({ "key" = "test" })
  type      = "api_key"
}
```

## Developing the provider

TODO
