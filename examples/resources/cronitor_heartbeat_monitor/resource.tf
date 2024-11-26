resource "cronitor_heartbeat_monitor" "this" {
  name     = "Some monitor"
  schedule = "* * * * *"
}
