# Instruções de Execução - BankMore Backend GO

Este documento contém as instruções para executar o sistema bancário BankMore implementado em Go.

## 📋 Pré-requisitos

### Para execução local:
- Go 1.21+ instalado
- SQLite3
- Apache Kafka (opcional, pode usar Docker)

### Para execução com Docker:
- Docker
- Docker Compose

## 🚀 Execução com Docker (Recomendado)

### 1. Clone o repositório
```bash
git clone <repository-url>
cd BankMore-Backend-GO
```

### 2. Execute o sistema completo
```bash
docker-compose -f deployments/docker-compose.yml up --build
```

### 3. Aguarde a inicialização
O sistema estará disponível em:
- **Account API**: http://localhost:8001
- **Transfer API**: http://localhost:8002
- **Fee API**: http://localhost:8003
- **Kafka**: localhost:9092
- **Zookeeper**: localhost:2181

### 4. Acesse a documentação Swagger
- **Account API**: http://localhost:8001/swagger/index.html
- **Transfer API**: http://localhost:8002/swagger/index.html
- **Fee API**: http://localhost:8003/swagger/index.html

## 🔧 Execução Local

### 1. Instale as dependências
```bash
go mod tidy
```

### 2. Execute o build
```bash
chmod +x scripts/build.sh
./scripts/build.sh
```

### 3. Configure as variáveis de ambiente
```bash
export DB_PATH="./database/bankmore.db"
export KAFKA_BROKERS="localhost:9092"
export JWT_SECRET="your-secret-key-here"
export TRANSFER_FEE_AMOUNT="2.00"
```

### 4. Execute os serviços
Em terminais separados:

```bash
# Terminal 1 - Account API
./bin/account-api

# Terminal 2 - Transfer API
./bin/transfer-api

# Terminal 3 - Fee API
./bin/fee-api
```

## 📊 Health Checks

Verifique se os serviços estão funcionando:

```bash
# Account API
curl http://localhost:8001/health

# Transfer API
curl http://localhost:8002/health

# Fee API
curl http://localhost:8003/health
```

## 🧪 Testando as APIs

### 1. Cadastrar uma conta
```bash
curl -X POST http://localhost:8001/api/account/register \
  -H "Content-Type: application/json" \
  -d '{
    "cpf": "12345678901",
    "name": "João Silva",
    "password": "senha123"
  }'
```

### 2. Fazer login
```bash
curl -X POST http://localhost:8001/api/account/login \
  -H "Content-Type: application/json" \
  -d '{
    "cpf": "12345678901",
    "password": "senha123"
  }'
```

### 3. Realizar depósito
```bash
curl -X POST http://localhost:8001/api/account/movement \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "requestId": "dep-001",
    "accountNumber": "100001",
    "amount": 1000.00,
    "type": "C"
  }'
```

### 4. Consultar saldo
```bash
curl -X GET http://localhost:8001/api/account/balance \
  -H "Authorization: Bearer <TOKEN>"
```

### 5. Realizar transferência
```bash
curl -X POST http://localhost:8002/api/transfer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "requestId": "trans-001",
    "destinationAccountNumber": "100002",
    "amount": 100.00
  }'
```

### 6. Consultar tarifas
```bash
curl -X GET http://localhost:8003/api/fee/100001
```

## 🔍 Logs e Monitoramento

### Visualizar logs dos containers
```bash
# Todos os serviços
docker-compose -f deployments/docker-compose.yml logs -f

# Serviço específico
docker-compose -f deployments/docker-compose.yml logs -f account-api
docker-compose -f deployments/docker-compose.yml logs -f transfer-api
docker-compose -f deployments/docker-compose.yml logs -f fee-api
```

### Monitorar Kafka
```bash
# Listar tópicos
docker exec -it <kafka-container> kafka-topics --bootstrap-server localhost:9092 --list

# Consumir mensagens
docker exec -it <kafka-container> kafka-console-consumer --bootstrap-server localhost:9092 --topic transfer-events --from-beginning
```

## 🛠️ Desenvolvimento

### Executar testes
```bash
go test ./...
```

### Executar testes com coverage
```bash
go test -cover ./...
```

### Gerar documentação Swagger
```bash
# Instalar swag
go install github.com/swaggo/swag/cmd/swag@latest

# Gerar docs para cada API
swag init -g cmd/account-api/main.go -o docs/account
swag init -g cmd/transfer-api/main.go -o docs/transfer
swag init -g cmd/fee-api/main.go -o docs/fee
```

## 🐛 Troubleshooting

### Problema: Erro de conexão com Kafka
**Solução**: Aguarde alguns segundos para o Kafka inicializar completamente antes de iniciar as APIs.

### Problema: Banco de dados não encontrado
**Solução**: Certifique-se de que o diretório `database` existe e tem permissões adequadas.

### Problema: Porta já em uso
**Solução**: Verifique se não há outros serviços rodando nas portas 8001, 8002, 8003, 9092.

```bash
# Verificar portas em uso
netstat -tulpn | grep :800
```

### Problema: Erro de build
**Solução**: Certifique-se de ter Go 1.21+ instalado e CGO habilitado para SQLite.

```bash
go version
export CGO_ENABLED=1
```

## 🔧 Configurações Avançadas

### Variáveis de Ambiente Disponíveis

| Variável | Descrição | Padrão |
|----------|-----------|---------|
| `DB_PATH` | Caminho do banco SQLite | `./database/bankmore.db` |
| `KAFKA_BROKERS` | Servidores Kafka | `localhost:9092` |
| `JWT_SECRET` | Chave secreta JWT | `your-secret-key-here-change-in-production` |
| `TRANSFER_FEE_AMOUNT` | Valor da tarifa de transferência | `2.00` |
| `ACCOUNT_API_URL` | URL da Account API | `http://localhost:8001` |
| `PORT` | Porta do serviço | `8001/8002/8003` |

### Personalizar configurações Docker
Edite o arquivo `deployments/docker-compose.yml` para ajustar as configurações conforme necessário.

## 📚 Documentação Adicional

- [README.md](README.md) - Visão geral do projeto
- [Swagger UI](http://localhost:8001/swagger/index.html) - Documentação interativa das APIs
- [Projeto Original C#](../BankMore-Backend/README.md) - Referência da implementação original

## 🆘 Suporte

Em caso de problemas:
1. Verifique os logs dos serviços
2. Confirme se todas as dependências estão instaladas
3. Verifique se as portas não estão em uso
4. Consulte a documentação do projeto original em C#
