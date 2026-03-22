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

See [repo structure](/Users/shubhamparamhans/Workspace/udv/docs/architecture/repo-structure.md) for detailed architecture.

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

- [Architecture Overview](/Users/shubhamparamhans/Workspace/udv/docs/architecture/system-overview.md)
- [Backend Design](/Users/shubhamparamhans/Workspace/udv/docs/architecture/backend.md)
- [Frontend Design](/Users/shubhamparamhans/Workspace/udv/docs/architecture/frontend.md)
- [Development Playbook](/Users/shubhamparamhans/Workspace/udv/docs/engineering/development-playbook.md)
- [MVP Scope](/Users/shubhamparamhans/Workspace/udv/docs/product/mvp-scope.md)
- [Query DSL Spec](/Users/shubhamparamhans/Workspace/udv/docs/architecture/query-dsl.md)
- [PostgreSQL Support](/Users/shubhamparamhans/Workspace/udv/docs/features/postgres-support/summary.md)

## Development

Follow the [Development Playbook](/Users/shubhamparamhans/Workspace/udv/docs/engineering/development-playbook.md) for step-by-step implementation guidance.

## License

MIT
