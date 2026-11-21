# BRLIA Sovereign Cloud - Politica de Seguranca

## Visao Geral

A seguranca da BRLIA Sovereign Cloud e construida sobre o principio de defesa em profundidade, com controles em todas as camadas da plataforma: fisica, rede, host, aplicacao e dados. Toda a carga de trabalho permanece em territorio nacional, em conformidade com o marco legal brasileiro de protecao de dados.

## Modelo de Ameacas

| Vetor | Mitigacao |
| --- | --- |
| Acesso nao autorizado ao Control Plane | Autenticacao JWT, RBAC por papel, MFA obrigatorio para operadores |
| Exfiltracao de dados | Criptografia em transito (TLS 1.3) e em repouso (AES-256) |
| Comprometimento de host | Imagens imutaveis, nftables, hardening CIS, attestacao de no |
| Ameaca interna | Registro de auditoria imutavel, segregacao de funcoes, revisao periodica de acessos |
| DDoS e abuso de API | Rate limiting no Nginx, WAF, quota por organizacao |
| Cadeia de suprimentos | govulncheck e gitleaks no CI, imagens base minimas (alpine), SBOM |

## Camadas de Controle

### Identidade e Acesso

- `services/identity` emite e valida tokens OIDC/JWT (HMAC-SHA256, rotacao de chave)
- `services/control-api` aplica RBAC baseado em papeis: `admin`, `operator`, `developer`, `viewer`
- Contas de servico recebem escopos minimos; sem credenciais compartilhadas
- Sessoes expiram em 1 hora, com refresh token rotativo de 24 horas

### Seguranca de Servicos (Servico `security`)

- `services/security` orquestra segredos e certificados; em dev usa armazenamento em memoria
- HashiCorp Vault (`infra/compose`, porta 8200) e a referencia para producao
- HSM certificado ICP-Brasil via `adapters/hsm` para chaves de assinatura e raiz
- Nenhuma credencial e gravada em arquivo de configuracao; tudo via variaveis de ambiente

### Seguranca de Rede

- Segmentacao gerencial, de controle, de dados e de armazenamento (ver `infra/networking/SOVEREIGN_NETWORK.md`)
- Trafego entre servicos em rede privada Docker/Kubernetes; sem exposicao publica direta
- Apenas `control-api` (via `infra/nginx`) e os portais recebem trafego externo
- Padroes nftables default-deny com excecoes explicitas documentadas

### Protecao de Dados

- Repouso: AES-256 no Ceph (`adapters/ceph`), discos com dm-crypt nos nos
- Transito: TLS 1.3 em todas as fronteiras; TLS mutuo entre servicos do plano de dados
- Banco de dados: PostgreSQL com credenciais unicas por servico e TLS obrigatorio
- Auditoria: `services/audit` registra eventos em fluxo NATS imutavel com retencao de 7 anos

## Gerenciamento de Chaves

1. Chaves raiz geradas e mantidas em HSM; nunca saem do dispositivo
2. Chaves intermediarias por regiao e por servico, com rotacao automatica de 90 dias
3. Envelope-encryption para dados de cliente: DEK por objeto, KEK no Vault/HSM
4. Revogacao imediata via matriz de chaves do servico `security`

## Resposta a Incidentes

1. **Deteccao** - alertas Prometheus, anomalias do servico `audit`, avisos de integridade
2. **Triagem** - classificacao de severidade (S1 a S4) e designacao de responsavel
3. **Contencao** - isolamento do no, revogacao de credenciais, bloqueio no Nginx/WAF
4. **Erradicacao** - remedicao da causa raiz, rotacao forcada de segredos afetados
5. **Recuperacao** - restauracao via `scripts/restore.sh` ou `scripts/disaster-recovery.sh`
6. **Pos-incidente** - relatorio em ate 5 dias uteis; notificacao ANPD quando aplicavel (LGPD)

## Divulgacao Responsavel

Vulnerabilidades devem ser reportadas ao canal sigiloso acordado com o operador da plataforma. Resposta inicial em ate 48 horas. Programas de recompensa nao se aplicam nesta fase do projeto.

## Checklist de Seguranca para Implantacao

- [ ] TLS configurado em todos os pontos de entrada
- [ ] Credenciais padrao removidas do `.env`
- [ ] `SERVER_MODE=release` no compose de producao (`infra/compose/docker-compose.production.yml`)
- [ ] Backups testados com `scripts/backup.sh` e `scripts/restore.sh`
- [ ] `bash tests/security/run.sh` executado no pipeline
