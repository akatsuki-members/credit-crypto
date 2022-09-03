variable "environment_name" {
  description = "Name of the environment. e.g. dev, qa, stage, prod"
  default     = "dev"
}

variable "additional_tags" {
  default = {
    Service   = "audit"
    Component = "messaging"
    Scope     = "solution"
    Project   = "credit-crypto"
  }
  description = "Additional resource tags"
  type        = map(string)
}