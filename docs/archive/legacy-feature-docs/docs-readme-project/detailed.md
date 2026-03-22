# Source: docs/README_PROJECT.md

# Universal Data Viewer (UDV)

A data visualization and query tool that works with PostgreSQL databases.

## Quick Start

### Prerequisites
- Go 1.22+
- Node.js 18+
- PostgreSQL 12+

### Backend Setup

```bash
go run ./cmd/server
```

Server starts on `http://localhost:8080`

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

Frontend starts on `http://localhost:3000`

## Project Structure

See [docs/repo_strucutre.md](docs/repo_strucutre.md) for detailed architecture.

```
udv/
├── cmd/server/          # Application entry point
├── internal/            # Core backend logic
├── frontend/            # React frontend
├── configs/             # Configuration files
├── docs/                # Architecture documents
└── tests/               # Test files
```

## Documentation

- [Architecture Overview](docs/technical.md)
- [Backend Design](docs/backend.md)
- [Frontend Design](docs/frontend.md)
- [Development Playbook](docs/development_plan.md)
- [MVP Scope](docs/mvp_scope.md)
- [Query DSL Spec](docs/query_dsl_spec.md)
- [PostgreSQL Adapter](docs/postgres_adapter_skeleton.md)

## Development

Follow the [Development Playbook](docs/development_plan.md) for step-by-step implementation guidance.

## License

MIT
