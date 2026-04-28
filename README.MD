# Game Tracker API

> Backend da aplicação Game Tracker — rastreamento de progresso em jogos.

API REST construída em Go com Gin, PostgreSQL e autenticação JWT.

## Stack

| Tecnologia | Versão | Uso |
|---|---|---|
| Go | 1.22+ | Linguagem |
| Gin | 1.12 | Framework HTTP |
| PostgreSQL | 15+ | Banco de dados |
| lib/pq | 1.11 | Driver PostgreSQL |
| golang-jwt/jwt | v5 | Autenticação |
| bcrypt | — | Hash de senhas |
| godotenv | 1.5 | Variáveis de ambiente |
| gin-contrib/cors | 1.7 | CORS |
| swaggo/gin-swagger | 1.6 | Documentação Swagger |

## Endpoints

### Autenticação (público)

| Método | Rota | Descrição |
|---|---|---|
| POST | `/register` | Criar conta |
| POST | `/login` | Login — retorna JWT |

### Usuário (autenticado)

| Método | Rota | Descrição |
|---|---|---|
| PUT | `/change-password` | Trocar senha |

### Jogos (autenticado)

| Método | Rota | Descrição |
|---|---|---|
| GET | `/games` | Listar jogos do usuário |
| POST | `/games` | Adicionar jogo |
| PUT | `/games/:id` | Editar jogo |
| DELETE | `/games/:id` | Excluir jogo |

### Documentação

| Rota | Descrição |
|---|---|
| `GET /swagger/*any` | Interface Swagger UI |

## Estrutura do projeto

```
game-tracker-api-go/
├── main.go               # Ponto de entrada, rotas, CORS
├── database/
│   ├── connection.go     # Conexão com PostgreSQL
│   └── migrate.go        # Executa migrations automaticamente
├── middlewares/
│   └── auth.go           # Validação JWT
├── utils/
│   └── jwt.go            # Geração e parsing de tokens
├── models/
│   ├── user.go
│   └── games.go
├── handlers/
│   ├── auth.go           # Register, Login, ChangePassword
│   └── games.go          # CRUD de jogos
├── docs/                 # Gerado pelo swaggo
└── migrations/           # 5 arquivos SQL
```

## Banco de dados

5 migrations executadas automaticamente na inicialização:

```
001_create_games.sql       — tabela de jogos
002_seed_games.sql         — dados de exemplo
003_create_users.sql       — tabela de usuários
004_add_platform.sql       — coluna platform em games
005_add_dates_to_games.sql — start_date / end_date
```

**Schema principal:**

```sql
users  → id, username (UNIQUE), password, created_at
games  → id, name, status, rating, notes, platform,
          user_id (FK), start_date, end_date, created_at
```

**Status de jogo:** `playing` | `completed` | `backlog`

## Rodando localmente

### Pré-requisitos

- Go 1.22+
- PostgreSQL 15+

### Setup

```bash
# Clone o repositório
git clone https://github.com/rafaelcorrea26/game-tracker-api-go.git
cd game-tracker-api-go

# Configure as variáveis de ambiente
cp .env.example .env
# Edite .env com suas credenciais

# Rode
go run main.go
```

As migrations criam todas as tabelas automaticamente na primeira execução.

**Swagger disponível em:** `http://localhost:8080/swagger/index.html`

## Variáveis de ambiente

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=sua_senha
DB_NAME=gamertracker
DB_SSLMODE=disable        # use "require" em produção

PORT=8080
JWT_SECRET=chave_forte_e_aleatoria_aqui
FRONTEND_URL=http://localhost:5173
```

## Deploy

Qualquer plataforma que execute binários Go:

- **Railway** — detecta Go automaticamente, configura variáveis de ambiente e faz deploy
- **Render** — mesma abordagem
- **VPS** — `go build -o server . && ./server`

### CORS em produção

O CORS está configurado para aceitar automaticamente:
- `http://localhost:5173` — desenvolvimento local
- `FRONTEND_URL` — URL de produção (variável de ambiente)
- `https://game-tracker-frontend-react-*.vercel.app` — previews do Vercel
