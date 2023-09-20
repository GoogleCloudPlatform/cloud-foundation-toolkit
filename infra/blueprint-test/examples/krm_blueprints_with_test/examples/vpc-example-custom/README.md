# Virtual Private Cloud blueprint

A Virtual Private Cloud (VPC)

## Setters

```
Setter        Usages
namespace     1
network-name  1
project-id    3
```

## Sub-packages

This package has no sub-packages.

## Resources

```
File           APIVersion                                  Kind            Name                Namespace
services.yaml  serviceusage.cnrm.cloud.google.com/v1beta1  Service         project-id-compute  projects
vpc.yaml       compute.cnrm.cloud.google.com/v1beta1       ComputeNetwork  network-name        networking
```

## Resource References

- [ComputeNetwork](https://cloud.google.com/config-connector/docs/reference/resource-docs/compute/computenetwork)
- [Service](https://cloud.google.com/config-connector/docs/reference/resource-docs/serviceusage/service)

## Usage

1.  Clone the package:
    ```
    kpt pkg get https://github.com/GoogleCloudPlatform/blueprints.git/catalog/networking/network/vpc@${VERSION}
    ```
    Replace `${VERSION}` with the desired repo branch or tag
    (for example, `main`).

1.  Move into the local package:
    ```
    cd "./vpc/"
    ```

1.  Edit the function config file(s):
    - setters.yaml

1.  Execute the function pipeline
    ```
    kpt fn render
    ```

1.  Initialize the resource inventory
    ```
    kpt live init --namespace ${NAMESPACE}"
    ```
    Replace `${NAMESPACE}` with the namespace in which to manage
    the inventory ResourceGroup (for example, `config-control`).

1.  Apply the package resources to your cluster
    ```
    kpt live apply
    ```

1.  Wait for the resources to be ready
    ```
    kpt live status --output table --poll-until current
    ```
