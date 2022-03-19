terraform {
  backend "s3" {
    bucket = "app-terraform-bucket"
    region = "eu-west-1"
    key    = "terraform.tfstate"
  }
}
provider "aws" {
  region = "eu-west-1"
}
data "aws_caller_identity" "current" {}

locals {
  account_id     = data.aws_caller_identity.current.account_id
  environment    = "test"
  lambda_handler = "whatson"
  name           = "whatson"
  region         = "eu-west-1"
}

data "aws_s3_object" "lambda_source" {
  bucket = "sab-lambda-artifact"
  key    = "whatson.zip"
}

data "aws_iam_policy_document" "lambda_policy_document" {
  policy_id = "${local.name}-lambda"
  version   = "2012-10-17"
  statement {
    effect  = "Allow"
    actions = [
      "sts:AssumeRole"
    ]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }

  statement {
    effect  = "Allow"
    actions = [
      "secretsmanager:DescribeSecret",
      "secretsmanager:GetSecretValue",
      "secretsmanager:ListSecrets"
    ]
    resources = ["arn:aws:secretsmanager:${local.region}:${local.account_id}:secret:test/secret/*"]
  }
}

resource "aws_iam_role" "lambda_role" {
  name               = "${local.name}-lambda"
  assume_role_policy = data.aws_iam_policy_document.lambda_policy_document.json
}

data "aws_iam_policy_document" "logs" {
  policy_id = "${local.name}-lambda-logs"
  version   = "2012-10-17"
  statement {
    effect  = "Allow"
    actions = ["logs:CreateLogStream", "logs:PutLogEvents"]

    resources = [
      "arn:aws:logs:${local.region}:${local.account_id}:log-group:/aws/lambda/${local.name}*:*"
    ]
  }
}

resource "aws_iam_policy" "logs" {
  name   = "${local.name}-lambda-logs"
  policy = data.aws_iam_policy_document.logs.json
}

resource "aws_iam_role_policy_attachment" "logs" {
  depends_on = [aws_iam_role.lambda_role, aws_iam_policy.logs]
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.logs.arn
}

resource "aws_cloudwatch_log_group" "log" {
  name              = "/aws/lambda/${local.name}"
  retention_in_days = 7
}

resource "aws_lambda_function" "whatson-lambda" {
  s3_bucket        = "sab-lambda-artifact"
  s3_key           = "whatson.zip"
  function_name    = local.name
  role             = aws_iam_role.lambda_role.arn
  handler          = local.lambda_handler
  source_code_hash = filebase64sha256(data.aws_s3_object.lambda_source)
  runtime          = "go1.x"
  memory_size      = 1024
  timeout          = 30
}




