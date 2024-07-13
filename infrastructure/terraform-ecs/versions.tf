terraform {
  required_version = ">= 1.4.6, < 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.58.0, < 6.0.0"
    }
  }
}
