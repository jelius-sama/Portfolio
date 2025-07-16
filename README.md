# 🚀 Portfolio — Personal Developer Site

A modern, self-hostable developer portfolio powered by a **Go backend** and **React frontend (via Bun & Vite)**.

This project is optimized for performance, portability, and developer experience, using a clean architecture with declarative metadata and zero server-side JavaScript dependencies.

---

## 📁 Project Structure

```
.
├── api             # Backend API handlers
├── assets          # Static assets (images, icons)
├── bin             # Final built Go binary
├── client          # Frontend React project (Bun + TypeScript)
├── cmd             # Go server entry point
├── config          # JSON-based configuration (server, client, metadata)
├── logger          # Logging utility
├── middleware      # Custom middleware (logger, crash recovery, etc.)
├── parser          # HTML head metadata parser
├── types           # Shared Go types
├── util            # Utility helpers
```

---

## 🧠 Features

- ✅ **Pure Go HTTP server** (no Node.js backend)
- ⚡ **React + Bun frontend** (blazing fast DX)
- 📦 **Portable binary** deployment (`bin/Portfolio-1.x.x`)
- 🔧 **Configurable via JSON**
- 🌍 **Cloudflare Tunnel | AWS EC2 | Any VPS With Linux** / Let’s Encrypt TLS-ready

---

## 🛠 Development Setup

### Backend (Go)
```bash
make dev
```

---

## 🧩 Metadata System (Dynamic page metadata not implemented)

Head `<title>`, `<meta>`, and `<link>` tags are defined per route using a single config file:

> `config/static.route.json`

- Global tags (wildcard path `"*"`).
- Merged with per-route overrides at runtime and build.

---

## 📦 Production Build

To generate a `.tar.gz` archive containing everything needed to deploy:

```bash
make archive_prod
```

This includes:
- `bin/`
- `config/`
- `LICENSE`
- `README.md`

---

## 🚀 Deploying

Extract the archive and run the binary:

```bash
tar -xzf archive_prod.tar.gz
./bin/Portfolio-1.x.x
```

Optional:
- Add a `systemd` service
- Use `certbot` for TLS
- Serve via Cloudflare Tunnel

---

## 📜 License

MIT — See [LICENSE](./LICENSE)

---

## ✨ Author

[Jelius Basumatary](https://jelius.dev) — Fullstack Web & App Developer

