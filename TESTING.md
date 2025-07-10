# Testes da Aplicação

Este documento descreve como executar os testes da aplicação de temperatura por CEP.

## Estrutura dos Testes

Os testes estão organizados por pacote e funcionalidade:

### 1. Testes de Utilitários (`pkg/utils/`)
- **Arquivo**: `pkg/utils/temperatures_test.go`
- **Funções testadas**:
  - `ConvertFahrenheitToCelsius()`
  - `ConvertCelsiusToFahrenheit()`
  - `ConvertCelsiusToKelvin()`
  - `ConvertKelvinToCelsius()`
  - `FormatTemperatures()`

### 2. Testes de ViaCEP (`pkg/viacep/`)
- **Arquivo**: `pkg/viacep/viacep_test.go`
- **Funções testadas**:
  - `FetchCEPData()`
  - Validação de CEP
  - Formatação de CEP
  - Unmarshal de JSON

### 3. Testes de Weather Search (`pkg/weather/`)
- **Arquivo**: `pkg/weather/search_test.go`
- **Funções testadas**:
  - `FecthSearchFromWeatherAPI()`
  - Busca de cidades
  - Tratamento de respostas vazias

### 4. Testes de Weather Data (`pkg/weather/`)
- **Arquivo**: `pkg/weather/weather_test.go`
- **Funções testadas**:
  - `FetchWeatherData()`
  - Estruturas de resposta
  - Coordenadas inválidas

### 5. Testes de Service (`pkg/weather/`)
- **Arquivo**: `pkg/weather/service_test.go`
- **Funções testadas**:
  - `GetTemperatureByCEP()`
  - `TemperatureHandler()`
  - Fluxo completo de CEP para temperatura

### 6. Testes de Integração (`main/`)
- **Arquivo**: `main_test.go`
- **Testes**:
  - Inicialização do servidor
  - Fluxo completo da aplicação
  - Tratamento de erros HTTP
  - Métodos HTTP
  - Headers de resposta

## Como Executar os Testes

### Executar todos os testes
```bash
go test ./...
```

### Executar testes de um pacote específico
```bash
# Testes de utilitários
go test ./pkg/utils

# Testes de ViaCEP
go test ./pkg/viacep

# Testes de weather
go test ./pkg/weather

# Testes de integração
go test ./main
```

### Executar testes com verbose
```bash
go test -v ./...
```

### Executar testes com cobertura
```bash
go test -cover ./...
```

### Executar testes com relatório de cobertura detalhado
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Configuração para Testes

### Variáveis de Ambiente
Para executar os testes que fazem chamadas reais para as APIs, você precisa configurar:

1. **Weather API Key**:
   ```bash
   export WEATHER_API_KEY="sua_chave_aqui"
   ```

2. **Arquivo .env**:
   ```
   WEATHER_API_KEY=sua_chave_aqui
   ```

### Testes que Requerem API Key
Alguns testes fazem chamadas reais para as APIs externas. Estes testes:
- São marcados com `t.Skipf()` se a API key não estiver disponível
- Podem falhar se as APIs estiverem indisponíveis
- São úteis para validação de integração

### Testes Unitários
A maioria dos testes são unitários e não requerem:
- API keys
- Conexão com internet
- Serviços externos

## Tipos de Teste

### 1. Testes Unitários
- Testam funções individuais
- Usam mocks quando necessário
- Não dependem de serviços externos

### 2. Testes de Integração
- Testam o fluxo completo da aplicação
- Fazem chamadas reais para APIs (quando possível)
- Validam a integração entre componentes

### 3. Testes de HTTP Handler
- Testam os endpoints HTTP
- Validam status codes e headers
- Testam diferentes métodos HTTP

### 4. Testes de JSON
- Validam marshaling/unmarshaling
- Testam estruturas de dados
- Verificam compatibilidade com APIs

## Exemplos de Execução

### Executar apenas testes unitários
```bash
go test -short ./...
```

### Executar testes com timeout
```bash
go test -timeout 30s ./...
```

### Executar testes paralelos
```bash
go test -parallel 4 ./...
```

### Executar testes específicos
```bash
go test -run TestConvertCelsiusToFahrenheit ./pkg/utils
```

## Troubleshooting

### Erro: "WEATHER_API_KEY is not set"
- Configure a variável de ambiente ou arquivo .env
- Ou execute apenas testes unitários com `-short`

### Erro: "network timeout"
- Verifique sua conexão com a internet
- Alguns testes fazem chamadas reais para APIs

### Erro: "API rate limit"
- As APIs têm limites de requisições
- Aguarde alguns minutos e tente novamente

## Cobertura de Testes

Para verificar a cobertura de testes:

```bash
# Gerar relatório de cobertura
go test -coverprofile=coverage.out ./...

# Visualizar no navegador
go tool cover -html=coverage.out -o coverage.html

# Ver no terminal
go tool cover -func=coverage.out
```

## Boas Práticas

1. **Execute os testes antes de fazer commit**
2. **Mantenha a cobertura de testes alta**
3. **Use testes de integração para validar o fluxo completo**
4. **Teste casos de erro e edge cases**
5. **Use mocks para testes unitários**
6. **Documente novos testes**