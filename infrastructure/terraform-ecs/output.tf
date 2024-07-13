output "alb_dns_name" {
  description = "The DNS name of the ALB"
  value       = module.alb.alb_dns_name
}

output "ecs_service_name" {
  description = "The name of the ECS service"
  value       = module.ecs.ecs_service_name
}

output "route53_record_name" {
  description = "The name of the Route 53 record"
  value       = module.route53.route53_record_name
}
