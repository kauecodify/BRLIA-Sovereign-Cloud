# OpenTofu

Provisionamento declarativo da base de infraestrutura.

```bash
cd infra/opentofu
tofu init
tofu plan  -var environment=dev
tofu apply -var environment=dev
```

O output `inventory_nodes` alimenta `infra/ansible/inventory`.
