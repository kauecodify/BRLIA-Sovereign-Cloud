# Contratos da BRLIA Sovereign Cloud

Este diretorio concentra os contratos publicos da plataforma.

## Estrutura

- `openapi/`: especificacoes OpenAPI 3.0 dos servicos.
- `events/`: catalogo de eventos assinados via NATS.
- `schemas/`: schemas JSON reutilizados por eventos e contratos.

## Regras

1. Toda mudanca de contrato passa por PR e revisao.
2. Mudancas incompativeis exigem incremento de versao major em `info.version`
   e periodo de transicao com suporte as duas versoes.
3. Eventos devem incluir `event_id`, `event_type`, `occurred_at` e `producer`.
