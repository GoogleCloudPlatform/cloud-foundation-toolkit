# Run blueprint-tests

Set environment variables:

```bash
export TF_VAR_org_id="your_org_id"
export TF_VAR_folder_id="your_folder_id"
export TF_VAR_billing_account="your_billing_account_id"
```

Create test project:

```bash
terraform -chdir=setup init
terraform -chdir=setup apply
```

Run tests:

```bash
go test [-v]
```

Cleanup test project:

```bash
terraform -chdir=setup destroy
```
