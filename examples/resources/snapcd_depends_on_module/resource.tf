resource "snapcd_depends_on_module" "example" {
  module_id            = snapcd_module.app.id
  depends_on_module_id = snapcd_module.database.id
}