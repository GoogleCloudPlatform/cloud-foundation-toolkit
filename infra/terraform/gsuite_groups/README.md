# gsuite_groups

## Requirements

- A super admin account within `phoogle.net`
- A project with the `admin.googleapis.com` API enabled

## Quickstart

Run this with your `phoogle.net` credentials.

```
gcloud auth application-default login \
  --scopes \
https://www.googleapis.com/auth/admin.directory.group,\
https://www.googleapis.com/auth/admin.directory.group.member,
```
