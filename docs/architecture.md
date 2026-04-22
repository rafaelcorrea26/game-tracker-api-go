## 📄 docs/architecture.md
```md
# 🧱 Arquitetura - Backend

## 🎯 Objetivo

API performática e segura.

---

## 📁 Estrutura


handlers/
database/
models/
middlewares/


---

## 🔄 Fluxo

1. Request entra
2. Middleware valida JWT
3. Handler executa lógica
4. Banco executa query
5. Retorna JSON

---

## 🔐 Segurança

- JWT
- Isolamento por usuário

---

## 🗄 Banco

- users
- games