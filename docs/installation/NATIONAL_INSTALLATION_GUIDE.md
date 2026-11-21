# BRLIA Sovereign Cloud - Manual Tecnico de Instalacao Nacional

## Visao Geral

Manual tecnico para instalacao, configuracao e operacao do BRLIA Sovereign Cloud em ambiente nacional.

## Requisitos Minimos

### Desenvolvimento
- 8 vCPUs, 32 GB RAM, 500 GB SSD
- Ubuntu 22.04 LTS ou RHEL 9
- Rede 1 Gbps

### Producao Piloto
- 3 servidores (16 vCPUs, 64 GB RAM, 2 TB NVMe)
- Balanceador de carga: 2 servidores (4 vCPUs, 8 GB RAM)
- Armazenamento distribuido: 10 TB
- Rede 10 Gbps

### Producao Completa
- Data Center Primario: 50+ servidores, 10+ racks
- Data Center Secundario: 30+ servidores, 6+ racks
- Armazenamento: 1 PB+
- Rede: 100 Gbps backbone

## Dependencias de Software

- Linux Kernel 5.15+
- Docker 24.0+, containerd 1.6+
- Kubernetes 1.28+
- PostgreSQL 15+, Redis 7.0+
- Ceph 17+
- NATS 2.9+
- Prometheus 2.47+, Grafana 10.0+
- Keycloak 22+, Vault 1.15+
- Go 1.21+, Node.js 18+

## Arquitetura de Deploy

### Topologia Recomendada

Data Center Primario (BR-SP-1 - Sao Paulo):
- Load Balancers (HAProxy)
- Control Plane (Kubernetes)
- Compute Nodes (KVM/Kubernetes)
- Database (PostgreSQL com Patroni)
- Message Broker (NATS)
- Storage (Ceph)
- Monitoring (Prometheus, Grafana)
- Identity (Keycloak)
- Security (Vault)

Data Center Secundario (BR-RJ-1 - Rio de Janeiro):
- Mesma arquitetura do Primario
- Replicacao de dados

Conectividade:
- Link dedicado 10 Gbps entre DCs
- Multiplos provedores de internet com BGP
- Conexao direta com IX.br

## Etapas de Instalacao

### Fase 1: Preparacao do Ambiente

#### Configuracao do Sistema Operacional

```bash
# Atualizar sistema
sudo apt update && sudo apt upgrade -y

# Instalar dependencias basicas
sudo apt install -y curl wget git build-essential jq htop net-tools

# Configurar timezone e locale
sudo timedatectl set-timezone America/Sao_Paulo
sudo apt install -y locales
sudo locale-gen pt_BR.UTF-8
sudo update-locale LANG=pt_BR.UTF-8 LC_ALL=pt_BR.UTF-8

# Configurar limites do sistema
sudo bash -c 'cat > /etc/security/limits.d/99-brlia.conf <<EOF
* soft nofile 1048576
* hard nofile 1048576
* soft nproc 65535
* hard nproc 65535
* soft memlock unlimited
* hard memlock unlimited
EOF'

# Desabilitar swap
sudo swapoff -a
sudo sed -i '/ swap / s/^/#/' /etc/fstab

# Configurar kernel parameters
sudo bash -c 'cat > /etc/sysctl.d/99-brlia.conf <<EOF
vm.swappiness=0
vm.max_map_count=262144
fs.file-max=1048576
net.ipv4.ip_local_port_range=10240 65535
net.ipv4.ip_forward=1
net.core.somaxconn=10000
net.core.netdev_max_backlog=5000
EOF'

sudo sysctl -p /etc/sysctl.d/99-brlia.conf
```

#### Configuracao de Rede

```bash
# Exemplo de configuracao para Ubuntu
sudo bash -c 'cat > /etc/netplan/99-brlia.yaml <<EOF
network:
  version: 2
  renderer: networkd
  ethernets:
    eth0:
      dhcp4: no
      addresses: [192.168.1.10/24]
      gateway4: 192.168.1.1
      nameservers:
        addresses: [8.8.8.8, 8.8.4.4]
    eth1:
      dhcp4: no
      addresses: [10.0.0.10/24]
EOF'

sudo netplan apply
```

#### Configuracao de Seguranca

```bash
# Instalar e configurar firewall
sudo apt install -y ufw
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow http
sudo ufw allow https
sudo ufw allow 6443/tcp  # Kubernetes API
sudo ufw allow 8080:8081/tcp  # BRLIA APIs
sudo ufw allow 5432/tcp  # PostgreSQL
sudo ufw allow 6379/tcp  # Redis
sudo ufw allow 4222/tcp  # NATS
sudo ufw allow 9090/tcp  # Prometheus
sudo ufw allow 3000/tcp  # Grafana
sudo ufw enable

# Configurar SSH
sudo bash -c 'cat > /etc/ssh/sshd_config.d/99-brlia.conf <<EOF
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
MaxAuthTries 3
LoginGraceTime 60
ClientAliveInterval 300
ClientAliveCountMax 2
EOF'

sudo systemctl restart sshd
```

### Fase 2: Instalacao de Dependencias

#### Docker

```bash
# Remover versoes antigas
sudo apt remove -y docker docker-engine docker.io containerd runc

# Instalar Docker
sudo apt install -y ca-certificates curl gnupg lsb-release
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Configurar Docker
sudo bash -c 'cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m",
    "max-file": "3"
  },
  "storage-driver": "overlay2",
  "default-ulimits": {
    "nofile": {
      "Name": "nofile",
      "Hard": 1048576,
      "Soft": 1048576
    }
  }
}
EOF'

sudo systemctl enable docker
sudo systemctl start docker
sudo systemctl restart docker
sudo usermod -aG docker $USER
newgrp docker
```

#### Kubernetes

```bash
# Instalar kubeadm, kubelet e kubectl
sudo apt install -y apt-transport-https ca-certificates curl
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.28/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg

echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.28/deb/ /" | sudo tee /etc/apt/sources.list.d/kubernetes.list

sudo apt update
sudo apt install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl

# Inicializar cluster (no primeiro node)
sudo kubeadm init --pod-network-cidr=10.244.0.0/16 --apiserver-advertise-address=192.168.1.10

# Configurar kubectl
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Instalar rede CNI (Calico)
kubectl apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml

# Verificar status
kubectl get nodes
kubectl get pods -A
```

#### PostgreSQL

```bash
# Instalar PostgreSQL
sudo apt install -y postgresql-15 postgresql-client-15 postgresql-contrib-15

# Configurar PostgreSQL
sudo -u postgres psql -c "CREATE USER brlia WITH PASSWORD 'secure-password';"
sudo -u postgres psql -c "CREATE DATABASE brlia_control WITH OWNER brlia;"
sudo -u postgres psql -c "CREATE DATABASE brlia_waterfall WITH OWNER brlia;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE brlia_control TO brlia;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE brlia_waterfall TO brlia;"

# Configurar pg_hba.conf
sudo bash -c 'cat >> /etc/postgresql/15/main/pg_hba.conf <<EOF
# BRLIA Access
host    brlia_control    brlia    192.168.1.0/24    md5
host    brlia_waterfall brlia    192.168.1.0/24    md5
host    brlia_control    brlia    10.0.0.0/24      md5
host    brlia_waterfall brlia    10.0.0.0/24      md5
EOF'

# Configurar postgresql.conf
sudo sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/g" /etc/postgresql/15/main/postgresql.conf
sudo sed -i "s/#max_connections = 100/max_connections = 500/g" /etc/postgresql/15/main/postgresql.conf
sudo sed -i "s/#shared_buffers = 128MB/shared_buffers = 4GB/g" /etc/postgresql/15/main/postgresql.conf

sudo systemctl restart postgresql
sudo systemctl enable postgresql
```

#### Redis

```bash
# Instalar Redis
sudo apt install -y redis-server

# Configurar Redis
sudo bash -c 'cat > /etc/redis/redis.conf <<EOF
bind 0.0.0.0
protected-mode yes
port 6379
tcp-keepalive 300
database 16
maxmemory 8gb
maxmemory-policy allkeys-lru
appendonly yes
EOF'

sudo systemctl restart redis
sudo systemctl enable redis
```

#### NATS

```bash
# Baixar e instalar NATS
wget https://github.com/nats-io/nats-server/releases/download/v2.9.22/nats-server-v2.9.22-linux-amd64.zip
unzip nats-server-v2.9.22-linux-amd64.zip
sudo mv nats-server-v2.9.22-linux-amd64 /usr/local/bin/nats-server
sudo chmod +x /usr/local/bin/nats-server

# Configurar NATS
sudo bash -c 'cat > /etc/nats-server.conf <<EOF
listen: 0.0.0.0:4222
http: 0.0.0.0:8222
jetstream: true
store_dir: /var/lib/nats
max_payload: 16MB
max_control_line: 32KB
max_connections: 10000
EOF'

# Criar servico systemd
sudo bash -c 'cat > /etc/systemd/system/nats-server.service <<EOF
[Unit]
Description=NATS Server
After=network.target

[Service]
Type=simple
User=nats
Group=nats
ExecStart=/usr/local/bin/nats-server --config /etc/nats-server.conf
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF'

sudo useradd nats
sudo mkdir -p /var/lib/nats
sudo chown nats:nats /var/lib/nats
sudo systemctl daemon-reload
sudo systemctl start nats-server
sudo systemctl enable nats-server
```

#### Keycloak

```bash
# Baixar e instalar Keycloak
wget https://github.com/keycloak/keycloak/releases/download/22.0.0/keycloak-22.0.0.tar.gz
tar -xzf keycloak-22.0.0.tar.gz
sudo mv keycloak-22.0.0 /opt/keycloak

# Criar usuario
sudo useradd -r -s /bin/false keycloak
sudo chown -R keycloak:keycloak /opt/keycloak

# Criar servico systemd
sudo bash -c 'cat > /etc/systemd/system/keycloak.service <<EOF
[Unit]
Description=Keycloak Identity and Access Management
After=network.target postgresql.service

[Service]
User=keycloak
Group=keycloak
ExecStart=/opt/keycloak/bin/kc.sh start --db=postgres --db-url=jdbc:postgresql://localhost:5432/keycloak --db-username=keycloak --db-password=secure-password --hostname=brlia-control-01 --http-port=8082 --https-port=8445
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF'

# Criar banco de dados
sudo -u postgres psql -c "CREATE USER keycloak WITH PASSWORD 'secure-password';"
sudo -u postgres psql -c "CREATE DATABASE keycloak WITH OWNER keycloak;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak;"

sudo systemctl daemon-reload
sudo systemctl start keycloak
sudo systemctl enable keycloak
```

#### Vault

```bash
# Baixar e instalar Vault
wget https://releases.hashicorp.com/vault/1.15.0/vault_1.15.0_linux_amd64.zip
unzip vault_1.15.0_linux_amd64.zip
sudo mv vault /usr/local/bin/
sudo chmod +x /usr/local/bin/vault

# Criar usuario
sudo useradd -r -s /bin/false vault
sudo mkdir -p /opt/vault/data
sudo chown -R vault:vault /opt/vault

# Criar servico systemd
sudo bash -c 'cat > /etc/systemd/system/vault.service <<EOF
[Unit]
Description=Vault Secret Manager
After=network.target

[Service]
User=vault
Group=vault
ExecStart=/usr/local/bin/vault server -config=/etc/vault.d/vault.hcl
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF'

# Configurar Vault
sudo mkdir -p /etc/vault.d
sudo bash -c 'cat > /etc/vault.d/vault.hcl <<EOF
storage "file" {
  path = "/opt/vault/data"
}

listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = true
}

disable_mlock = true
ui = true

api_addr = "http://192.168.1.10:8200"
EOF'

sudo systemctl daemon-reload
sudo systemctl start vault
sudo systemctl enable vault

# Inicializar Vault
vault operator init -key-shares=5 -key-threshold=3
```

#### Prometheus e Grafana

```bash
# Instalar Prometheus
sudo apt install -y prometheus prometheus-node-exporter

# Configurar Prometheus
sudo bash -c 'cat > /etc/prometheus/prometheus.yml <<EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node-exporter'
    static_configs:
      - targets: ['localhost:9100']
  - job_name: 'brlia-control'
    static_configs:
      - targets: ['localhost:8080']
  - job_name: 'brlia-waterfall'
    static_configs:
      - targets: ['localhost:8081']
EOF'

sudo systemctl restart prometheus
sudo systemctl enable prometheus

# Instalar Grafana
sudo apt install -y grafana

# Configurar Grafana
sudo bash -c 'cat > /etc/grafana/grafana.ini <<EOF
[server]
http_port = 3000
http_addr = 0.0.0.0

[database]
type = postgres
host = localhost:5432
name = grafana
user = grafana
password = secure-password
EOF'

# Criar banco de dados
sudo -u postgres psql -c "CREATE USER grafana WITH PASSWORD 'secure-password';"
sudo -u postgres psql -c "CREATE DATABASE grafana WITH OWNER grafana;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE grafana TO grafana;"

sudo systemctl restart grafana-server
sudo systemctl enable grafana-server
```

### Fase 3: Deploy do BRLIA Control Plane

#### Configuracao do Ambiente

```bash
# Clonar repositorio
cd /opt
sudo git clone https://github.com/kauecodify/BRLIA-Sovereign-Cloud.git
sudo chown -R $USER:$USER /opt/BRLIA-Sovereign-Cloud
cd /opt/BRLIA-Sovereign-Cloud

# Criar arquivos de configuracao
mkdir -p /opt/brlia/config

# Configurar Control API
cat > /opt/brlia/config/control-api.yaml <<EOF
server:
  port: 8080
  mode: release
  shutdown_timeout: 30s

database:
  host: localhost
  port: 5432
  username: brlia
  password: secure-password
  database: brlia_control
  ssl_mode: disable
  pool_size: 50

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

auth:
  enabled: true
  jwt_secret: "brlia-sovereign-cloud-secret-change-me"
  jwt_expiry: 15m
  refresh_secret: "brlia-refresh-secret-change-me"
  refresh_expiry: 168h
  issuer: "brlia-control-api"
EOF'

# Configurar Waterfall
cat > /opt/brlia/config/waterfall.yaml <<EOF
server:
  port: 8081
  mode: release
  shutdown_timeout: 30s

database:
  host: localhost
  port: 5432
  username: brlia
  password: secure-password
  database: brlia_waterfall
  ssl_mode: disable
  pool_size: 50

redis:
  host: localhost
  port: 6379
  password: ""
  db: 1

auth:
  enabled: true
  jwt_secret: "brlia-waterfall-secret-change-me"
  jwt_expiry: 15m
  refresh_secret: "brlia-waterfall-refresh-secret-change-me"
  refresh_expiry: 168h
  issuer: "brlia-waterfall-api"

waterfall:
  worker_count: 20
  queue_worker_count: 10
  scheduler_interval: 1s
  capacity_check_interval: 5s
  failure_check_interval: 10s
  cost_optimizer_interval: 30s
  max_retries: 3
  retry_delay: 1s
  enable_ai: false
  ai_check_interval: 1m
EOF'
```

#### Build e Deploy

```bash
# Instalar Go
if ! command -v go &> /dev/null; then
    wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    source ~/.bashrc
fi

# Buildar Control API
cd /opt/BRLIA-Sovereign-Cloud/services/control-api
go mod download
go build -o /opt/brlia/bin/control-api main.go

# Buildar Waterfall
cd /opt/BRLIA-Sovereign-Cloud/services/waterfall
go mod download
go build -o /opt/brlia/bin/waterfall main.go

# Criar servicos systemd
sudo bash -c 'cat > /etc/systemd/system/brlia-control-api.service <<EOF
[Unit]
Description=BRLIA Control API
After=network.target postgresql.service redis.service

[Service]
User=root
WorkingDirectory=/opt/BRLIA-Sovereign-Cloud/services/control-api
ExecStart=/opt/brlia/bin/control-api
Environment=CONFIG_PATH=/opt/brlia/config/control-api.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF'

sudo bash -c 'cat > /etc/systemd/system/brlia-waterfall.service <<EOF
[Unit]
Description=BRLIA Waterfall Engine
After=network.target postgresql.service redis.service brlia-control-api.service

[Service]
User=root
WorkingDirectory=/opt/BRLIA-Sovereign-Cloud/services/waterfall
ExecStart=/opt/brlia/bin/waterfall
Environment=CONFIG_PATH=/opt/brlia/config/waterfall.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF'

# Inicializar banco de dados
# Aplicar migracoes usando goose, flyway ou similar

# Iniciar servicos
sudo systemctl daemon-reload
sudo systemctl start brlia-control-api
sudo systemctl start brlia-waterfall
sudo systemctl enable brlia-control-api
sudo systemctl enable brlia-waterfall

# Verificar status
sudo systemctl status brlia-control-api
sudo systemctl status brlia-waterfall
```

### Fase 4: Deploy do Admin Console

```bash
# Instalar Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# Buildar Admin Console
cd /opt/BRLIA-Sovereign-Cloud/apps/admin-console
npm install
npm run build

# Criar server.js
cat > /opt/BRLIA-Sovereign-Cloud/apps/admin-console/server.js <<EOF
const express = require('express');
const path = require('path');
const app = express();
const port = process.env.PORT || 3001;

app.use(express.static(path.join(__dirname, 'dist')));
app.get('*', (req, res) => {
    res.sendFile(path.join(__dirname, 'dist', 'index.html'));
});

app.listen(port, () => {
    console.log(`BRLIA Admin Console running on port ${port}`);
});
EOF'

# Criar servico systemd
sudo bash -c 'cat > /etc/systemd/system/brlia-admin-console.service <<EOF
[Unit]
Description=BRLIA Admin Console
After=network.target brlia-control-api.service brlia-waterfall.service

[Service]
User=root
WorkingDirectory=/opt/BRLIA-Sovereign-Cloud/apps/admin-console
ExecStart=/usr/bin/node /opt/BRLIA-Sovereign-Cloud/apps/admin-console/server.js
Environment=NODE_ENV=production
Environment=PORT=3001
Environment=API_URL=http://localhost:8080
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF'

npm run build

sudo systemctl daemon-reload
sudo systemctl start brlia-admin-console
sudo systemctl enable brlia-admin-console
```

### Fase 5: Configuracao de Seguranca

#### Certificados SSL

```bash
# Instalar certbot
sudo apt install -y certbot python3-certbot-nginx

# Obter certificados
sudo certbot certonly --standalone -d brlia.cloud.br

# Configurar renovacao automatica
sudo certbot renew --dry-run
```

#### Configuracao de Autenticacao

```bash
# Acessar Keycloak em http://localhost:8082
# Criar realm brlia
# Criar clientes: brlia-control-api, brlia-waterfall, brlia-admin-console
# Criar usuarios administradores
```

### Fase 6: Monitoramento

#### Configuracao do Prometheus

```bash
# Configurar alertas
cat > /etc/prometheus/brlia-rules.yml <<EOF
groups:
- name: brlia-alerts
  rules:
  - alert: ControlAPIDown
    expr: up{job="brlia-control"} == 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "Control API is down"
      description: "Control API has been down for more than 5 minutes"

  - alert: WaterfallEngineDown
    expr: up{job="brlia-waterfall"} == 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "Waterfall Engine is down"
      description: "Waterfall Engine has been down for more than 5 minutes"
EOF'

# Atualizar configuracao
sudo bash -c 'cat >> /etc/prometheus/prometheus.yml <<EOF

rule_files:
  - /etc/prometheus/brlia-rules.yml
EOF'

sudo systemctl restart prometheus
```

### Fase 7: Alta Disponibilidade

#### Kubernetes HA

```bash
# Inicializar cluster com multiplos control planes
sudo kubeadm init --control-plane-endpoint "192.168.1.100:6443" --upload-certs

# Obter token de join
kubeadm init phase upload-certs --upload-certs
kubeadm token create --print-join-command

# Adicionar outros control planes
# No segundo e terceiro node:
sudo kubeadm join 192.168.1.100:6443 --token <token> --discovery-token-ca-cert-hash <hash> --control-plane --certificate-key <cert-key>

# Adicionar worker nodes
sudo kubeadm join 192.168.1.100:6443 --token <token> --discovery-token-ca-cert-hash <hash>

# Configurar HAProxy
sudo apt install -y haproxy

cat > /etc/haproxy/haproxy.cfg <<EOF
global
    log /dev/log local0
    chroot /var/lib/haproxy
    user haproxy
    group haproxy
    daemon

defaults
    mode tcp
    timeout connect 5s
    timeout client 50s
    timeout server 50s

frontend kubernetes
    bind *:6443
    default_backend kubernetes-masters

backend kubernetes-masters
    balance roundrobin
    server master1 192.168.1.10:6443 check
    server master2 192.168.1.11:6443 check
    server master3 192.168.1.12:6443 check
EOF'

sudo systemctl restart haproxy
```

#### PostgreSQL HA

```bash
# Instalar Patroni
sudo apt install -y patroni

# Configurar Patroni
cat > /etc/patroni/patroni.yml <<EOF
scope: brlia
name: brlia-db-01

restapi:
  listen: 0.0.0.0:8008
  connect_address: 192.168.1.10:8008

etcd:
  hosts: 192.168.1.10:2379,192.168.1.11:2379,192.168.1.12:2379

bootstrap:
  dcs:
    ttl: 30
    loop_wait: 10
    retry_timeout: 10
    maximum_lag_on_failover: 1048576
    postgresql:
      use_pg_rewind: true
      use_slots: true
      parameters:
        max_connections: 500
        shared_buffers: 4GB

  initdb:
  - encoding: UTF8
  - data-checksums

  pg_hba:
  - host replication replicator 192.168.1.0/24 md5
  - host all all 192.168.1.0/24 md5

postgresql:
  listen: 0.0.0.0:5432
  connect_address: 192.168.1.10:5432
  data_dir: /var/lib/postgresql/15/main
  bin_dir: /usr/lib/postgresql/15/bin
  authentication:
    replication:
      username: replicator
      password: "replication-password"
    superuser:
      username: postgres
      password: "secure-password"
EOF'

sudo systemctl start patroni
sudo systemctl enable patroni
```

### Fase 8: Testes

#### Testes de Funcionalidade

```bash
# Testar Control API
curl -X GET http://localhost:8080/health
curl -X GET http://localhost:8080/ready

# Testar Waterfall Engine
curl -X GET http://localhost:8081/health
curl -X GET http://localhost:8081/ready

# Testar Admin Console
curl -X GET http://localhost:3001

# Testar autenticacao
curl -X POST http://localhost:8080/api/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin-password"}'

# Testar criacao de organizacao
curl -X POST http://localhost:8080/api/v1/organizations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Organization", "description": "Test"}'
```

#### Testes de Desempenho

```bash
# Instalar k6
sudo apt install -y k6

# Criar script de teste
cat > /opt/brlia/tests/load-test.js <<EOF
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 100,
  duration: '5m',
};

export default function () {
  const healthRes = http.get('http://localhost:8080/health');
  check(healthRes, {
    'health check passed': (r) => r.status === 200,
  });

  const orgRes = http.get('http://localhost:8080/api/v1/organizations', {
    headers: {
      'Authorization': 'Bearer <token>',
    },
  });
  check(orgRes, {
    'list organizations passed': (r) => r.status === 200,
  });

  sleep(1);
}
EOF'

# Executar teste
k6 run /opt/brlia/tests/load-test.js
```

### Fase 9: Operacao e Manutencao

#### Procedimentos Operacionais

```bash
# Iniciar servicos
sudo systemctl start brlia-control-api
sudo systemctl start brlia-waterfall
sudo systemctl start brlia-admin-console

# Parar servicos
sudo systemctl stop brlia-control-api
sudo systemctl stop brlia-waterfall
sudo systemctl stop brlia-admin-console

# Reiniciar servicos
sudo systemctl restart brlia-control-api
sudo systemctl restart brlia-waterfall
sudo systemctl restart brlia-admin-console

# Verificar status
sudo systemctl status brlia-control-api
sudo systemctl status brlia-waterfall
sudo systemctl status brlia-admin-console

# Verificar logs
sudo journalctl -u brlia-control-api -f
sudo journalctl -u brlia-waterfall -f
sudo journalctl -u brlia-admin-console -f
```

#### Atualizacao do Sistema

```bash
# Atualizar codigo
cd /opt/BRLIA-Sovereign-Cloud
sudo git pull origin main

# Rebuildar
cd /opt/BRLIA-Sovereign-Cloud/services/control-api
go build -o /opt/brlia/bin/control-api main.go

cd /opt/BRLIA-Sovereign-Cloud/services/waterfall
go build -o /opt/brlia/bin/waterfall main.go

# Reiniciar servicos
sudo systemctl restart brlia-control-api
sudo systemctl restart brlia-waterfall
```

#### Backup e Recuperacao

```bash
# Criar script de backup
cat > /opt/brlia/scripts/backup.sh <<EOF
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR=/opt/brlia/backup/$DATE

mkdir -p $BACKUP_DIR/config
cp -r /opt/brlia/config $BACKUP_DIR/config/

mkdir -p $BACKUP_DIR/database
sudo -u postgres pg_dump brlia_control > $BACKUP_DIR/database/brlia_control.sql
sudo -u postgres pg_dump brlia_waterfall > $BACKUP_DIR/database/brlia_waterfall.sql

cd /opt/brlia/backup
tar -czf $DATE.tar.gz $DATE
rm -rf $DATE
EOF

chmod +x /opt/brlia/scripts/backup.sh

# Agendar backup diario
sudo bash -c 'cat > /etc/cron.d/brlia-backup <<EOF
0 2 * * * root /opt/brlia/scripts/backup.sh
EOF'

# Restaurar de backup
cd /opt/brlia/backup
LATEST=$(ls -t *.tar.gz | head -1)
tar -xzf $LATEST
sudo -u postgres psql brlia_control < database/brlia_control.sql
sudo -u postgres psql brlia_waterfall < database/brlia_waterfall.sql
sudo systemctl restart brlia-control-api
sudo systemctl restart brlia-waterfall
```

## Conformidade e Seguranca

### Requisitos de Conformidade

1. **LGPD**: Todos os dados de usuarios brasileiros sao armazenados no Brasil, criptografia em repouso e transito, controle de acesso e auditoria
2. **Marco Civil da Internet**: Logs de conexao por 1 ano, logs de aplicacao por 6 meses, neutralidade da rede
3. **Normas do Banco Central**: MFA, trilha de auditoria imutavel, alta disponibilidade (99.95%), recuperacao de desastres

### Politicas de Seguranca

1. **Controle de Acesso**: Zero Trust, MFA, RBAC, rotacao periodica de credenciais
2. **Criptografia**: TLS 1.3, AES-256, KMS, HSM
3. **Monitoramento**: SIEM, IDS/IPS, WAF, DDoS protection
4. **Auditoria**: Logs imutaveis, trilha de auditoria completa, relatorios periodicos

## Suporte

- **Suporte Tecnico**: suporte@brlia.cloud.br
- **Incidentes Criticos**: +55 11 99999-9999
- **Documentacao**: https://docs.brlia.cloud.br

## Historico de Revisoes

| Versao | Data | Autor | Descricao |
|--------|------|-------|-----------|
| 1.0 | 2026-11-21 | BRLIA Team | Versao inicial |
