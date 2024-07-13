variable "domain_name" {
  description = "The domain name for the ACM certificate."
}

variable "hosted_zone_name" {
  description = "The name of the hosted zone in Route 53."
}

variable "alb_dns_name" {
  description = "The DNS name of the ALB."
}

variable "alb_zone_id" {
  description = "The zone ID of the ALB."
}
