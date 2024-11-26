resource "cronitor_http_monitor" "this" {
  name     = "Some monitor"
  schedule = "every 5 minutes"
  url      = "https://registry.terraform.io/providers/henrywhitaker3/cronitor/latest"
  method   = "GET"
  assertions = [
    "response.code = 200"
  ]
}

# Create a notification list and a monitor that uses it
resource "cronitor_notification_list" "this" {
  name = "Demo"
  webhooks = [
    "https://registry.terraform.io/providers/henrywhitaker3/cronitor/latest"
  ]
}

resource "cronitor_http_monitor" "this" {
  name     = "Some monitor"
  schedule = "every 5 minutes"
  url      = "https://registry.terraform.io/providers/henrywhitaker3/cronitor/latest"
  method   = "GET"
  assertions = [
    "response.code = 200"
  ]
  notify = [
    cronitor_notification_list.this.key,
  ]
}
