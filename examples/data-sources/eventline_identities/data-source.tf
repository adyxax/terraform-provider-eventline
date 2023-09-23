data "eventline_project" "main" {
  name = "main"
}

data "eventline_identities" "example" {
  project_id = data.project.main.id
}
