resource "kubernetes_service_account" "tiller" {
  metadata {
    name = "tiller"
    namespace = "kube-system"
  }
}

resource "kubernetes_cluster_role_binding" "tiller" {
  depends_on = ["kubernetes_service_account.tiller"]
  metadata {
    name = "tiller"
  }
  role_ref {
    kind = "ClusterRole"
    name = "cluster-admin"
    api_group = "rbac.authorization.k8s.io"
  }
  subject {
    kind = "ServiceAccount"
    name = "tiller"
    namespace = "kube-system"
    api_group = ""
  }
}

resource "null_resource" "helm_init" {
  depends_on = ["kubernetes_cluster_role_binding.tiller", "kubernetes_service_account.tiller"]

  provisioner "local-exec" {
    command = <<EOM
gcloud container clusters get-credentials ${data.terraform_remote_state.gke.cluster_name} \
  --project ${module.variables.project_id} --region ${module.variables.region[terraform.workspace]} && \
helm init --service-account ${kubernetes_cluster_role_binding.tiller.metadata.0.name} && \
kubectl -n kube-system patch deployment tiller-deploy -p '{"spec": {"template": {"spec": {"automountServiceAccountToken": true}}}}'
EOM
  }
}
