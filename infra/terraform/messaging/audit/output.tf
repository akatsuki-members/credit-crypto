output "sqs_arn" {
  description = "SQS Arn"
  value       = aws_sqs_queue.audit_orders_queue.arn
}

output "sqs_url" {
  description = "SQS Url"
  value       = aws_sqs_queue.audit_orders_queue.id
}