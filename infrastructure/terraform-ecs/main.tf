module "network" {
  source = "./network"
}

module "alb" {
  source = "./alb"

  vpc_id        = module.network.vpc_id
  subnet1_id    = module.network.subnet1_id
  subnet2_id    = module.network.subnet2_id
  lb_sg_id      = aws_security_group.lb_sg.id
  certificate_arn = module.route53.certificate_arn
}

module "ecs" {
  source = "./ecs"

  vpc_id            = module.network.vpc_id
  subnet1_id        = module.network.subnet1_id
  subnet2_id        = module.network.subnet2_id
  ecs_sg_id         = aws_security_group.ecs_sg.id
  target_group_arn  = module.alb.target_group_arn
  http_listener_arn = module.alb.http_listener_arn
  https_listener_arn = module.alb.https_listener_arn

  db_user           = var.db_user
  db_password       = var.db_password
  db_hostname       = var.db_hostname
  db_name           = var.db_name
  jwt_secret        = var.jwt_secret
}

module "route53" {
  source = "./route53"

  domain_name      = var.domain_name
  hosted_zone_name = var.hosted_zone_name
  alb_dns_name     = module.alb.alb_dns_name
  alb_zone_id      = module.alb.alb_zone_id
}


resource "aws_security_group" "ecs_sg" {
  vpc_id = module.network.vpc_id

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "lb_sg" {
  vpc_id = module.network.vpc_id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
