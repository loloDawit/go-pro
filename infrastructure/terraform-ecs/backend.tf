terraform {
  backend "s3" {
    bucket         = "go-pro-app-terraform-state"
    key            = "go-pro-app/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}
