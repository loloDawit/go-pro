variable "vpc_id" {
  description = "The VPC ID."
}

variable "subnet1_id" {
  description = "The first subnet ID."
}

variable "subnet2_id" {
  description = "The second subnet ID."
}

variable "ecs_sg_id" {
  description = "The ECS security group ID."
}

variable "target_group_arn" {
  description = "The target group ARN."
}

variable "http_listener_arn" {
  description = "The HTTP listener ARN."
}

variable "https_listener_arn" {
  description = "The HTTPS listener ARN."
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