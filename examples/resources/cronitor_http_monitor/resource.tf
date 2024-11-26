resource "cronitor_http_monitor" "this" {
  name     = "Some monitor"
  schedule = "every 5 minutes"
  url      = "https://registry.terraform.io/providers/henrywhitaker3/cronitor/latest"
  method   = "GET"
  assertions = [
    "response.code = 200"
  ]
}
