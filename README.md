[![codecov](https://codecov.io/gh/loloDawit/go-pro/branch/main/graph/badge.svg)](https://codecov.io/gh/loloDawit/go-pro)


## Project Description: Comprehensive Backend API Development and Deployment

### Overview
This project focuses on developing a robust backend API in Golang, ensuring thorough testing, and deploying it using modern CI/CD practices with GitHub Actions. The deployment targets AWS ECS, utilizing AWS ECR for Docker image storage and Terraform for infrastructure management. The project emphasizes secure and efficient handling of environment variables, automated testing, and seamless integration and deployment workflows.

### Key Components

#### 1. **Backend API in Golang**
   - Developed a comprehensive backend API using Golang, incorporating best practices for performance and scalability.
   - Implemented JWT-based authentication for secure user management.

#### 2. **Automated Testing and CI/CD Integration**
   - Set up automated tests using Go’s testing framework and Codecov for coverage reporting.
   - Configured GitHub Actions workflows for continuous integration and deployment:
     - **Go CI Workflow**: Runs tests on code push and pull requests to ensure code quality and reliability.
     - **Docker Publish Workflow**: Builds and pushes Docker images to AWS ECR, triggered by the successful completion of the Go CI workflow.

#### 3. **Containerization and Deployment**
   - Utilized Docker for containerizing the Golang application.
   - Set up QEMU and Docker Buildx for multi-platform builds.
   - Configured AWS credentials and login for secure access to AWS services.
   - Built and tagged Docker images with both commit SHA and 'latest' tags for version control and easy rollbacks.

#### 4. **AWS Infrastructure Management with Terraform**
   - Deployed the application on AWS ECS Fargate using a robust Terraform configuration.
   - Managed infrastructure components including VPC, subnets, security groups, CloudWatch log groups, IAM roles, and ECS cluster setup.
   - Configured AWS ACM for SSL/TLS certificates to secure the application endpoints.

#### 5. **Advanced CI/CD Features**
   - Captured and logged success or failure status of AWS ECS service updates, providing detailed error outputs for troubleshooting.
   - Ensured that the Docker build and deployment steps only run if the Go tests pass, maintaining the integrity of the deployment process.
   - Utilized environment files to securely manage and pass configuration data during the CI/CD process.

### Achievements
   - Ensured robust and secure deployment practices using GitHub Actions and AWS services.
   - Automated the build, test, and deployment pipeline, significantly reducing manual intervention and potential errors.
   - Improved application security and reliability through thorough testing and secure handling of environment variables.

This project demonstrates a comprehensive approach to backend development, testing, and deployment, leveraging modern DevOps practices to achieve a highly efficient and reliable system.

Cost/day analysis for the services used on AWS

### 1. Amazon ECS with Fargate

**ECS Fargate Pricing:**
- **vCPU Hours:** \$0.04048 per vCPU hour
- **GB Hours:** \$0.004445 per GB hour

Assuming we have one task running with the following configuration:
- 0.25 vCPU
- 0.5 GB Memory

**Cost Calculation:**
- vCPU Cost: 0.25 vCPU * \$0.04048 per hour = \$0.01012 per hour
- Memory Cost: 0.5 GB * \$0.004445 per hour = \$0.0022225 per hour

**Total Cost per Task per Hour:**
- \$0.01012 + \$0.0022225 = \$0.0123425 per hour

**Total Cost per Task per Day:**
- \$0.0123425 * 24 hours = \$0.29622 per day

### 2. Amazon ECR

**ECR Pricing:**
- **Storage:** \$0.10 per GB-month

Assuming we use 1 GB of storage for your Docker images.

**Cost per Day:**
- \$0.10 / 30 days ≈ \$0.00333 per day

### 3. Amazon ACM

**ACM Pricing:**
- Free for public SSL/TLS certificates provisioned through ACM.

### 4. Amazon CloudWatch

**CloudWatch Pricing:**
- **Logs:** \$0.50 per GB ingested
- **Metrics:** \$0.30 per metric/month

Assuming we ingest 1 GB of logs per day and have 10 custom metrics.

**Log Ingestion Cost per Day:**
- \$0.50 per GB ≈ \$0.50 per day

**Metrics Cost per Day:**
- 10 metrics * \$0.30 / 30 days ≈ \$0.10 per day

### 5. VPC, Application Load Balancer, Subnet, Internet Gateway, Route Table

**VPC and Subnets:**
- There is no direct cost for VPC and subnets.

**Application Load Balancer (ALB):**
- **ALB Pricing:**
  - **Load Balancer Hours:** \$0.0225 per ALB hour
  - **LCU (Load Balancer Capacity Units):** \$0.008 per LCU hour

Assuming average usage of 1 LCU per hour.

**ALB Cost per Day:**
- Load Balancer Hours: \$0.0225 * 24 hours = \$0.54 per day
- LCU Hours: \$0.008 * 24 hours = \$0.192 per day

**Total ALB Cost per Day:**
- \$0.54 + \$0.192 = \$0.732 per day

**Internet Gateway and Route Table:**
- No direct cost.

### 6. Route 53

**Route 53 Pricing:**
- **Hosted Zone:** \$0.50 per hosted zone / month
- **DNS Queries:** \$0.40 per million queries

Assuming 1 hosted zone and 1 million DNS queries per day.

**Hosted Zone Cost per Day:**
- \$0.50 / 30 days ≈ \$0.0167 per day

**DNS Query Cost per Day:**
- 1 million queries * \$0.40 / 1 million ≈ \$0.40 per day

**Total Route 53 Cost per Day:**
- \$0.0167 + \$0.40 = \$0.4167 per day

### 7. Terraform Backup State

**Amazon S3:**
- **Storage:** \$0.023 per GB-month

Assuming 10 GB of storage.

**Storage Cost per Day:**
- 10 GB * \$0.023 / 30 days ≈ \$0.0077 per day

**Amazon DynamoDB:**
- **DynamoDB Pricing:**
  - **Write Capacity Unit:** \$0.00065 per WCU hour
  - **Read Capacity Unit:** \$0.00013 per RCU hour

Assuming 1 WCU and 1 RCU.

**Cost per Day:**
- WCU Cost: 1 WCU * \$0.00065 * 24 hours = \$0.0156 per day
- RCU Cost: 1 RCU * \$0.00013 * 24 hours = \$0.00312 per day

**Total DynamoDB Cost per Day:**
- \$0.0156 + \$0.00312 = \$0.01872 per day

### Total Estimated Cost per Day

Summing up all the estimated costs:

- **ECS Fargate:** \$0.29622 per day
- **ECR:** \$0.00333 per day
- **ACM:** Free
- **CloudWatch:** \$0.50 (logs) + \$0.10 (metrics) ≈ \$0.60 per day
- **ALB:** \$0.732 per day
- **Route 53:** \$0.4167 per day
- **S3:** \$0.0077 per day
- **DynamoDB:** \$0.01872 per day

**Total Estimated Cost per Day:**
- \$0.29622 + \$0.00333 + \$0.60 + \$0.732 + \$0.4167 + \$0.0077 + \$0.01872 ≈ \$2.08 per day

