# BRLIA Sovereign Cloud

Plataforma de computacao em nuvem soberana brasileira com arquitetura propria, capacidade de operacao global e governanca institucional.

## Visao Geral

A BRLIA Sovereign Cloud e uma infraestrutura de computacao em nuvem de nova geracao, projetada para competir diretamente com AWS, Azure e Google Cloud, superando-as em indicadores tecnicos, economicos, operacionais e estrategicos.

## Estrutura do Projeto

```
brlia-sovereign-cloud/
├── apps/                    # Aplicacoes de interface
│   ├── customer-portal/     # Portal do cliente
│   ├── admin-console/       # Console administrativo
│   ├── developer-portal/    # Portal do desenvolvedor
│   └── operations-center/   # Centro de operacoes
│
├── services/                # Servicos principais
│   ├── control-api/         # Control Plane
│   ├── waterfall/           # Motor de orquestracao
│   ├── identity/            # Gerenciamento de identidade
│   ├── compute/             # Compute Plane
│   ├── storage/             # Data Plane
│   ├── network/             # Rede
│   ├── security/            # Seguranca
│   ├── billing/             # Faturamento
│   ├── audit/               # Auditoria
│   ├── ai-platform/         # Plataforma de IA
│   ├── hpc/                 # Computacao de alto desempenho
│   ├── quantum-broker/      # Broker de computacao quantica
│   └── government/          # Servicos governamentais
│
├── adapters/                # Adaptadores de infraestrutura
│   ├── kubernetes/
│   ├── kvm/
│   ├── ceph/
│   ├── slurm/
│   ├── hsm/
│   └── datacenter/
│
├── infra/                   # Infraestrutura
│   ├── compose/
│   ├── kubernetes/
│   ├── opentofu/
│   ├── ansible/
│   ├── monitoring/
│   └── networking/
│
├── contracts/               # Contratos de API
│   ├── openapi/
│   ├── events/
│   └── schemas/
│
├── database/
│   ├── migrations/
│   └── seeds/
│
├── docs/                    # Documentacao
│   ├── architecture/
│   ├── security/
│   ├── operations/
│   ├── compliance/
│   ├── benchmarks/
│   ├── deployment/
│   └── installation/
│
├── tests/                   # Testes
│   ├── unit/
│   ├── integration/
│   ├── e2e/
│   ├── security/
│   ├── performance/
│   └── resilience/
│
├── scripts/                 # Scripts de automacao
│
├── LICENSE
└── README.md
```

## Inicio Rapido

Consulte o [Manual Tecnico de Instalacao Nacional](docs/installation/NATIONAL_INSTALLATION_GUIDE.md) para instrucoes detalhadas de deploy.

## Arquitetura

- **BRLIA Control**: Gerenciamento global de organizacoes, clientes, projetos e recursos
- **BRLIA Waterfall**: Motor inteligente de orquestracao e distribuicao de recursos
- **BRLIA Compute**: Execucao de cargas de trabalho (VMs, containers, Kubernetes, GPU, HPC)
- **BRLIA Data**: Infraestrutura de dados (armazenamento, bancos, streaming, backup)
- **BRLIA Operations**: Operacao autonoma e continua (monitoramento, SRE, FinOps)

## Tecnologias

- **Container Runtime**: containerd, Kubernetes
- **Virtualizacao**: KVM
- **Armazenamento**: Ceph
- **Banco de Dados**: PostgreSQL, Redis
- **Mensageria**: NATS
- **Observabilidade**: Prometheus, Grafana, OpenTelemetry
- **Autenticacao**: Keycloak, OAuth 2.0, OIDC
- **Seguranca**: Vault, HSM
- **Infraestrutura como Codigo**: OpenTofu, Ansible

## Contribuindo

Este e um projeto de infraestrutura critica nacional.

## Licenca

BSD 3-Clause License - Copyright (c) 2026, BRLIA
