# Frontend - DeFi Vision

This is the frontend for DeFi Vision, a crypto portfolio analyzer. It is built with React, TypeScript, Tailwind CSS, and Vite.

## Features
- Analyze any EVM wallet address
- Display token holdings, DeFi positions, risk analysis, and AI recommendations
- Responsive design for desktop and mobile
- Interactive cards, charts, and badges

## Getting Started

### Prerequisites
- Node.js (v18+ recommended)
- Bun (if using Bun)

### Installation
```bash
cd frontend
npm install # or bun install
```

### Running Locally
```bash
npm run dev # or bun run dev
```
The app will be available at http://localhost:3000

### Configuration
- The frontend expects the backend API to be running at http://localhost:8080
- You can change the API URL in the code if needed

## Folder Structure
- `src/` - Main source code
- `components/` - UI components
- `pages/` - Page components
- `hooks/` - Custom React hooks
- `lib/` - Utility functions

## Development
- Uses Tailwind CSS for styling
- Hot reload enabled
- Linting and formatting via ESLint and Prettier

## Deployment
- Build with `npm run build` or `bun run build`
- Deploy static files to Vercel, Netlify, or any static host

## License
MIT
