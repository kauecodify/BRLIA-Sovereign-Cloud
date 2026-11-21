# Manifestos Kubernetes

Manifestos base para o plano de controle BRLIA.

```bash
kubectl apply -f infra/kubernetes/namespace.yaml
kubectl apply -f infra/kubernetes/
```

ConfigMap `control-api-config` e Secret `control-api-secrets` devem ser
criados por loja de segredos (Vault) ou a partir do `.env` ambiente.
