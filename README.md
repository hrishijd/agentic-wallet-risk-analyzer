# agentic-wallet-risk-analyzer

# DeFi Vision

DeFi Vision is a comprehensive crypto portfolio analyzer and dashboard. It provides real-time insights, risk analysis, and AI-powered recommendations for any EVM wallet address. The platform aggregates token balances, DeFi positions, and risk factors from multiple sources, presenting them in a modern, user-friendly interface.

## Features
- Analyze any EVM wallet address
- View token holdings and DeFi positions
- AI-generated risk analysis and recommendations
- Interactive dashboard with charts and cards
- Fast, secure backend API
- Responsive frontend (React + Tailwind)

## Architecture
- **Frontend:** React, TypeScript, Tailwind CSS, Vite
- **Backend:** Go, Gin, Zapper API integration
- **API:** `/analyze?address=<wallet_address>` returns portfolio, risk, and recommendations

## Setup
1. Clone the repository
2. See `frontend/README.md` and `backend/README.md` for setup instructions

## Development
- Frontend and backend run independently
- CORS is enabled for local development

## License
MIT

## Contact
For support or contributions, open an issue or pull request.