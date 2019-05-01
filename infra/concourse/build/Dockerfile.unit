FROM alpine:3.8

RUN apk add --no-cache --update \
    bash \
    make \
    python=2.7.15-r1 \
    py-pip=10.0.1-r0

ADD ./build/data/requirements.txt .

RUN pip install -r requirements.txt
