data "project" "main" {
  name = "main"
}

data "eventline_jobs" "example" {
  project_id = data.project.main.id
}
