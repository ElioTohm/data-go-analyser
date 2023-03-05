terraform {
  required_version = ">= 0.13.1"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.22"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = ">= 2.12, < 3.0"
    }
    null = {
      source  = "hashicorp/null"
      version = ">= 2.0"
    }
  }
}

provider "aws" {
  region = "eu-west-1"
}
data "aws_ecr_authorization_token" "token" {}

provider "docker" {
  registry_auth {
    address  = "321696485191.dkr.ecr.eu-west-1.amazonaws.com"
    username = data.aws_ecr_authorization_token.token.user_name
    password = data.aws_ecr_authorization_token.token.password
  }
}