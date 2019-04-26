resource "helm_release" "cert_manager" {

  depends_on = ["null_resource.helm_init"]

  name      = "cert-manager"
  chart     = "stable/cert-manager"
  namespace = "kube-system"
  keyring   = ""

  values = [
    <<EOF
image:
  tag: v0.5.2
EOF
  ]
}

resource "null_resource" "cert_manager_cluster_issuer" {

  depends_on = ["helm_release.concourse", "helm_release.cert_manager"]

  provisioner "local-exec" {
  command = <<EOM
cat <<EOF | kubectl apply -f -
apiVersion: certmanager.k8s.io/v1alpha1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: ${var.lets_encrypt_email[terraform.workspace]}
    privateKeySecretRef:
      name: letsencrypt-prod
    http01: {}
EOF
EOM
  }
}

resource "null_resource" "cert_manager_certificate" {

  depends_on = ["helm_release.concourse", "helm_release.cert_manager"]

  provisioner "local-exec" {
  command = <<EOM
cat <<EOF | kubectl apply -f -
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: ${replace(local.concourse_host, ".", "-")}
  namespace: default
spec:
  secretName: concourse-web-tls
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  commonName: ${local.concourse_host}
  dnsNames:
  - ${local.concourse_host}
  acme:
    config:
    - http01:
        ingress: concourse-web
      domains:
      - ${local.concourse_host}
EOF
EOM
  }
}
