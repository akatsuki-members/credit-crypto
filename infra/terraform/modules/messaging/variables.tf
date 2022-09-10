variable "environment_name" {
  description = "Name of the environment. e.g. dev, qa, stage, prod"
  default     = "dev"
  type        = string
}

variable "additional_tags" {
  default = {
    Component = "messaging"
    Scope     = "solution"
    Project   = "credit-crypto"
  }
  description = "Additional resource tags"
  type        = map(string)
}

variable "orders_topic_name" {
  type    = string
  default = "orders"
}

variable "audit_orders_queue_name" {
  type    = string
  default = "audit-orders"
}
