variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

variable "environment_name" {
  description = "Name of the environment. e.g. dev, qa, stage, prod"
  default     = "dev"
}

variable "orders_topic_name" {
  description = "Name of the orders topic"
  default     = "orders"
  type        = string
}
