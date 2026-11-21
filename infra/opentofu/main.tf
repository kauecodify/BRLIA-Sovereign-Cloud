# Provisionamento da base fisica da nuvem soberana.
# Os providers de nuvem/bare-metal serao plugados conforme cada datacenter
# nacional parceiro; provisoriamente o estado do inventario e exportado
# como outputs para o Ansible.


locals {
  cluster_nodes = concat(
    [for i in range(var.control_plane_count) : "${var.region}-control-${i}"],
    [for i in range(var.worker_count) : "${var.region}-worker-${i}"]
  )
}

output "inventory_nodes" {
  description = "Nos declarados para o Ansible"
  value       = local.cluster_nodes
}

output "environment" {
  value = var.environment
}
