output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}

output "cluster_name" {
  value = azurerm_kubernetes_cluster.aks.name
}

output "node_resource_group" {
  value = azurerm_kubernetes_cluster.aks.node_resource_group
}
