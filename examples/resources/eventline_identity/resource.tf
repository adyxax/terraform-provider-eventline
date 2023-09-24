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
