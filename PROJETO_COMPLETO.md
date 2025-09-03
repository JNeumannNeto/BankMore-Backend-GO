# BankMore Backend GO - Projeto Completo

## 📋 Resumo do Projeto

Este projeto é uma implementação completa do sistema bancário BankMore em Go, baseado no projeto original em C#. O sistema implementa uma arquitetura de microserviços com comunicação assíncrona via Kafka.

## 🏗️ Arquitetura Implementada

### Microserviços
1. **Account API** (Porta 8001) - Gerenciamento de contas e autenticação
2. **Transfer API** (Porta 8002) - Processamento de transferências
3. **Fee API** (Porta 8003) - Processamento de tarifas

### Tecnologias Utilizadas
- **Go 1.21+** - Linguagem principal
- **Gin** - Framework web HTTP
- **GORM** - ORM para acesso aos dados
- **SQLite** - Banco de dados
- **Kafka** - Mensageria assíncrona
- **JWT** - Autenticação e autorização
- **Docker** - Containerização
- **Swagger** - Documentação das APIs

## 📁 Estrutura do Projeto

```
BankMore-Backend-GO/
├── cmd/                          # Aplicações principais
│   ├── account-api/             # Account API
│   ├── transfer-api/            # Transfer API
│   └── fee-api/                 # Fee API
├── internal/                    # Código interno
│   ├── shared/                  # Código compartilhado
│   │   ├── models/             # Modelos compartilhados
│   │   ├── middleware/         # Middlewares (JWT, CORS)
│   │   ├── utils/              # Utilitários (CPF, Hash)
│   │   └── kafka/              # Cliente Kafka
│   ├── account/                # Domínio de Contas
│   ├── transfer/               # Domínio de Transferências
│   └── fee/                    # Domínio de Tarifas
├── database/                   # Scripts de banco
├── deployments/               # Configurações Docker
├── scripts/                   # Scripts de build
├── docs/                      # Documentação
└── tests/                     # Testes (a implementar)
```

## 🔧 Funcionalidades Implementadas

### Account API
- ✅ Cadastro de contas
- ✅ Autenticação com JWT
- ✅ Movimentações (depósitos e saques)
- ✅ Consulta de saldo
- ✅ Inativação de contas
- ✅ Validação de CPF
- ✅ Hash de senhas com salt
- ✅ Controle de idempotência

### Transfer API
- ✅ Transferências entre contas
- ✅ Validações de saldo
- ✅ Rollback em caso de erro
- ✅ Publicação de eventos Kafka
- ✅ Comunicação com Account API

### Fee API
- ✅ Processamento de tarifas
- ✅ Consumo de eventos Kafka
- ✅ Débito automático de tarifas
- ✅ Consulta de tarifas por conta

### Recursos Compartilhados
- ✅ Middleware JWT
- ✅ Validação de CPF
- ✅ Utilitários de hash
- ✅ Cliente Kafka (Producer/Consumer)
- ✅ Modelos de erro padronizados
- ✅ Health checks

## 🚀 Como Executar

### Opção 1: Docker (Recomendado)
```bash
# Executar todo o sistema
docker-compose -f deployments/docker-compose.yml up --build

# Ou usando Makefile
make docker-up
```

### Opção 2: Execução Local
```bash
# Instalar dependências
go mod tidy

# Build
make build

# Executar serviços (em terminais separados)
make run-account
make run-transfer
make run-fee
```

## 📊 Endpoints Principais

### Account API (8001)
- `POST /api/account/register` - Cadastrar conta
- `POST /api/account/login` - Login
- `POST /api/account/movement` - Movimentação
- `GET /api/account/balance` - Consultar saldo
- `PUT /api/account/deactivate` - Inativar conta

### Transfer API (8002)
- `POST /api/transfer` - Realizar transferência

### Fee API (8003)
- `GET /api/fee/{accountNumber}` - Consultar tarifas
- `GET /api/fee/fee/{id}` - Consultar tarifa específica

## 🔒 Segurança

- **JWT Authentication** em todos os endpoints protegidos
- **Validação de CPF** com dígitos verificadores
- **Hash de senhas** com salt único
- **Validações de entrada** em todos os endpoints
- **CORS** configurado para desenvolvimento

## 📈 Diferenças do Projeto Original C#

### Adaptações para Go
1. **Estrutura de pastas** seguindo convenções Go
2. **Gin** ao invés de ASP.NET Core
3. **GORM** ao invés de Dapper
4. **Sarama** ao invés de Confluent.Kafka
5. **Logrus** para logging estruturado
6. **Graceful shutdown** nativo

### Melhorias Implementadas
1. **Context propagation** para cancelamento
2. **Health checks** nativos
3. **Makefile** para automação
4. **Docker multi-stage builds**
5. **Structured logging** com campos contextuais

## 🧪 Testes

### Estrutura Preparada
- Framework de testes nativo do Go
- Mocks preparados para dependências
- Testes de integração estruturados

### Executar Testes
```bash
make test
make test-coverage
```

## 🐳 Docker

### Serviços Configurados
- **zookeeper** - Coordenação Kafka
- **kafka** - Broker de mensagens
- **sqlite-db** - Banco de dados
- **account-api** - API de contas
- **transfer-api** - API de transferências
- **fee-api** - API de tarifas

### Comandos Docker
```bash
make docker-up      # Iniciar serviços
make docker-down    # Parar serviços
make docker-logs    # Ver logs
```

## 🔧 Desenvolvimento

### Comandos Úteis
```bash
make help           # Ver todos os comandos
make dev-setup      # Configurar ambiente
make fmt            # Formatar código
make lint           # Lint código
make swagger        # Gerar documentação
make health-check   # Verificar saúde dos serviços
```

### Variáveis de Ambiente
Ver arquivo `.env.example` para configurações disponíveis.

## 📚 Documentação

- **Swagger UI** disponível em cada API
- **README.md** com visão geral
- **INSTRUCOES_EXECUCAO.md** com guia detalhado
- **Makefile** com comandos automatizados

## 🎯 Próximos Passos

### Melhorias Futuras
- [ ] Implementar testes unitários completos
- [ ] Adicionar métricas com Prometheus
- [ ] Implementar circuit breaker
- [ ] Adicionar rate limiting
- [ ] Logs centralizados
- [ ] Monitoramento com Grafana
- [ ] Testes de carga
- [ ] Criptografia de dados sensíveis

### Otimizações
- [ ] Cache distribuído com Redis
- [ ] Connection pooling otimizado
- [ ] Compressão de responses
- [ ] Paginação em consultas
- [ ] Índices de banco otimizados

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature
3. Implemente os testes
4. Execute `make fmt` e `make lint`
5. Commit suas mudanças
6. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT.

---

**Projeto desenvolvido com ❤️ em Go, seguindo as melhores práticas de arquitetura de microserviços**
