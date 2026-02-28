# GoWA : Golang Web Application

## Architecture

- [Frontend](./frontend/) : <img src="./frontend/public/next.svg" alt="ts-logo" height="12"/>

- [Backend](./backend/) : <img src="./frontend/public/gologo.svg" alt="go-logo" height="12"/>

- [Container](./compose.yml) : Docker

## Install

```bash
docker compose build
```

## Run

```bash
docker compose up nextapp
```

## Testing

### Unit Tests (Frontend)

Run Jest tests locally:

```bash
cd frontend && npm test
```

### E2E Tests (Frontend + Backend)

Run full E2E tests with backend and database orchestrated together:

```bash
npm run test:e2e
```

Or from the frontend directory:

```bash
docker compose --profile test up frontend-test --abort-on-container-exit
```

Clean up containers after testing:

```bash
npm run test:e2e:down
```

**Note:** The E2E test service automatically handles:

- Starting the PostgreSQL database
- Starting the backend API server
- Waiting for services to be healthy before running tests
- Using internal container networking (backend at `http://goapp:8000`)

## Learning Resources

- [tutorial](https://dev.to/francescoxx/go-typescript-full-stack-web-app-with-nextjs-postgresql-and-docker-42ln)
