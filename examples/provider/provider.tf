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
