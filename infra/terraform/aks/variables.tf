variable "location" {
  type        = string
  description = "Azure region"
}

variable "resource_group_name" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "dns_prefix" {
  type = string
}

variable "node_count" {
  type    = number
  default = 2
}

variable "node_size" {
  type    = string
  default = "Standard_D4as_v5"
}

variable "tags" {
  type    = map(string)
  default = {}
}

variable "subscription_id" {
  type        = string
  description = "Azure Subscription ID"
}

variable "kubernetes_version" {
  type        = string
  description = "AKS Kubernetes Version"
}