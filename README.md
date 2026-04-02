# Kasseapparat

![Kasseapparat Logo](doc/kasseapparat.svg)

> _Kasseapparat_ is the Danish term for cash register.

It is a simple point of sale (POS) system aimed at [demoparties](https://en.wikipedia.org/wiki/Demoscene#Parties). Based on [Partymeister](https://github.com/partymeister), rewritten after moving to [Granola](https://gitlab.com/granola-compo/granola) for [Evoke](https://www.evoke.eu/).

## Tooling

- [Go](https://go.dev)
  - [Gin Web Framework](https://gin-gonic.com)
  - [GORM](https://gorm.io)
- [React](https://react.dev)
  - [Flowbite React](https://flowbite-react.com)
  - [Tailwind CSS](https://tailwindcss.com)
  - [React Admin](https://marmelab.com/react-admin/)
- [SQLite](https://www.sqlite.org)
- [mise](https://mise.jdx.dev/)

## Quickstart

We use `mise` to automatically manage all tool versions (Go, Node, etc.) and project tasks.

```bash
# 1. Install mise (if not already installed)
curl https://mise.run | sh

# 2. Setup the project (installs dependencies and starts infra)
mise run setup

# 3. Start the development server (hot-reload for backend & frontend)
mise run dev
```

## Documentation

- [Developer Guide](doc/dev.md)
- [Admin Documentation](doc/admin.md)
- [User Documentation](doc/manual.md)
- [SumUp Integration Documentation](doc/sumup.md)
- [Image Signing and SBOM Attestations](doc/supply-chain.md)
