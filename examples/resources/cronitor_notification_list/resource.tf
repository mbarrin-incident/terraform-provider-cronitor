resource "cronitor_notification_list" "this" {
  name = "Demo"
  webhooks = [
    "https://registry.terraform.io/providers/henrywhitaker3/cronitor/latest"
  ]
}
