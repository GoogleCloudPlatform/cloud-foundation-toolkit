resource "google_compute_global_address" "concourse" {
  name = "concourse-${terraform.workspace}"
}

resource "google_dns_record_set" "concourse" {
  name = "${var.concourse_subdomain[terraform.workspace]}.${data.google_dns_managed_zone.tips_cft_infra.dns_name}"
  type = "A"
  ttl = 300
  managed_zone = "${data.google_dns_managed_zone.tips_cft_infra.name}"
  rrdatas = ["${google_compute_global_address.concourse.address}"]
}

# Service Account for main project.

resource "google_service_account" "concourse_cft_team" {
  account_id   = "concourse-cft-team"
  display_name = "concourse-cft-team"
}

resource "google_project_iam_binding" "concourse_cft_team_storage" {
  role = "roles/storage.admin"
  members = [
    "serviceAccount:${google_service_account.concourse_cft_team.email}",
  ]
}

resource "google_service_account_key" "concourse_cft_team" {
  service_account_id = "${google_service_account.concourse_cft_team.id}"
}

# Service Account in phoogle.net organization to run integration tests with.

resource "kubernetes_secret" "concourse_cft_team_service_accounts" {
  metadata {
    namespace = "concourse-cft"
    name      = "sa"
  }
  data {
    google  = "${base64decode(google_service_account_key.concourse_cft_team.private_key)}"
  }
  depends_on = ["helm_release.concourse"]
}

locals {
  concourse_host = "${substr(google_dns_record_set.concourse.name, 0, length(google_dns_record_set.concourse.name) - 1)}"
}

resource "helm_release" "concourse" {
  depends_on = ["null_resource.helm_init"]

  name      = "concourse"
  chart     = "stable/concourse"
  version   = "3.0.0"
  namespace = "default"
  keyring   = ""

  values = [
    <<EOF
imageTag: 4.2.1
concourse:
  web:
    externalUrl: https://${local.concourse_host}
    auth:
      oidc:
        enabled: true
        displayName: Google CFT CICD
        issuer: https://accounts.google.com
      mainTeam:
        localUser: concourse
    postgres:
      host: ${data.terraform_remote_state.postgres.ip_address}
    kubernetes:
      teams:
      - cft
  worker:
    baggageclaim:
      driver: overlay
web:
  service:
    type: NodePort
  ingress:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: gce
      kubernetes.io/ingress.global-static-ip-name: ${google_compute_global_address.concourse.name}
    hosts:
    - ${local.concourse_host}
    tls:
    - secretName: concourse-web-tls
      hosts:
      - ${local.concourse_host}
secrets:
  create: false
postgresql:
  enabled: false
EOF
  ]
}
