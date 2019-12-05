# CFT GitHub

## Operating

This module requires setting a GitHub personal access token.

First, create a [personal access token](https://help.github.com/en/enterprise/2.17/user/authenticating-to-github/creating-a-personal-access-token-for-the-command-line#creating-a-token).

Then, export it:

```
export GITHUB_TOKEN=aaaaaaa
```

Note, because of the many resources involved, you might need to run Terraform with `-refresh=false`.
