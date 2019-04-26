module "variables" { source = "../variables" }

variable "postgres_concourse_user_password" {
  description = "PostgreSQL password to be associated with the concourse user."
}
