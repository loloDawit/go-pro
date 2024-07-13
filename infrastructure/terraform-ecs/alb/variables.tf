variable "vpc_id" {
  description = "The VPC ID."
}

variable "subnet1_id" {
  description = "The first subnet ID."
}

variable "subnet2_id" {
  description = "The second subnet ID."
}

variable "lb_sg_id" {
  description = "The Load Balancer security group ID."
}

variable "certificate_arn" {
  description = "The ARN of the ACM certificate."
}
