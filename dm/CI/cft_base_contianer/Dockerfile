FROM gcr.io/cloud-builders/gcloud

RUN apt-get update && apt-get -y install make \
    && apt-get -y install gettext-base \
    && pip install --upgrade pip \
    && pip install setuptools \
    && git clone https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit \
    && cd cloud-foundation-toolkit/dm \
    && rm -rf templates \
    && make prerequisites \
    && make build \
    && make install \
    && pip install tox \
    && pip install pytest \
    && make cft-venv \
    && make template-prerequisites \
    && make cft-prerequisites \
    && . venv/bin/activate \
    && ./src/cftenv \
    && pwd \
    && cft --version \
    && bats -v

WORKDIR /cloud-foundation-toolkit/dm
