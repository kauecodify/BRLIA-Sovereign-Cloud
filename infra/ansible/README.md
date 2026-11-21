# Ansible

Configuracao de nos bare-metal/VMs da nuvem soberana.

```bash
cd infra/ansible
ansible-playbook -i inventory/hosts.yml playbook.yml --check
ansible-playbook -i inventory/hosts.yml playbook.yml
```

O inventario deve ser gerado a partir do output do OpenTofu em producao.
