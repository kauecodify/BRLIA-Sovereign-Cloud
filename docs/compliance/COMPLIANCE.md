# BRLIA Sovereign Cloud - Conformidade Regulatoria

## Escopo

Este documento descreve como a plataforma atende aos marcos regulatorios brasileiros aplicaveis a infraestrutura de nuvem soberana. Itens marcados como decisao de operador dependem da implantacao e servem de guia.

## Marco Legal Aplicavel

| Norma | Relevancia | Status na Implementacao |
| --- | --- | --- |
| LGPD (Lei 13.709/2018) | Dados pessoais de titulares | Controles tecnicos implementados; DPO e processos sao responsabilidade do operador |
| Marco Civil da Internet (Lei 12.965/2014) | Registros de acesso e dados de conexao | Logs de acesso e auditoria com retencao configuravel |
| Decreto 11.871/2023 (Seguranca Cibernetica) | Requisitos para servicos essenciais | Arquitetura de segmentacao e resposta a incidentes documentada |
| Decreto 10.046/2019 | Compartilhamento de dados no governo | API de `services/government` com trilha de auditoria por solicitacao |
| MP 2.200-2/2001 (ICP-Brasil) | Assinatura e certificados digitais | Suporte a HSM e certificados ICP-Brasil via `adapters/hsm` |

## LGPD - Controles Tecnicos

### Principios Implementados

- **Minimizacao** - coleta de dados limitada ao cadastro de organizacoes, projetos e usuarios do plano de controle
- **Finalidade** - dados sao usados apenas para operacao da plataforma e faturamento
- **Seguranca** - ver `docs/security/SECURITY.md` (criptografia, controle de acesso, auditoria)
- **Transparencia** - informacoes de tratamento disponiveis nos portais de cliente

### Direitos dos Titulares

O servico `identity` e o `customer-portal` expoe fluxos para:

1. Acesso aos dados do titular armazenados na conta
2. Correcao de dados cadastrais incorretos
3. Exclusao de dados dispensaveis (respeitando retencoes legais)
4. Portabilidade via exportacao organizada do projeto

### Retencao e Descarte

| Dado | Proposito | Retencao |
| --- | --- | --- |
| Diario de auditoria | Seguranca e investigacao | 7 anos (imutavel) |
| Dados de sessao | Autenticacao | expiracao em 24 horas |
| Metricas Prometheus | Monitoramento | 1 ano |
| Dados de cliente (objetos) | Decisao do cliente | Ciclo de vida gerenciado pelo cliente |

## Decreto de Seguranca Cibernetica

- **Gestao de riscos** - matriz de risco mantida em `docs/security/SECURITY.md` e revisada trimestralmente
- **Resposta a incidentes** - fluxo de 6 etapas definido, com testes periodicos (`tests/resilience/run.sh`)
- **Continuidade** - rotina de backup e DR em `scripts/` e `docs/operations/RUNBOOKS.md`
- **Subcontratacao** - todos os componentes executam em infraestrutura sob jurisdicao brasileira

## Comprovacao de Soberania

Devido a definicao da implementacao de regioes brasileiras:

- A regiao **br-sul-1** (padrao) e as regioes **br-suldeste-1** e **br-nordeste-1** ficam inteiramente dentro do territorio nacional
- O plano de rede (`infra/networking/SOVEREIGN_NETWORK.md`) usa nodes de propriedade do operador, sem transito internacional
- Telemetria e alertas ficam dentro do cluster soberano (Prometheus/Grafana locais)

## Auditorias

1. **Interna** - revisao trimestral de acessos, rotacao de chaves e testes de restauracao
2. **Externa** - recomendada antes da oferta formal para clientes governamentais
3. **Tecnica** - varreduras automatizadas de vulnerabilidade (`tests/security/run.sh`, govulncheck no CI)

## Pendencias Conhecidas

- Certificacao formal (ISO 27001, SOC 2) depende de implantacao e auditoria externa
- Mapa de fluxo de dados detalhado por titular e responsabilidade do operador de cada regiao
- Processo formal de DPO e canal de titulares sao requisitos organizacionais, nao de software
