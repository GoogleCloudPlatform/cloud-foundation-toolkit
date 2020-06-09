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

1. GKE Cluster with Config Connector
1. [Helm](../../../README.md#helm)

## USAGE

All steps are run from this directory.

1. Review and update the values in `./lb/values.yaml`. Note that if you change cluster name and locaiton, you will need to change how they are used in `gcloud container clusters get-credentials` commands below. [PROJECT_ID] refers to the project where all the GCP resources will be created.  

1. Configure clusters and load balancing resources with Helm

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
    # The command uses cluster names based on the values passed in the previous step
    kubectl wait --for=condition=Ready containercluster/cluster-na containercluster/cluster-eu
    ```

1. Configure first cluster

    ```bash
    # switch the context to the first cluster. The command uses cluster name and zone based on the values used to create the clusters.
    gcloud container clusters get-credentials cluster-na --zone=us-central1-a

    # validate your chart
    helm lint ./cluster/ --set projectId=[PROJECT_ID] --set localMessage="Hello from North America!"

    # install your chart
    helm install ./cluster/ --set projectId=[PROJECT_ID] --set localMessage="Hello from North America!" --generate-name

     # annotate service account
    kubectl annotate sa -n autoneg-system default iam.gke.io/gcp-service-account=autoneg-system@[PROJECT_ID].iam.gserviceaccount.com

    # ensure pods are ready
    kubectl wait --for=condition=Ready pods --all
    ```

1. Repeat the step for the other cluster

    ```bash
    # switch the context to the second cluster. The command uses cluster name and zone based on the values used to create the clusters.
    gcloud container clusters get-credentials cluster-eu --zone=europe-west2-a

    # validate your chart
    helm lint ./cluster/ --set projectId=[PROJECT_ID] --set localMessage="Hello from Europe"


    # install your chart
    helm install ./cluster/ --set projectId=[PROJECT_ID] --set localMessage="Hello from Europe" --generate-name

     # annotate service account
    kubectl annotate sa -n autoneg-system default iam.gke.io/gcp-service-account=autoneg-system@[PROJECT_ID].iam.gserviceaccount.com

    # ensure pods are ready
    kubectl wait --for=condition=Ready pods --all
    ```

1. Switch the context to the main cluster and run verify that multi-cluster ingress is configured

    ```bash
    # switch the context to the main cluster
    gcloud container clusters get-credentials [CLUSTER NAME] --zone=[CLUSTER ZONE]

    # curl the external address of the forwarding rule. Note that it might take around 5-10 minutes for load balancing to start working.
    # You will see the message ("Hello from North America" or "Hello from Europe" based on your location).
    curl $(kubectl get  computeforwardingrule -o=jsonpath='{.items[0].spec.ipAddress.addressRef.external}')

1. Clean up the installation:

    ```bash
    # switch the context to the first cluster. The command uses cluster name and zone based on the values used to create the clusters.
    gcloud container clusters get-credentials cluster-na --zone=us-central1-a

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

     # switch the context to the main cluster
    gcloud container clusters get-credentials [CLUSTER NAME] --zone=[CLUSTER ZONE]

    # list Helm releases to obtain release name
    helm list

    # delete release specifying release name from the previous command output. Note that can take a few minutes before all K8s resources are fully deleted.
    helm delete [release_name]
    ```

## LICENSE

Apache 2.0 - See [LICENSE](/LICENSE) for more information.