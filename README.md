# Tenant API

## Run
```bash
git clone https://github.com/ahmedkltn/tenant-api.git
cd tenant-api
go mod tidy
go run ./cmd
```

## Tests
```bash
go test ./cmd -v
```

## Seed credentials

| Email | Password | Tenant | Role |
| --- | --- | --- | --- |
| admin@tenant1.com | password123 | tenant-1 | admin |
| viewer@tenant1.com | password123 | tenant-1 | viewer |
| admin@tenant2.com | password123 | tenant-2 | admin |