variable "region" {
  description = "The AWS region to deploy resources in."
  default     = "us-east-1"
}

variable "domain_name" {
  description = "The domain name for the ACM certificate."
  default     = "app.addiscoding.com"
}

variable "hosted_zone_name" {
  description = "The name of the hosted zone in Route 53."
  default     = "addiscoding.com"
}

variable "db_user" {
  description = "The database user"
  type        = string
}

variable "db_password" {
  description = "The database password"
  type        = string
  sensitive   = true
}

variable "db_hostname" {
  description = "The database hostname"
  type        = string
}

variable "db_name" {
  description = "The database name"
  type        = string
}

variable "jwt_secret" {
  description = "JWT secret"
  type        = string
  sensitive   = true
}