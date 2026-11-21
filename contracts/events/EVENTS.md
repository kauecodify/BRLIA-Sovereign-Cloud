# Catalogo de Eventos BRLIA

Eventos sao publicados via NATS com JetStream. Assuntos seguem o padrao
`brlia.<dominio>.<evento>`.

## Assuntos

| Assunto                          | Produtor        | Consumidor esperado         |
|----------------------------------|-----------------|-----------------------------|
| `brlia.workload.submitted`       | waterfall       | compute, audit              |
| `brlia.workload.scheduled`       | waterfall       | compute, audit              |
| `brlia.workload.completed`       | waterfall       | billing, audit              |
| `brlia.workload.failed`          | waterfall       | operations-center, audit    |
| `brlia.resource.allocated`       | waterfall       | billing, audit              |
| `brlia.resource.released`        | waterfall       | billing, audit              |
| `brlia.identity.user.created`    | identity        | control-api, audit          |
| `brlia.identity.role.assigned`   | identity        | control-api, audit          |
| `brlia.billing.invoice.issued`   | billing         | customer-portal, audit      |
| `brlia.security.secret.rotated`  | security        | audit                       |

## Envelope

Todo evento usa o schema `schemas/event-envelope.schema.json`.
