# Apply these pipelines

See full document at [go/cft-module-ci][cft-module-ci]

Take care to log in using:

    fly login --target cft -n cft -c https://concourse.infra.cft.tips

Validate the pipeline:

    fly -t cft validate-pipeline -c pipelines/<pipeline>.yml

Enforce the pipeline:

    make startup-scripts

[cft-module-ci]: http://goto.google.com/cft-module-ci
