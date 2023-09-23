data "eventline_project" "main" {
  name = "main"
}

resource "eventline_identity" "example" {
  name       = "example"
  project_id = data.project.main.id

  connector = "eventline"
  data      = "{\n    \"key\": \"test\"\n  }"
  type      = "api_key"
}
