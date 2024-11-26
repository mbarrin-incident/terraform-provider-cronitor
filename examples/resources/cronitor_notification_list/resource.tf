resource "cronitor_notification_list" "this" {
  name = "Demo"
  webhooks = [
    # The name of a webhooks integration setup in the cronitor dashboard
    "name"
  ]
}
