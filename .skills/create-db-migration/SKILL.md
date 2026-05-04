---
name: create-db-migration
description: Create a PostgreSQL migration for this repository when the task involves schema changes in database/, including new tables, indexes, constraints, backfills, or safe alterations for the multi-tenant gated community SaaS.
---

# Create DB Migration

Use this skill when adding or modifying SQL schema in `database/`.

## Goal

Produce a migration that is:

- safe for multi-tenant production data
- readable by humans
- aligned with the current schema style
- easy to apply locally with the repo's existing commands

## Repo Conventions To Follow

- SQL lives in `database/`.
- Existing schema files are handwritten and heavily commented.
- The domain model is organization-first, with many community resources carrying both `organization_id` and `community_id`.
- Jamaica defaults matter where country/timezone defaults are introduced.
- Audit columns commonly include `created_at`, `updated_at`, and sometimes `deleted_at`.
- Existing schema uses:
  - `gen_random_uuid()`
  - check constraints for statuses
  - explicit indexes
  - trigger-managed `updated_at`

## Workflow

1. Read the relevant schema area first.
   Start with:
   - `database/schema-with-notes.sql`
   - `database/gated-community-schema.sql`
   - any related alteration files already in `database/`

2. Choose the migration shape.
   Typical options:
   - create new table(s)
   - add nullable column, backfill, then tighten later
   - add constraints
   - add indexes for tenant-scoped access paths
   - add supporting trigger wiring for `updated_at`

3. Name the file clearly.
   Use the numbered migration workflow in `database/migrations/`, for example:
   - `database/migrations/004_add_community_visitors.sql`
   - `database/migrations/005_add_households_status.sql`
   - `database/migrations/006_add_invoicing_reviews.sql`

4. Make tenant boundaries explicit in the schema.
   - include `organization_id` where operational scoping depends on it
   - include `community_id` for community resources
   - add composite uniqueness and indexes that match query patterns

5. Add constraints for business invariants.
   Common examples:
   - status check constraints
   - scoped uniqueness
   - non-null ownership references
   - date ordering rules where safe to enforce in SQL

6. Consider rollout safety.
   - avoid destructive rewrites when additive changes work
   - if backfill is needed, make the steps obvious
   - do not silently break seed data
   - if the change affects active developer workflows, update `database/seed.sql`

7. Update nearby schema notes when the business rule is non-obvious.
8. Update the schema snapshot when the intended fresh-install shape changed.
9. Update the implementation checklist if the migration completes a tracked item.

## Strong Defaults

- Prefer explicit column lists and named constraints.
- Prefer indexes for:
  - `(organization_id)`
  - `(organization_id, community_id)`
  - tenant-scoped lists ordered by recent creation/update
- Avoid `SELECT *` in any accompanying SQL examples.
- Use `JSONB` only for flexible settings/metadata, not core relationships.

## Local Execution Commands

Use the repo's existing commands:

```sh
make db-up
make db-init-docker
make db-migrate-docker
make db-seed-docker
```

If using a local `DATABASE_URL` instead of Docker:

```sh
make db-init
make db-seed
```

## Migration Review Checklist

- Is the change additive where possible?
- Does it preserve tenant isolation?
- Are constraints named and meaningful?
- Are indexes present for expected access paths?
- Does it fit current timestamp, UUID, and soft-delete conventions?
- Will local seed/setup still work after the change?
- Is the migration reflected in `database/gated-community-schema.sql` when needed for fresh bootstrap?

## Done Criteria

- Migration file exists in `database/migrations/` with a clear numbered name.
- Schema change is tenant-safe and reviewable.
- Execution path is realistic using the repo's Make targets.
- Any subtle domain rule is reflected in adjacent schema notes.
