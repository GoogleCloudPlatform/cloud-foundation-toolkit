# Merge Dependabot Pull Reqeusts

## Prerequisites:
Install and Authenticate Github CLI (https://github.com/cli/cli)

## Usage
```
merge_dependabot_prs.sh
  -o [GitHub organization | Default: terraform-google-modules]
  -f [Repository name(s) contains filter | Default: NONE]
  -l [label to apply to failed checks | Default: dependabot-checks-failed]
```
