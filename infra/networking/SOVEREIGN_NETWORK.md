# Rede Soberana BRLIA

## Princípios

1. Enderecamento proprio em faixas privadas por regiao, com roteamento entre
   datacenters nacionais por enlaces dedicados ou criptografados (IPsec/WireGuard).
2. Nenhum fluxo de gerenciamento cruza a internet aberta.
3. Segmentacao: `mgmt` (administracao), `control` (APIs internas),
   `data` (trafego de carga) e `storage` (trafego Ceph).
4. Observabilidade obrigatoria de todos os bordas (flow logs + metricas).

## Enderecamento planejado

| Regiao       | Sub-rede        | Uso              |
|--------------|-----------------|------------------|
| br-sul-1     | 10.20.0.0/16    | Datacenter sul   |
| br-suldeste-1| 10.21.0.0/16    | Datacenter sudeste |
| br-nordeste-1| 10.22.0.0/16    | Datacenter nordeste |
| docker dev   | 172.20.0.0/16   | Ambiente local   |

## Firewall (referencia nftables)

Regra base de referencia para nos:

```
table inet brlia {
  chain input {
    type filter hook input priority 0; policy drop;
    ct state established,related accept
    iif lo accept
    ip saddr 10.20.0.0/16 tcp dport { 22, 2379, 2380, 6443 } accept
  }
  chain forward {
    type filter hook forward priority 0; policy drop;
  }
}
```
