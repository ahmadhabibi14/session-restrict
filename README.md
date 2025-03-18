# A Dummy Session Restriction app

> Built with:
> - Golang
> - Fiber
> - PostgreSQL
> - Redis

## How to run ?

1. Setup Dependencies
```bash
make setup
```
2. Download Go dependencies
```bash
go mod download
```
3. Start Docker containers
```bash
docker compose up -d
```
4. Load database
```bash
dbmate up
```
5. Start development
```bash
air

pnpm dev
```
