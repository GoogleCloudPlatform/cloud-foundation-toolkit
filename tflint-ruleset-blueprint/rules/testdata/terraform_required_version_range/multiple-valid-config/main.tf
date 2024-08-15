terraform {
  required_version = ">=1.6"
}

terraform {
  required_version = ">= 1.6, < 2.0"
}

terraform {
  required_version = "~>1.6"
}

terraform {
  required_version = "~>1.9"
}

terraform {
  backend "gcs" {
    bucket = "UPDATE_ME"
    prefix = "UPDATE_ME"
  }
}

terraform {
}
