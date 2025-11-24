# Borg

## 🧠 Overview

- **Backend:** Go (`net/http` + embedded FS)
- **Frontend:** React + Rsbuild + pnpm
- **Database:** PostgreSQL
- **Migrations:** [Goose](https://github.com/pressly/goose)
- **SQL generation:** [sqlc](https://sqlc.dev)

---

## 🚀 Quick Start

### Running In Development Mode

> Requires Go ≥ 1.25, pnpm, make and container engine running locally.

```bash
make dev
```

---

## 🧩 Project Structure

```
.
├── cmd/
│   └── borg/               # main entrypoint
├── pkg/
│   ├── db/                 # database setup and interfaces
│   ├── config/             # configuration management
│   └── router/             # http routes and handlers
└── web/                    # React SPA source
```
