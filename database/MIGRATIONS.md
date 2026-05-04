# Database Migrations

This repo now uses two complementary database entrypoints:

- `database/gated-community-schema.sql` for fresh bootstrap of a brand-new database
- `database/migrations/*.sql` for additive schema evolution over time

## Recommended workflow

For local resettable development:

1. Use `make db-init` or `make db-init-docker` to bootstrap a fresh database from `database/gated-community-schema.sql`
2. Use `make db-seed` or `make db-seed-docker` if you want seed data

For ongoing schema changes:

1. Create a new numbered migration in `database/migrations/`
2. Keep migrations additive and safe to run once in sequence
3. Apply them with `make db-migrate` or `make db-migrate-docker`
4. Periodically fold stable changes back into `database/gated-community-schema.sql` so fresh environments stay easy to bootstrap

## Naming

Use sortable filenames such as:

- `001_add_resident_invitations.sql`
- `002_add_household_role.sql`
- `003_add_invoicing_module.sql`

## Rules

- Prefer additive changes over destructive rewrites
- If a migration changes business rules, update schema notes and related docs in the same change
- Keep tenant scoping explicit in tables, indexes, and constraints
- Use `schema_migrations` to record applied files rather than relying on memory or manual notes
- Keep `database/seed.sql` aligned with the current schema and feature set
- Seed data should include at least 20 rows per seeded entity overall and should use varied, realistic values rather than uniform placeholders
- Treat seed data as an operational testing fixture, not just a smoke-test bootstrap
- Seed only the entities needed by currently implemented application flows unless a change explicitly requires early data for a future module

## Current migration baseline

The current migration set assumes:

- older environments may already exist from the pre-migration `gated-community-schema.sql`
- fresh environments can still bootstrap directly from the schema snapshot
- new changes should land as migrations first, then be folded back into the schema snapshot as needed
