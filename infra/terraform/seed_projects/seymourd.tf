module "seymourd_seed" {
  source = "../../modules/seed_project"

  username = "seymourd"

  owner_emails = ["user:seymourd@google.com"]

  org_id          = "${var.org_id}"
  seed_folder_id  = "${google_folder.seed_root_folder.name}"
  billing_account = "${module.variables.phoogle_billing_account}"
}
