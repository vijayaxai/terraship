# Example AWS Infrastructure

provider "aws" {
  region = var.region
}

variable "region" {
  description = "AWS region"
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  default     = "dev"
}

# S3 Bucket with encryption
resource "aws_s3_bucket" "example" {
  bucket = "terraship-example-${var.environment}"

  tags = {
    Name        = "terraship-example"
    Environment = var.environment
    Owner       = "DevOps Team"
    Project     = "Terraship Demo"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "example" {
  bucket = aws_s3_bucket.example.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_versioning" "example" {
  bucket = aws_s3_bucket.example.id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_public_access_block" "example" {
  bucket = aws_s3_bucket.example.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# VPC
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name        = "terraship-vpc"
    Environment = var.environment
    Owner       = "DevOps Team"
    Project     = "Terraship Demo"
  }
}

# Private Subnet
resource "aws_subnet" "private" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "${var.region}a"

  tags = {
    Name        = "terraship-private-subnet"
    Environment = var.environment
    Owner       = "DevOps Team"
    Project     = "Terraship Demo"
    Type        = "Private"
  }
}

# Output values
output "bucket_name" {
  value = aws_s3_bucket.example.id
}

output "vpc_id" {
  value = aws_vpc.main.id
}
