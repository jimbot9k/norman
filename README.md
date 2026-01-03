# Norman
> Schema & access safety auditing for your database - know what your database guarantees, and where those guarantees are fragile or missing.

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](LICENSE)

## Why

- **Trust & Transparency** — Norman tells you exactly what it can see in your database schema, nothing more
- **DDL-Level Analysis** — Analyzes schema structure and guarantees without inspecting row data
- **Confidence Levels** — Opinionated but cautious; every finding includes confidence scoring
- **Machine & Human Readable** — Export to JSON for automation or Mermaid ERD for documentation
- **PostgreSQL-First** — Deep support for PostgreSQL with MySQL support in development

## Install

### Binaries

Download pre-built binaries from [GitHub Releases](https://github.com/jimbot9k/norman/releases).

### Go

```bash
go install github.com/jimbot9k/norman@latest
```

## Quickstart

```bash
# Show help
norman --help

# Generate inventory reports (JSON + Mermaid ERD)
norman --conn "postgres://user:password@localhost:5432/mydb"

# Specify output directory
norman --conn "postgres://user:password@localhost:5432/mydb" --output-dir ./reports/

# Generate only JSON report
norman --conn "postgres://user:password@localhost:5432/mydb" --report-types json

# Generate only Mermaid ERD
norman --conn "postgres://user:password@localhost:5432/mydb" --report-types mermaid
```

## Configuration

Norman is configured entirely via command-line flags:

| Flag | Default | Description |
|------|---------|-------------|
| `--conn` | *(required)* | Database connection string |
| `--output-dir` | `./norman/` | Directory to output reports to |
| `--report-types` | `all` | Comma-separated list of report types (`json`, `mermaid`, `all`) |

### Supported Databases

| Database | Connection String Format |
|----------|-------------------------|
| PostgreSQL | `postgres://user:password@host:port/database` |
| MySQL | `mysql://user:password@tcp(host:port)/database` |

### Output Formats

- **JSON** — Machine-readable schema inventory with full metadata
- **Mermaid** — ERD diagram in Mermaid syntax (`.mmd`) for documentation

## Roadmap

Norman is building toward a credible **v1.0** release focused on schema & access safety auditing.

| Version | Focus | Status |
|---------|-------|--------|
| **v0.1** | Schema Inventory — tables, columns, keys, indexes, views, functions | 🚧 In Progress |
| v0.2 | Integrity Gaps — missing FKs, constraints, indexes | Planned |
| v0.3 | Normalization Smells — duplication, coupling patterns | Planned |
| v0.4 | Multi-Tenancy Safety — tenant discriminators, isolation gaps | Planned |
| v0.5 | PostgreSQL RLS Analysis — policy correctness, bypass risks | Planned |
| v0.6 | Access & Privilege Risk — over-privileged roles, GRANT analysis | Planned |
| v0.7 | Schema Entropy — orphaned objects, naming inconsistencies | Planned |
| v0.8 | Migration Risk — evolution blockers, dangerous patterns | Planned |
| **v1.0** | Full Schema & Access Audit — stable CLI, JSON schema, HTML reports | Planned |

## License

Mozilla Public License Version 2.0 — see [LICENSE](LICENSE)
