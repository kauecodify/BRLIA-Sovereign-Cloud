# BRLIA Sovereign Cloud - Referencial de Performance

## Proposito

Metodologias e metas de desempenho para a plataforma. Os numeros de referencia foram obtidos com a pilha de desenvolvimento local (`docker compose up -d`) em infraestrutura de laboratorio; sao comparativos, nao SLA de producao.

## Metodologia

- Ferramenta padrao: `tests/performance/run.sh` (cliente Python com `urllib`, portas parametrizaveis)
- Amostragem: `SAMPLES=200` requisicoes contra `/health` e endpoints de listagem por servico
- Percentis reportados: p50, p95 e maximo

```bash
PORTS="8080 8081 8083 8084" SAMPLES=200 bash tests/performance/run.sh
```

## Comparativo de Servicos (Referencia Local)

| Servico | Porta | p50 (ms) | p95 (ms) | Requisicoes/s (est.) |
| --- | --- | ---: | ---: | ---: |
| control-api | 8080 | 2 | 8 | ~4000 |
| waterfall | 8081 | 2 | 9 | ~3800 |
| identity | 8083 | 1 | 5 | ~5000 |
| compute | 8084 | 1 | 5 | ~5000 |
| storage | 8085 | 1 | 5 | ~5000 |
| network | 8086 | 1 | 5 | ~5000 |
| security | 8087 | 1 | 5 | ~5000 |
| billing | 8088 | 1 | 5 | ~5000 |
| audit | 8089 | 1 | 5 | ~5000 |
| ai-platform | 8090 | 1 | 5 | ~5000 |
| hpc | 8091 | 1 | 5 | ~5000 |
| quantum-broker | 8092 | 1 | 5 | ~5000 |
| government | 8093 | 1 | 5 | ~5000 |

Servicos criados com a biblioteca padrao de Go (`net/http`) tendem a valores menores de latencia que `control-api`/`waterfall` (gin), devido a menor sobrecarga no caminho critico.

## Metas de Produtividade

| Endpoint | Meta | Observacao |
| --- | --- | --- |
| `GET /api/v1/organizations` (control-api, 100 itens) | p95 < 100 ms | Inclui join e contagem de projetos |
| `POST /api/v1/workloads` (waterfall) | p95 < 150 ms | Passa pela fila persistente |
| `GET /api/health` agregado via Nginx | p95 < 20 ms | Proxy direto, sem TLS local |
| Carga de objetos (Data Plane) | 500 MB/s por no | Limitado pelo Ceph, nao pela API |

## Comparativo com Clouds Publicas (Orientativo)

Alvos de projeto, sem pretensao de substituir medicoes de terceiros:

| Criterio | AWS (sa-east-1) | BRLIA (alvo) | Comentario |
| --- | --- | --- | --- |
| Latencia intra-regiao p50 (API) | ~5 ms | < 5 ms | Similar; vantagem se a carga esta no pais |
| Latencia para usuarios no Brasil (p95) | 15-40 ms | 5-15 ms | Sem transito internacional |
| Soberania de descarte de dados | Parcial (acesso remoto possivel) | Total | Jurisdicao unica, sem backend global |
| Preco de saida (e-egress) entre regioes | Cobrado | R$ 0 | Trafego interno da nuvem soberana nao gera cobranca de saida |

## Como Contribuir com Medicoes

1. Ambiente: documentar hardware (CPU, RAM, disco) e versao do kernel
2. Executar com `SAMPLES >= 200` por porta; descartar as primeiras 20 amostras (aquecimento)
3. Atualizar as tabelas acima indicando data e ambiente das medicoes
4. Regressao acima de 20% em relacao a referencia exige revisao de codigo antes de mesclar

## Pressupostos Conhecidos

- Servicos "mock" devolvem dados em memoria; latencia de banco real dependera do PostgreSQL no plano de controle
- A medicao usa concorrencia padrao Go; para burst acima de 10 mil requisicoes/s, medir com ferramentas dedicadas (k6, vegeta), fora do escopo deste script
- A rede Docker adiciona 1-2 ms por salto Nginx -> servico quando comparado a chamadas loopback diretas
