data "eventline_project" "main" {
  name = "main"
}

data "eventline_identities" "example" {
  project_id = data.eventline_project.main.id
}
