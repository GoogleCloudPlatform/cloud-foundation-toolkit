variable "test" {
  default = ""
}

resource "random_pet" "hello" {
}

output "test" {
  value = var.test
}

output "current_ws" {
  value = terraform.workspace
}
