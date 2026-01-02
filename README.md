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

- **Core Service**: Recebe notificações de criação de pedidos
- **Operation Service**: Notifica mudanças de status de pagamento
- **Mercado Pago**: Processamento de pagamentos via QR Code

## 🔗 Integração Serverless (AWS Lambda)

✅ **PRONTO PARA USO**: A autenticação serverless já está totalmente configurada!

### **🛠️ Código Implementado**
O código foi atualizado seguindo o padrão do monolítico `tc-golunch-api`:

1. **ServerlessAuthGateway**: Implementado para comunicação com Lambda
2. **ServerlessAuthMiddleware**: Middleware de autenticação serverless
3. **main.go**: Atualizado para usar serverless auth em vez de JWT local

### **🔧 Configuração das URLs**

**⚠️ PREREQUISITO**: Primeiro faça deploy do `tc-golunch-serverless` para gerar as URLs reais!

```bash
# 1. Deploy serverless (OBRIGATÓRIO primeiro)
cd ../tc-golunch-serverless
terraform init
terraform apply
# Isso cria funções Lambda e gera URLs reais do API Gateway

# 2. Obter URLs reais geradas
terraform output
# Output: api_gateway_url = "https://abc123def.execute-api.us-east-1.amazonaws.com"

# 3. ENTÃO configurar variáveis locais com URLs reais:
export LAMBDA_AUTH_URL="https://abc123def.execute-api.us-east-1.amazonaws.com/auth"
export SERVICE_AUTH_LAMBDA_URL="https://abc123def.execute-api.us-east-1.amazonaws.com/service-auth"

# Variáveis existentes (mantidas)
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DATABASE="golunch_payments"
export PAYMENT_SERVICE_PORT="8082"
export ORDER_SERVICE_URL="http://localhost:8081"
export OPERATION_SERVICE_URL="http://localhost:8083"

# Mercado Pago (necessárias)
export MP_ACCESS_TOKEN="seu-mercado-pago-token"
export MP_USER_ID="seu-user-id"
export MP_POS_ID="seu-pos-id"
```

### **📦 Deploy Kubernetes**

⚠️ **PREREQUISITO**: Deploy do `tc-golunch-serverless` ANTES de fazer deploy Kubernetes!

**Passo-a-passo completo:**

```bash
# PASSO 1: Deploy Serverless (OBRIGATÓRIO primeiro)
cd ../tc-golunch-serverless
terraform init
terraform apply

# PASSO 2: Obter URLs reais do API Gateway
terraform output
# Exemplo output: api_gateway_url = "https://abc123def.execute-api.us-east-1.amazonaws.com"

# PASSO 3: Atualizar ConfigMap com URLs REAIS
cd ../tc-golunch-payment-service
vim k8s/payment-service-configmap.yaml

# SUBSTITUIR estas linhas (são templates):
# LAMBDA_AUTH_URL: "https://your-api-gateway-id.execute-api.region.amazonaws.com/auth"
# SERVICE_AUTH_LAMBDA_URL: "https://your-api-gateway-id.execute-api.region.amazonaws.com/service-auth"

# POR URLs reais obtidas no terraform output:
# LAMBDA_AUTH_URL: "https://abc123def.execute-api.us-east-1.amazonaws.com/auth"
# SERVICE_AUTH_LAMBDA_URL: "https://abc123def.execute-api.us-east-1.amazonaws.com/service-auth"

# PASSO 4: Deploy Kubernetes
kubectl apply -f k8s/
```

**Estrutura já configurada:**
```yaml
# k8s/payment-service-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: payment-service-config
data:
  LAMBDA_AUTH_URL: "https://your-api-gateway-id.execute-api.region.amazonaws.com/auth"
  SERVICE_AUTH_LAMBDA_URL: "https://your-api-gateway-id.execute-api.region.amazonaws.com/service-auth"
  # ... outras variáveis
```

### **✅ Verificação da Configuração**

Após configurar as variáveis, teste a integração:

```bash
# 1. Inicie o serviço
go run cmd/api/main.go

# 2. Teste health check
curl -X GET http://localhost:8082/ping

# 3. Teste endpoint protegido (requer autenticação via Lambda)
curl -X POST http://localhost:8082/payment \
  -H "Authorization: Bearer <token-do-lambda>" \
  -H "Content-Type: application/json" \
  -d '{"order_id": "123", "amount": 50.00}'
```

### **🔄 Migração Gradual**

A implementação mantém **compatibilidade total** com o código existente:
- ✅ Mesmas interfaces de autenticação
- ✅ Mesmos endpoints e responses  
- ✅ Zero breaking changes para clientes
- ✅ Fallback automático se Lambda não disponível

