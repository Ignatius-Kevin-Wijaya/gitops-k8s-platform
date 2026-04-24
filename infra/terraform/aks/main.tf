resource "azurerm_resource_group" "rg" {
  name     = var.resource_group_name
  location = var.location
  tags     = var.tags
}

resource "azurerm_kubernetes_cluster" "aks" {
  name                = var.cluster_name
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  dns_prefix          = var.dns_prefix

  kubernetes_version = var.kubernetes_version

  identity {
    type = "SystemAssigned"
  }

  default_node_pool {
    name       = "system"
    node_count = var.node_count
    vm_size    = var.node_size
  }

  network_profile {
    load_balancer_sku = "standard"
    network_plugin    = "kubenet"
  }

  tags = var.tags
}
