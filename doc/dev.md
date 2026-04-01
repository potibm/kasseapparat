# Kasseapparat: Developer Guide

Welcome to the development documentation for Kasseapparat! We use a highly automated, modern development environment to make contributing as frictionless as possible.

## 🛠 Prerequisites

You only need two things installed on your machine to work on this project:

- [Docker](https://docs.docker.com/get-docker/) (for running infrastructure like databases and mail servers)
- [mise](https://mise.jdx.dev/) (The polyglot tool manager)

**That's it**! You do not need to manually install Go, Node.js, Yarn, or any linters. mise will automatically download and use the exact versions specified in our mise.toml to ensure 100% reproducible builds across all machines.

## 🚀 Getting Started

To initialize the project, download all dependencies, and start the infrastructure containers, simply run:

```bash
mise run setup
```

Afterward, boot up the entire stack (Backend + Frontend) with hot-reloading enabled:

```bash
mise run dev
```

_Note: This uses **overmind** under the hood to stream all logs into a single terminal window._

## 📋 Task Reference

Our tasks are neatly categorized. You can always run `mise tasks` to see an interactive list of all available commands.

Here are the most important ones:

### Development & Testing

- `mise run be:dev` - Start only the backend with live-reload (air).
- `mise run fe:dev` - Start only the frontend with live-reload.
- `mise run test` - Run all frontend and backend tests in parallel.
- `mise run e2e:run` - Boot a clean database and run Playwright end-to-end tests.

### Linting & Formatting

We enforce strict code quality rules for both Go and TypeScript.

- `mise run lint` - Run all linters (Backend & Frontend) in parallel.
- `mise run lint --fix` - Automatically fix formatting issues and linting errors.

### Infrastructure

- `mise run infra:up` - Start the Docker Compose stack (Mailhog, OpenObserve, RedisInsight).
- `mise run infra:down` - Stop the stack.

### Building & Docker

- `mise run build` - Compile both frontend and backend into the /dist directory.
- `mise run docker:build` - Build the production-ready Docker image (kasseapparat:latest).
- `mise run docker:build --clean` - Force a fresh build without cache and pull the latest base images.
