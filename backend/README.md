# Backend - DeFi Vision

This is the backend for DeFi Vision, a crypto portfolio analyzer. It is built with Go, mux, and integrates with Zapper and other APIs.

## Features
- REST API for portfolio analysis
- Fetches token balances and DeFi positions from Zapper
- AI-powered risk analysis and recommendations
- CORS enabled for local development

## Getting Started

### Prerequisites
- Go 1.20+

### Installation
```bash
cd backend
go mod tidy
```

### Running Locally
```bash
go run cmd/main.go
```
The API will be available at http://localhost:8080

### API Endpoints
- `GET /analyze?address=<wallet_address>`: Returns JSON with recommended_tokens, risk_score, reasoning, token_balances, app_balances
- `GET /positions?address=<wallet_address>`: Returns raw positions data

## Configuration
- Environment variables can be set in `.env`
- See `cmd/main.go` for server setup

## Development
- Modular Go codebase
- Logging enabled
- Error handling for API call

## Deployment
- Build with `go build`
- Deploy to any cloud or server

## License
MIT
