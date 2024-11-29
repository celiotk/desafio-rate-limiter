# Rate Limiter com Redis e Janela Fixa

Este projeto implementa um **Rate Limiter** em Go utilizando Redis para armazenamento de estado. Ele controla o número de requisições por segundo com base no endereço IP e/ou no token de acesso enviado na requisição.

## Como Funciona

1. **Janela Fixa (Fixed Window)**:
   - A lógica do rate limiter utiliza o modelo de **janela fixa**.
   - Uma nova janela de tempo é iniciada assim que uma requisição é recebida de um IP ou token cuja janela ainda não foi registrada.
   - Durante a janela de tempo, o sistema monitora a quantidade de requisições e bloqueia requisições excedentes.

2. **Prioridade do Token**:
   - Caso a requisição contenha um **token** (especificado no cabeçalho `API_KEY: <TOKEN>`), as configurações do token terão prioridade sobre as configurações de IP.
   - Se não houver configurações específicas para o token, serão utilizadas configurações padrão definidas pelas variáveis `DEFAULT_TOKEN_RATE_LIMIT` e `DEFAULT_TOKEN_RATE_INTERVAL`.

3. **Tempo de Bloqueio**:
   - Quando o limite de requisições é atingido, o IP ou token será **bloqueado** por um período de tempo definido:
     - Para IPs: configurado com `IP_BLOCK_TIME`.
     - Para tokens: configurado com `TOKEN_BLOCK_TIME`.
   - Durante o bloqueio, todas as requisições retornarão o código HTTP `429` com a mensagem:
     ```
     you have reached the maximum number of requests or actions allowed within a certain time frame
     ```

## Configuração

O rate limiter pode ser configurado no arquivo `.env` ou passadas por meio de variáveis de ambiente através do `docker-compose.yaml`

### Exemplo de configuração
Abaixo está um exemplo de configuração:

```
IP_RATE_LIMIT=10
IP_RATE_INTERVAL=1
IP_BLOCK_TIME=60
DEFAULT_TOKEN_RATE_LIMIT=30
DEFAULT_TOKEN_RATE_INTERVAL=30
TOKEN_BLOCK_TIME=60
TOKENS=ABC123:10/1,DEF456:20
```

### Explicação das Variáveis

**Limites por IP**
- `IP_RATE_LIMIT`: Número máximo de requisições permitidas por intervalo de tempo para um único IP.  
- `IP_RATE_INTERVAL`: Intervalo de tempo (em segundos) da janela fixa para IPs.  
- `IP_BLOCK_TIME`: Tempo de bloqueio (em segundos) para um IP que ultrapassou o limite.

**Limites por Token**
- `DEFAULT_TOKEN_RATE_LIMIT`: Limite padrão de requisições por intervalo de tempo para tokens que não possuem configurações específicas.  
- `DEFAULT_TOKEN_RATE_INTERVAL`: Intervalo padrão (em segundos) da janela fixa para tokens sem configurações específicas.  
- `TOKEN_BLOCK_TIME`: Tempo de bloqueio (em segundos) para um token que ultrapassou o limite.

**Configurações de Tokens Específicos** 
- `TOKENS`: Lista de tokens com suas configurações específicas, no formato `<token>:<limite>/<intervalo>`. Exemplo:
    - `ABC123:10/1`: O token `ABC123` permite 10 requisições por segundo.
    - `DEF456:20/2`: O token `DEF456` permite 20 requisições a cada 2 segundos.
## Mecanismo de Persistência
O rate limiter foi projetado para suportar diferentes mecanismos de persistência. Ele utiliza uma interface chamada RateLimiterStorage, que abstrai as operações de armazenamento.

### Implementações Disponíveis:
1. **Redis:**
    - Redis é a implementação utilizado por padrão.
2. **In-Memory (Memória):**
    - Uma implementação em memória está disponível no arquivo `in_memory_storage.go`.
    - Essa implementação é utilizada nos testes automatizados para evitar dependência de serviços externos.

Para trocar o mecanismo de persistência, basta implementar a interface `RateLimiterStorage` e registrar a nova implementação.

## Testes Automatizados
O sistema possui testes automatizados que validam o funcionamento do rate limiter. O teste principal está localizado no arquivo:

```
internal/infra/web/middleware/rate_limiter_test.go
```