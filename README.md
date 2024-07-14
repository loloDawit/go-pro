[![codecov](https://codecov.io/gh/loloDawit/go-pro/branch/main/graph/badge.svg)](https://codecov.io/gh/loloDawit/go-pro)

Here's a more detailed project description that encapsulates all the work you've done:

---

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
