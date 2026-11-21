# Adaptador hsm

Adaptador de infraestrutura soberana BRLIA para **modulos de hardware de seguranca para chaves criptograficas nacionais**.

## Uso

```bash
export HSM_ENDPOINT=https://<alvo>   # endereco do alvo gerenciado
go run main.go
```

O adaptador consulta o alvo configurado e imprime um relatorio JSON de
`unconfigured`, `reachable`, `degraded` ou `unreachable`. A integracao
completa (provisionamento, teardown e reconciliacao) sera implementada
progressivamente a partir desse ponto de coleta de estado.

## Contrato

- Entrada: variavel de ambiente `HSM_ENDPOINT`.
- Saida: JSON com `adapter`, `endpoint`, `status`, `detail` e `checked_at`.
- Exit codes: `0` quando conclui a verificacao (mesmo em falha de rede) e
  `1` apenas em erros de execucao.
