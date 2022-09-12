module "messaging_module" {
  source                  = "./modules/messaging"
  orders_topic_name       = "orders"
  audit_orders_queue_name = "audit-orders"
}
