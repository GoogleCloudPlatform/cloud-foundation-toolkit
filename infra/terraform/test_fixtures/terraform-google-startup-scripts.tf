resource "random_id" "startup_scripts_github_webhook_token" {
  byte_length = 20
}

data "template_file" "startup_scripts_github_webhook_url" {

  template = "https://concourse.infra.cft.tips/api/v1/teams/cft/pipelines/$${pipeline}/resources/pull-request/check/webhook?webhook_token=$${webhook_token}"

  vars {
    pipeline = "terraform-google-startup-scripts"
    webhook_token = "${random_id.startup_scripts_github_webhook_token.hex}"
  }
}

resource "kubernetes_secret" "ci_startup_scripts" {
  metadata {
    namespace = "concourse-cft"
    name = "startup-scripts"
  }
  data {
    github_webhook_token = "${random_id.startup_scripts_github_webhook_token.hex}"
  }
}
