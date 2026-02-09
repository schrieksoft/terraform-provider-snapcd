data "snapcd_runner" "myrunner" {
  name = "myrunner"
}

resource "snapcd_source_refresher_preselection" "example" {
  source_url = "https://github.com/schrieksoft/snapcd-samples.git"
  runner_id  = data.snapcd_runner.myrunner.id
}
