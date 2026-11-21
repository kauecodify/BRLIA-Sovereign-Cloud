variable "environment" {
  description = "Nome do ambiente (dev, staging, production)"
  type        = string
  default     = "dev"
}

variable "region" {
  description = "Regiao soberana principal"
  type        = string
  default     = "br-sul-1"
}

variable "control_plane_count" {
  description = "Quantidade de nos do control plane"
  type        = number
  default     = 3
}

variable "worker_count" {
  description = "Quantidade de nos de trabalho"
  type        = number
  default     = 6
}
