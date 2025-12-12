# SPECTRA 🌐

**Intelligent Fiber Infrastructure System**

A comprehensive fiber optic infrastructure management system with GIS mapping, network management, and real-time monitoring capabilities.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?logo=postgresql)

---

## ✨ Features

- **GIS Mapping** — Interactive map visualization with Mapbox GL JS
- **Node Management** — OLT, ODC, ODP, Closure, Pole, Customer points
- **Cable Tracking** — ADSS, Duct, Drop cables with fiber core management
- **Splice Matrix** — Connection management between fiber cores
- **Customer Monitoring** — ONT status tracking (Online, Offline, LOS)
- **GeoJSON Export** — Direct integration with mapping libraries

---

## 🏗️ Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│    Frontend     │────▶│     Backend     │────▶│   PostgreSQL    │
│  React + Vite   │     │    Go (Gin)     │     │                 │
│  :5173          │     │    :8080        │     │    spectra      │
└─────────────────┘     └─────────────────┘     └─────────────────┘
```

---

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+
- PostgreSQL 15+
- Mapbox Access Token

### 1. Clone Repository

```bash
git clone https://github.com/yogisyahroni/spectra.git
cd spectra
```

### 2. Setup Database

```bash
# Create database
psql -U postgres -c "CREATE DATABASE spectra;"
```

### 3. Start Backend

```bash
cd backend
cp .env.example .env
# Edit .env with your database credentials

go mod tidy
go run ./cmd/server/main.go
```

Backend runs at: <http://localhost:8080>

### 4. Start Frontend

```bash
cd frontend
cp .env.example .env
# Add your Mapbox token to .env

npm install
npm run dev
```

Frontend runs at: <http://localhost:5173>

---

## 📁 Project Structure

```
spectra/
├── backend/                 # Go API Server
│   ├── cmd/server/          # Entry point
│   ├── internal/
│   │   ├── config/          # Environment config
│   │   ├── database/        # DB connection & migrations
│   │   ├── handlers/        # HTTP handlers
│   │   ├── middleware/      # CORS, logging
│   │   ├── models/          # Data models
│   │   ├── repository/      # Database operations
│   │   └── routes/          # API routes
│   └── .env.example
│
├── frontend/                # React Application
│   ├── src/
│   │   ├── components/      # UI components
│   │   ├── pages/           # Page components
│   │   ├── services/        # API client
│   │   ├── types/           # TypeScript types
│   │   └── lib/             # Utilities
│   └── .env.example
│
└── README.md
```

---

## 🔌 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| GET/POST | `/api/nodes` | List/Create nodes |
| GET/PUT/DELETE | `/api/nodes/{id}` | Get/Update/Delete node |
| GET | `/api/nodes/nearby` | Get nearby nodes |
| GET | `/api/geojson/nodes` | Export nodes as GeoJSON |
| GET/POST | `/api/cables` | List/Create cables |
| GET | `/api/cables/{id}/cores` | Get cable cores |
| GET/POST | `/api/connections` | List/Create splices |
| GET/POST | `/api/customers` | List/Create customers |
| GET | `/api/customers/los` | Get LOS customers |

---

## 🛠️ Tech Stack

### Backend

- **Go** — High-performance API server
- **pgx** — PostgreSQL driver
- **net/http** — HTTP router (Go 1.22+)

### Frontend

- **React 18** — UI library
- **Vite** — Build tool
- **TypeScript** — Type safety
- **Tailwind CSS** — Styling
- **Mapbox GL JS** — Map visualization
- **TanStack Query** — Data fetching
- **Lucide React** — Icons

---

## 🔐 Environment Variables

### Backend (`backend/.env`)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=spectra
SERVER_PORT=8080
```

### Frontend (`frontend/.env`)

```env
VITE_API_URL=http://localhost:8080
VITE_MAPBOX_TOKEN=pk.your_mapbox_token
```

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

## 👤 Author

**Yogi Syahroni**  
GitHub: [@yogisyahroni](https://github.com/yogisyahroni)
