# BankMore - Sistema Bancário com Microserviços em Go

Sistema bancário completo implementado com arquitetura de microserviços em Go, baseado no projeto original em C#.

## 🏗️ Arquitetura

```
BankMore-Backend-GO - Sistema Bancário
│
├── 📁 cmd/
│   ├── account-api/                  # API de Contas (Porta 8001)
│   ├── transfer-api/                 # API de Transferências (Porta 8002)
│   └── fee-api/                      # API de Tarifas (Porta 8003)
│
├── 📁 internal/
│   ├── shared/                       # Código compartilhado
│   │   ├── models/                   # Modelos compartilhados
│   │   ├── middleware/               # Middlewares (JWT, CORS)
│   │   ├── utils/                    # Utilitários (CPF, Hash)
│   │   └── kafka/                    # Cliente Kafka
│   │
│   ├── account/                      # Domínio de Contas
│   │   ├── domain/                   # Entidades de domínio
│   │   ├── handlers/                 # Handlers HTTP
│   │   ├── repository/               # Repositórios
│   │   └── service/                  # Serviços de negócio
│   │
│   ├── transfer/                     # Domínio de Transferências
│   │   ├── domain/                   # Entidades de domínio
│   │   ├── handlers/                 # Handlers HTTP
│   │   ├── repository/               # Repositórios
│   │   └── service/                  # Serviços de negócio
│   │
│   └── fee/                          # Domínio de Tarifas
│       ├── domain/                   # Entidades de domínio
│       ├── handlers/                 # Handlers HTTP
│       ├── repository/               # Repositórios
│       └── service/                  # Serviços de negócio
│
├── 📁 database/
│   └── init.sql                      # Script de inicialização do banco
│
├── 📁 deployments/
│   └── docker-compose.yml            # Orquestração dos serviços
│
├── 📁 scripts/
│   └── build.sh                      # Scripts de build
│
├── 📄 go.mod                         # Dependências Go
├── 📄 go.sum                         # Checksums das dependências
└── 📖 README.md                      # Documentação
```

### Fluxo de Comunicação:
```
[Cliente] 
    ↓ HTTP/JWT
[Account API] ←→ [SQLite Database]
    ↓ Kafka (Transfer Events)
[Transfer API] ←→ [SQLite Database]
    ↓ Kafka (Fee Events)
[Fee API] ←→ [SQLite Database]
```

## 🛠️ Tecnologias Utilizadas

- **Go 1.21+**: Linguagem principal
- **Gin**: Framework web HTTP
- **SQLite**: Banco de dados
- **GORM**: ORM para acesso aos dados
- **JWT-Go**: Autenticação e autorização
- **Sarama**: Cliente Kafka para Go
- **Docker**: Containerização
- **Swagger**: Documentação das APIs (gin-swagger)
- **Testify**: Framework de testes

## 🔧 Funcionalidades Implementadas

### Requisitos Funcionais ✅
- [x] Cadastro e autenticação de usuários
- [x] Realização de movimentações (depósitos e saques)
- [x] Transferências entre contas
- [x] Consulta de saldo
- [x] Sistema de tarifas

### Requisitos Técnicos ✅
- [x] **Clean Architecture**: Estrutura de domínio bem definida
- [x] **JWT**: Autenticação em todos os endpoints
- [x] **Idempotência**: Prevenção de operações duplicadas
- [x] **Kafka**: Comunicação assíncrona
- [x] **Cache**: Implementado com cache em memória
- [x] **Validações**: CPF, senhas, valores, etc.
- [x] **Docker**: Containerização completa

## 🚀 Como Executar

### Pré-requisitos
- Go 1.21+ instalado
- Docker e Docker Compose instalados

### Build do Projeto

1. **Clone o repositório**
```bash
git clone <repository-url>
cd BankMore-Backend-GO
```

2. **Instale as dependências**
```bash
go mod tidy
```

3. **Execute o build**
```bash
./scripts/build.sh
```

### Executando com Docker Compose

1. **Execute o sistema completo**
```bash
docker-compose -f deployments/docker-compose.yml up --build
```

2. **Aguarde a inicialização**
- Kafka: http://localhost:9092
- Account API: http://localhost:8001
- Transfer API: http://localhost:8002
- Fee API: http://localhost:8003

### Acessando as APIs

- **Account API Swagger**: http://localhost:8001/swagger/index.html
- **Transfer API Swagger**: http://localhost:8002/swagger/index.html
- **Fee API Swagger**: http://localhost:8003/swagger/index.html

## 📋 Endpoints Principais

### Account API (Porta 8001)

#### POST `/api/account/register`
Cadastra uma nova conta corrente
```json
{
  "cpf": "12345678901",
  "name": "João Silva",
  "password": "senha123"
}
```

#### POST `/api/account/login`
Realiza login e retorna JWT token
```json
{
  "cpf": "12345678901",
  "password": "senha123"
}
```

#### POST `/api/account/movement`
Realiza movimentação na conta (requer autenticação)
```json
{
  "requestId": "uuid-unique",
  "accountNumber": "123456",
  "amount": 100.00,
  "type": "C"
}
```

#### GET `/api/account/balance`
Consulta saldo da conta (requer autenticação)

### Transfer API (Porta 8002)

#### POST `/api/transfer`
Realiza transferência entre contas (requer autenticação)
```json
{
  "requestId": "uuid-unique",
  "destinationAccountNumber": "654321",
  "amount": 50.00
}
```

### Fee API (Porta 8003)

#### GET `/api/fee/{accountNumber}`
Consulta tarifas por número da conta

#### GET `/api/fee/fee/{id}`
Consulta tarifa específica por ID

## 🗄️ Estrutura do Banco de Dados

### Tabelas Principais

- **contacorrente**: Dados das contas
- **movimento**: Movimentações financeiras
- **transferencia**: Histórico de transferências
- **tarifa**: Registro de tarifas cobradas
- **idempotencia**: Controle de idempotência

## 🔒 Segurança

### Autenticação JWT
- Todos os endpoints protegidos requerem token JWT
- Token contém informações da conta logada
- Validação de expiração e assinatura

### Validações Implementadas
- **CPF**: Validação completa com dígitos verificadores
- **Senhas**: Hash com salt único por usuário
- **Valores**: Apenas valores positivos
- **Contas**: Verificação de existência e status ativo

## 🔄 Fluxo de Transferência

1. **Validações iniciais** (conta origem, destino, valor)
2. **Verificação de saldo** na conta origem
3. **Débito na conta origem** via Account API
4. **Crédito na conta destino** via Account API
5. **Registro da transferência** no banco de dados
6. **Publicação no Kafka** para cobrança de tarifa
7. **Fee API**: Processa tarifa e debita automaticamente

## 📊 Monitoramento e Logs

- Logs estruturados em todos os serviços
- Rastreamento de operações via RequestId
- Métricas de performance disponíveis

## 🧪 Testes

### Executando Testes
```bash
go test ./...
```

### Testes com Coverage
```bash
go test -cover ./...
```

## 🐳 Docker

### Serviços no Docker Compose
- **zookeeper**: Coordenação do Kafka
- **kafka**: Broker de mensagens
- **sqlite-db**: Banco de dados compartilhado
- **account-api**: API de contas
- **transfer-api**: API de transferências
- **fee-api**: API de tarifas

## 🔧 Configurações

### Variáveis de Ambiente
- `DB_PATH`: Caminho do banco SQLite
- `KAFKA_BROKERS`: Servidores Kafka
- `JWT_SECRET`: Chave secreta JWT
- `TRANSFER_FEE_AMOUNT`: Valor da tarifa

## 📈 Diferenças do Projeto Original C#

### Adaptações para Go
- **Gin** ao invés de ASP.NET Core
- **GORM** ao invés de Dapper
- **Sarama** ao invés de Confluent.Kafka
- **Testify** ao invés de xUnit
- Estrutura de pastas seguindo convenções Go

### Melhorias Implementadas
- Graceful shutdown em todos os serviços
- Context propagation para cancelamento
- Structured logging com logrus
- Health checks nativos

---

**Desenvolvido com ❤️ em Go seguindo as melhores práticas**
