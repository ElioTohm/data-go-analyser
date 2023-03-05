locals {
  name = "elio-s3-test"
  ecr_image_tag = "0.0.6"
  tags = {
    Environment = "test"
  }
}

##########################################
# Lambda Function (with various triggers)
##########################################

module "docker_image" {
  source  = "terraform-aws-modules/lambda/aws//modules/docker-build"
  version = "4.10.1"

  create_ecr_repo  = true
  ecr_force_delete = true
  ecr_repo         = local.name
  ecr_repo_lifecycle_policy = jsonencode({
    "rules" : [
      {
        "rulePriority" : 1,
        "description" : "Keep only the last 2 images",
        "selection" : {
          "tagStatus" : "any",
          "countType" : "imageCountMoreThan",
          "countNumber" : 2
        },
        "action" : {
          "type" : "expire"
        }
      }
    ]
  })
  platform    = "linux/amd64"
  image_tag   = local.ecr_image_tag
  source_path = ".."
}

module "lambda_function" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "4.10.1"

  function_name  = "trigger"
  image_uri      = module.docker_image.image_uri
  package_type   = "Image"
  architectures  = ["x86_64"]
  create_package = false
  attach_policy_json = true
  policy_json        = <<-EOT
    {
      "Version": "2012-10-17",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "sqs:GetQueueUrl",
            "sqs:SendMessage"
          ],
          "Resource": "${aws_sqs_queue.queue.arn}"
        },
        {
          "Effect": "Allow",
          "Action": "s3:GetObject",
          "Resource": "${aws_s3_bucket.bucket.arn}/*"
          
        }
      ]
    }
  EOT

}

##################################
# Cloudwatch Events (EventBridge)
##################################
resource "aws_cloudwatch_event_rule" "scan_ami" {
  name          = "EC2CreateImageEvent"
  description   = "EC2 Create Image Event..."
  event_pattern = <<EOF
  {
    "source": ["aws.ec2"],
    "detail-type": ["AWS API Call via CloudTrail"],
    "detail": {
      "eventSource": ["ec2.amazonaws.com"],
      "eventName": ["CreateImage"]
    }
  }
  EOF
}

resource "aws_s3_bucket" "bucket" {
  bucket = local.name
}

resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_function.lambda_function_arn
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.bucket.arn
}

resource "aws_s3_bucket_notification" "my-trigger" {
  bucket = aws_s3_bucket.bucket.id

  lambda_function {
    lambda_function_arn = module.lambda_function.lambda_function_arn
    events              = ["s3:ObjectCreated:*"]
    filter_suffix       = ".csv"
  }
}

resource "aws_sqs_queue" "queue" {
  name                      = local.name
  delay_seconds             = 90
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10

  tags = local.tags
}
