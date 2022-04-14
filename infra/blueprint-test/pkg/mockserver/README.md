# Mock Server

This package allows tests to be run much faster by using a [mock server](https://www.mock-server.com/)
instead of directly requesting results from the GCP APIs.

In particular, if a request hasn't changed since the last time it was run
then a response can be returned immediately from the cache.

## To Do

- [ ] Improve security by using a dynamically generated cert and inserting it into gcloud/Terraform.

## Logic

When matching requests, we should ignore:

- Request:
  - headers.authorization
  - headers.User-Agent
  - headers.X-Goog-User-Project
- Responses:
  - headers (all)

