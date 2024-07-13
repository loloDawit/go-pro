resource "aws_lb" "go_pro_app_lb" {
  name               = "go-pro-app-lb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [var.lb_sg_id]
  subnets            = [var.subnet1_id, var.subnet2_id]
}

resource "aws_lb_target_group" "go_pro_app_tg" {
  name     = "go-pro-app-tg"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = var.vpc_id
  target_type = "ip"

  health_check {
    path                = "/health"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
    matcher             = "200-399"
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.go_pro_app_lb.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.go_pro_app_tg.arn
  }
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.go_pro_app_lb.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"

  certificate_arn = var.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.go_pro_app_tg.arn
  }
}

output "alb_zone_id" {
  value = aws_lb.go_pro_app_lb.zone_id
}


output "alb_dns_name" {
  value = aws_lb.go_pro_app_lb.dns_name
}

output "http_listener_arn" {
  value = aws_lb_listener.http.arn
}

output "https_listener_arn" {
  value = aws_lb_listener.https.arn
}

output "target_group_arn" {
  value = aws_lb_target_group.go_pro_app_tg.arn
}
