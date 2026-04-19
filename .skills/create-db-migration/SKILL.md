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
   Use a descriptive filename such as:
   - `database/alter-residents-add-invite-token.sql`
   - `database/create-community-visitors.sql`
   - `database/alter-households-add-status.sql`

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

7. Update nearby schema notes when the business rule is non-obvious.

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
make db-seed-docker
make db-docker-run SCRIPT=database/your-migration.sql
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

## Done Criteria

- Migration file exists in `database/` with a clear name.
- Schema change is tenant-safe and reviewable.
- Execution path is realistic using the repo's Make targets.
- Any subtle domain rule is reflected in adjacent schema notes.
