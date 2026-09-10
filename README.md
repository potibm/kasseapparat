# Kasseapparat

![Kasseapparat Logo](docs/kasseapparat.svg)

> _Kasseapparat_ is the Danish term for cash register.

> **Part of the Apparat Suite** ⚙️
> This tool is part of a decoupled set of single-purpose event management tools built for demoparties.
> [➔ Read more about the full Apparat suite here](https://github.com/potibm/apparat)

It is a simple point of sale (POS) system aimed at [demoparties](https://en.wikipedia.org/wiki/Demoscene#Parties). Based on [Partymeister](https://github.com/partymeister), rewritten after moving to [Granola](https://gitlab.com/granola-compo/granola) for [Evoke](https://www.evoke.eu/).

## Tooling

- [Go](https://go.dev)
  - [Gin Web Framework](https://gin-gonic.com)
  - [GORM](https://gorm.io)
  - [Cobra](https://cobra.dev) & [Viper](https://github.com/spf13/viper)
- [React](https://react.dev)
  - [Vite](https://vitejs.dev/)
  - [React Admin](https://marmelab.com/react-admin/)
  - [Flowbite React](https://flowbite-react.com) & [Tailwind CSS](https://tailwindcss.com)
- [SQLite](https://www.sqlite.org)
- Observability
  - [Sentry](https://sentry.io)
  - [OpenTelemetry](https://opentelemetry.io)
- Development & Ops
  - [mise](https://mise.jdx.dev/)
  - [Docker](https://www.docker.com)

## Quickstart

We use `mise` to automatically manage all tool versions (Go, Node, etc.) and project tasks.

```bash
# 1. Install mise (if not already installed)
curl https://mise.run | sh

# 2. Setup the project (installs dependencies and starts infra)
mise run setup

# 3. Start local services (Traefik, Redis, MinIO, etc.)
mise run infra:up

# 4. Start the development server (hot-reload for backend & frontend)
.mise/tasks/dev
```

## Documentation

- [Developer Guide](docs/dev.md)
- [Admin Documentation](docs/admin.md)
- [User Documentation](docs/manual.md)
- [SumUp Integration Documentation](docs/sumup.md)
- [Image Signing and SBOM Attestations](docs/supply-chain.md)

## Credits & License

- **Logo Icon:** [Cash Register](https://fontawesome.com/icons/cash-register) by [Font Awesome](https://fontawesome.com) is licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).
- **Software:** Licensed under MIT.
