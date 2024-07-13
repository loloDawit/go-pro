resource "aws_ecs_cluster" "go_pro_app_cluster" {
  name = "go-pro-app-cluster"
}

resource "aws_iam_role" "ecs_task_execution_role" {
  name = "ecsTaskExecutionRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy" "ecs_task_execution_policy" {
  role = aws_iam_role.ecs_task_execution_role.id

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:BatchCheckLayerAvailability"
        ],
        Resource = "*"
      }
    ]
  })
}

resource "aws_ecr_repository" "go_pro_app" {
  name = "go_pro_app"
}

resource "aws_cloudwatch_log_group" "go_pro_app_log_group" {
  name              = "/ecs/go_pro_app"
  retention_in_days = 1
}

resource "aws_ecs_task_definition" "go_pro_app_task" {
  family                   = "go-pro-app-task"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([
    {
      name      = "go_pro_app"
      image     = "${aws_ecr_repository.go_pro_app.repository_url}:latest"
      essential = true
      portMappings = [
        {
          containerPort = 8080
          hostPort      = 8080
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = "/ecs/go_pro_app"
          awslogs-region        = "us-east-1"
          awslogs-stream-prefix = "ecs"
        }
      }
      environment = [
        {
          name  = "DB_USER"
          value = var.db_user
        },
        {
          name  = "DB_PASSWORD"
          value = var.db_password
        },
        {
          name  = "DB_HOSTNAME"
          value = var.db_hostname
        },
        {
          name  = "DB_NAME"
          value = var.db_name
        },
        {
          name  = "JWT_SECRET"
          value = var.jwt_secret
        },
        {
          name  = "CONFIG_DIRECTORY"
          value = "/app/config"
        },
        # {
        #   name  = "ENV"
        #   value = "production"
        # },
        {
          name  = "DEPLOYMENT"
          value = "ecs"
        }
      ]
    }
  ])
}

resource "aws_ecs_service" "go_pro_app_service" {
  name            = "go-pro-app-service"
  cluster         = aws_ecs_cluster.go_pro_app_cluster.id
  task_definition = aws_ecs_task_definition.go_pro_app_task.arn
  desired_count   = 1

  network_configuration {
    subnets         = [var.subnet1_id, var.subnet2_id]
    security_groups = [var.ecs_sg_id]
    assign_public_ip = true
  }

  launch_type = "FARGATE"

  load_balancer {
    target_group_arn = var.target_group_arn
    container_name   = "go_pro_app"
    container_port   = 8080
  }

  depends_on = [var.http_listener_arn, var.https_listener_arn]
}

resource "aws_appautoscaling_target" "ecs" {
  max_capacity       = 10
  min_capacity       = 1
  resource_id        = "service/${aws_ecs_cluster.go_pro_app_cluster.name}/${aws_ecs_service.go_pro_app_service.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "scale_up" {
  name                   = "scale-up"
  policy_type            = "TargetTrackingScaling"
  resource_id            = aws_appautoscaling_target.ecs.resource_id
  scalable_dimension     = aws_appautoscaling_target.ecs.scalable_dimension
  service_namespace      = aws_appautoscaling_target.ecs.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }

    target_value = 50.0
    scale_in_cooldown  = 300
    scale_out_cooldown = 300
  }
}

resource "aws_appautoscaling_policy" "scale_down" {
  name                   = "scale-down"
  policy_type            = "TargetTrackingScaling"
  resource_id            = aws_appautoscaling_target.ecs.resource_id
  scalable_dimension     = aws_appautoscaling_target.ecs.scalable_dimension
  service_namespace      = aws_appautoscaling_target.ecs.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageMemoryUtilization"
    }

    target_value = 50.0
    scale_in_cooldown  = 300
    scale_out_cooldown = 300
  }
}

output "ecs_cluster_id" {
  value = aws_ecs_cluster.go_pro_app_cluster.id
}

output "ecs_task_definition_arn" {
  value = aws_ecs_task_definition.go_pro_app_task.arn
}

output "ecs_service_name" {
  value = aws_ecs_service.go_pro_app_service.name
}
