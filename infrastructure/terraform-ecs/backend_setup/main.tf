provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "terraform_state" {
  bucket = "go-pro-app-terraform-state"

  tags = {
    Name        = "terraform-state"
    Environment = "production"
  }
}

resource "aws_s3_bucket_versioning" "terraform_state_versioning" {
  bucket = aws_s3_bucket.terraform_state.bucket

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "terraform_state_lifecycle" {
  bucket = aws_s3_bucket.terraform_state.bucket

  rule {
    id     = "log"
    status = "Enabled"

    noncurrent_version_transition {
      storage_class = "GLACIER"
      noncurrent_days = 30
    }

    noncurrent_version_expiration {
      noncurrent_days = 365
    }
  }
}

resource "aws_dynamodb_table" "terraform_locks" {
  name         = "terraform-locks"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = {
    Name        = "terraform-locks"
    Environment = "development"
  }
}
