# prueba-tecnica

```
prueba-tecnica
├─ cmd
│  └─ api
│     └─ main.go
├─ docker-compose.yml
├─ go.mod
├─ go.sum
├─ internal
│  ├─ adapters
│  │  └─ storage
│  │     ├─ customer_repository.go
│  │     ├─ db.go
│  │     └─ workorder_repository.go
│  └─ core
│     ├─ domain
│     │  ├─ customer.go
│     │  └─ workorder.go
│     ├─ ports
│     │  └─ ports.go
│     └─ services
│        └─ services.go
├─ Makefile
├─ migrations
│  ├─ 001_create_initial_tables.down.sql
│  └─ 001_create_initial_tables.up.sql
└─ README.md

```