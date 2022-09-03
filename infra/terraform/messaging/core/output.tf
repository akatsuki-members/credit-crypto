output "sns_id" {
  value       = aws_sns_topic.orders-fifo.id
  description = "SNS Topic Id"
}

output "sns_arn" {
  value       = aws_sns_topic.orders-fifo.arn
  description = "SNS Topic Arn"
}