resource "aws_sqs_queue" "audit_orders_queue" {
  name       = "audit-orders-queue"
  fifo_queue = true

  tags = merge(
    var.additional_tags,
    {
      Queue = "audit-orders"
    }
  )
}

resource "aws_sns_topic_subscription" "orders_sqs_audit" {
  topic_arn = data.aws_sns_topic.orders
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.audit_orders_queue.arn
}
