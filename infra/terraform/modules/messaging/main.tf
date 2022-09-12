resource "aws_sns_topic" "orders_fifo" {
  name            = "${var.orders_topic_name}.fifo"
  fifo_topic      = true
  delivery_policy = <<EOF
{
  "http": {
    "defaultHealthyRetryPolicy": {
      "minDelayTarget": 20,
      "maxDelayTarget": 20,
      "numRetries": 3,
      "numMaxDelayRetries": 0,
      "numNoDelayRetries": 0,
      "numMinDelayRetries": 0,
      "backoffFunction": "linear"
    },
    "disableSubscriptionOverrides": false,
    "defaultThrottlePolicy": {
      "maxReceivesPerSecond": 1
    }
  }
}
EOF
  tags = merge(
    var.additional_tags,
    {
      Topic = var.orders_topic_name
    }
  )
}

resource "aws_sqs_queue" "audit_orders" {
  name       = "${var.audit_orders_queue_name}.fifo"
  fifo_queue = true

  tags = merge(
    var.additional_tags,
    {
      Queue   = var.audit_orders_queue_name
      Service = "audit"
    }
  )
}

resource "aws_sns_topic_subscription" "orders_sqs_audit" {
  topic_arn = aws_sns_topic.orders_fifo.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.audit_orders.arn
}
