# 💳 GoLunch Payment Service

Microsserviço responsável pelo processamento de pagamentos da lanchonete GoLunch. Este serviço gerencia integrações com processadores de pagamento, processa webhooks e mantém o status dos pagamentos.

## 🎯 Responsabilidades

- **Processamento de Pagamentos**: Criação e gerenciamento de pagamentos via QR Code
- **Integração Mercado Pago**: Geração de QR Codes e processamento de pagamentos
- **Webhooks**: Recebimento e processamento de notificações de pagamento
- **Status de Pagamento**: Controle do fluxo de status dos pagamentos
- **Histórico de Transações**: Armazenamento e consulta de transações

## 🏗️ Arquitetura

O serviço segue os princípios da **Arquitetura Hexagonal** com as seguintes camadas:

- **Entities**: Regras de negócio fundamentais
- **Use Cases**: Lógica de negócio específica
- **Gateways**: Interfaces para acesso a dados externos
- **Controllers**: Coordenação entre camadas
- **Handlers**: Gerenciamento de requisições HTTP
- **External/Infrastructure**: Implementações concretas (MongoDB, APIs externas)

## 🗄️ Banco de Dados

- **MongoDB**: Banco de dados NoSQL para flexibilidade de dados de pagamento
- **Coleções**:
  - `payments`: Dados dos pagamentos
  - `transactions`: Histórico de transações
  - `webhooks`: Log de webhooks recebidos

## 🚀 Endpoints Disponíveis

### Pagamentos
- `POST /payment` - Criar novo pagamento
- `GET /payment/:id` - Consultar pagamento por ID

### Webhooks
- `POST /webhook/payment/check` - Webhook do Mercado Pago

### Health Check
- `GET /ping` - Health check do serviço

## 🔧 Configuração Local

1. **Clone o repositório**
2. **Configure as variáveis de ambiente**:
   ```bash
   export MONGODB_URI="mongodb://localhost:27017"
   export MONGODB_DATABASE="golunch_payments"
   export MERCADO_PAGO_ACCESS_TOKEN="your_access_token"
   export MERCADO_PAGO_SELLER_APP_USER_ID="your_seller_id"
   ```

3. **Execute o banco de dados**:
   ```bash
   docker run -d -p 27017:27017 --name mongodb mongo:latest
   ```

4. **Execute a aplicação**:
   ```bash
   go run cmd/api/main.go
   ```

## 📋 Dependências

- **Go** 1.24.3
- **MongoDB** 7.0+
- **Gin** - Framework web
- **MongoDB Driver** - Driver para MongoDB
- **Resty** - Cliente HTTP para APIs externas
- **Swagger** - Documentação da API

## 🧪 Testes

```bash
# Executar todos os testes
go test ./...

# Executar testes com cobertura
go test -cover ./...

# Executar testes BDD
go test -tags=bdd ./...
```

## 📊 Cobertura de Testes

- **Meta**: 80% de cobertura
- **BDD**: Implementado para cenários de processamento de pagamento
- **Testes Unitários**: Todos os use cases e controllers
- **Testes de Integração**: Webhooks e APIs externas

## 🐳 Docker

```bash
# Build da imagem
docker build -t tc-golunch-payment-service .

# Executar container
docker run -p 8082:8082 tc-golunch-payment-service
```

## 📈 Monitoramento

- **Health Check**: `GET /ping`
- **Swagger UI**: `GET /swagger/index.html`
- **Logs**: Estruturados em JSON
- **Métricas**: Tempo de resposta, taxa de sucesso

## 🔄 CI/CD

O serviço possui pipeline CI/CD configurado com:
- Validação de código
- Execução de testes
- Análise de cobertura
- Build e deploy automático
- Proteção de branch main

## 🔐 Segurança

- **Tokens de Acesso**: Armazenados como secrets
- **Validação de Webhooks**: Verificação de assinatura
- **HTTPS**: Comunicação segura
- **Rate Limiting**: Proteção contra abuso

## 📝 Documentação da API

A documentação completa da API está disponível via Swagger UI em:
`http://localhost:8082/swagger/index.html`

## 🔗 Integração com Outros Serviços

- **Order Service**: Recebe notificações de criação de pedidos
- **Production Service**: Notifica mudanças de status de pagamento
- **Mercado Pago**: Processamento de pagamentos via QR Code

