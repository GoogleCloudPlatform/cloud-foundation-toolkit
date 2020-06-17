# Autoneg

==================================================

## NAME

  autoneg

## SYNOPSIS

Autoneg solution uses network endpoint groups to provision load balancing across multiple clusters.
This solution uses `./cluster/templates/autoneg.yaml` from [gke-autoneg-controller](https://github.com/GoogleCloudPlatform/gke-autoneg-controller).
For demonstration purposes it uses a docker image with a simple Node.js service (bulankou/node-hello-world) with a single endpoint that prints out a message.

## CONSUMPTION

  1. Clone GoogleCloudPlatform/cloud-foundation-toolkit repository
  
      ```bash
      git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit.git
      ```

  1. Go to the autoneg folder:

      ```bash
      cd cloud-foundation-toolkit/config-connector/solutions/networking/helm/autoneg
      ```

## REQUIREMENTS

1. [Helm](../../../README.md#helm)
1. GKE Cluster with Config Connector. This solution assumes that all resources are installed in the same project, where the cluster with Config Connector is installed, and that load balancing resources are installed on the same cluster where Config Connector is installed. If you would like to configure your resources in a different project, the easiest approach would be to give your Config Connector service account (`cnrm-system`) owner permissions on this target project.
1. If your Config Connector version is earlier than [1.12.0](https://github.com/GoogleCloudPlatform/k8s-config-connector/releases) you need to apply [this workaround](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/78#issuecomment-577285402) to `iampolicymembers.iam.cnrm.cloud.google.com` CRD.
1. `compute.googleapis.com`, `container.googleapis.com` and `cloudresourcemanager.googleapis.com` APIs should be enabled on the project managed by Config Connector, in addition to the default services enabled.

## USAGE

All steps are run from this directory.

1. Create the clusters.

    Review and update the values in `./clusters/values.yaml`. Note that if you change cluster name and location, you will need to change how they are used in `gcloud container clusters get-credentials` commands below. [PROJECT_ID] refers to the project where all the GCP resources will be created.  

    ```bash
    # validate your chart
    helm lint ./clusters/ --set projectId=[PROJECT_ID]

    # check the output of your chart
    helm template ./clusters/ --set projectId=[PROJECT_ID]

    # install your chart
    helm install ./clusters/ --set projectId=[PROJECT_ID] --generate-name
    ```

1. Create load balancing resources.

    Review and update the values in `./lb/values.yaml`. [PROJECT_ID] refers to the project where all the GCP resources will be created.

    ```bash
    # validate your chart
    helm lint ./lb/ --set projectId=[PROJECT_ID]

    # check the output of your chart
    helm template ./lb/ --set projectId=[PROJECT_ID]

    # install your chart
    helm install ./lb/ --set projectId=[PROJECT_ID] --generate-name
    ```

1. Wait for clusters to be created

    ```bash
    # The command uses cluster names based on the values passed in the ealier step
    kubectl wait --for=condition=Ready containercluster/cluster-na containercluster/cluster-eu
    ```

1. Configure first cluster

    ```bash
    # switch the context to the first cluster. The command uses cluster name and zone based on the values used to create the clusters.
    gcloud container clusters get-credentials cluster-na --zone=us-central1-b

    # validate your chart
    helm lint ./workload/ --set projectId=[PROJECT_ID] --set localMessage="Hello from North America\!"

    # install your chart
    helm install ./workload/ --set projectId=[PROJECT_ID] --set localMessage="Hello from North America\!" --generate-name

     # annotate service account
    kubectl annotate sa -n autoneg-system default iam.gke.io/gcp-service-account=autoneg-system@[PROJECT_ID].iam.gserviceaccount.com

    # ensure pods are ready
    kubectl wait --for=condition=Ready pods --all

    # check the service and ensure that `anthos.cft.dev/autoneg-status` annotation is present in the output
    kubectl get svc node-app-backend -o=jsonpath='{.metadata.annotations}'
    ```

1. Repeat the step for the other cluster

    ```bash
    # switch the context to the second cluster. The command uses cluster name and zone based on the values used to create the clusters.
    gcloud container clusters get-credentials cluster-eu --zone=europe-west2-a

    # validate your chart
    helm lint ./workload/ --set projectId=[PROJECT_ID] --set localMessage="Hello from Europe\!"


    # install your chart
    helm install ./workload/ --set projectId=[PROJECT_ID] --set localMessage="Hello from Europe\!" --generate-name

     # annotate service account
    kubectl annotate sa -n autoneg-system default iam.gke.io/gcp-service-account=autoneg-system@[PROJECT_ID].iam.gserviceaccount.com

    # ensure pods are ready
    kubectl wait --for=condition=Ready pods --all

    # check the service and ensure that `anthos.cft.dev/autoneg-status` annotation is present in the output
    kubectl get svc node-app-backend -o=jsonpath='{.metadata.annotations}'
    ```

1. Switch the context to the cluster that contains the configs for load balancing resources and run verify that multi-cluster ingress is configured

    ```bash
    # switch the context to the main cluster
    gcloud container clusters get-credentials [CLUSTER NAME] --zone=[CLUSTER ZONE]

    # if you created the load balancing resources in the namespace, other than default, switch the context to that namespace
    kubectl config set-context --current --namespace [NAMESPACE]

    # verify that your backend service has 2 backends attached (select index of "global" if prompted)
    gcloud compute backend-services describe node-app-backend-service
    ```

    The backends section of the output should list both backends, for example:

    ```yaml
    backends:
    - balancingMode: RATE
      capacityScaler: 1.0
      group: https://www.googleapis.com/compute/v1/projects/<project_id>/zones/us-central1-b/networkEndpointGroups/k8s1-37f1db7d-default-node-app-backend-80-486adca6
      maxRatePerEndpoint: 100.0
    - balancingMode: RATE
      capacityScaler: 1.0
      group: https://www.googleapis.com/compute/v1/projects/<project_id>/zones/europe-west2-a/networkEndpointGroups/k8s1-292a63d7-default-node-app-backend-80-636c84c5
      maxRatePerEndpoint: 100.0
      connectionDraining:
      drainingTimeoutSec: 300
    ```

    Verify that load balancing resources are forwarding the request to the backend:

    ```bash
    # curl the external address of the forwarding rule. Note that it might take around 5-10 minutes for load balancing to start working.
    # You will see the message ("Hello from North America!" or "Hello from Europe!" based on your location).
    curl $(kubectl get computeforwardingrule -o=jsonpath='{.items[0].spec.ipAddress.addressRef.external}')

1. Clean up the installation:

    ```bash
    # switch the context to the first cluster. The command uses cluster name and zone based on the values used to create the clusters.
    gcloud container clusters get-credentials cluster-na --zone=us-central1-b

    # list Helm releases to obtain release name
    helm list

    # delete release specifying release name from the previous command output. Note that can take a few minutes before all K8s resources are fully deleted.
    helm delete [release_name]
    
     # switch the context to the second cluster. The command uses cluster name and zone based on the values used to create the clusters.
    gcloud container clusters get-credentials cluster-eu --zone=europe-west2-a

    # list Helm releases to obtain release name
    helm list

    # delete release specifying release name from the previous command output. Note that can take a few minutes before all K8s resources are fully deleted.
    helm delete [release_name]

     # switch the context to the cluster that contains the configs for load balancing resources
    gcloud container clusters get-credentials [CLUSTER NAME] --zone=[CLUSTER ZONE]

    # if you created the load balancing resources in the namespace, other than default, switch the context to that namespace
    kubectl config set-context --current --namespace [NAMESPACE]

    # list Helm releases to obtain release name
    helm list

    # delete release specifying release name from the previous command output. Note that can take a few minutes before all K8s resources are fully deleted.
    helm delete [release_name]
    ```

## LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.