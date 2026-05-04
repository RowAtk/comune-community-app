---
name: refresh-seed-data
description: Refresh or regenerate database/seed.sql for this repository when schema or implemented application flows change, keeping the seed realistic, medium-sized, and limited to tables the app actually uses today unless future modules explicitly need early data.
---

# Refresh Seed Data

Use this skill when updating `database/seed.sql`, `database/seed-credentials.md`, or nearby seed strategy docs.

## Goal

Keep local seed data:

- aligned with the current schema
- aligned with currently implemented application flows
- realistic enough for manual and automated testing
- easy to reapply during local development

## Repo Conventions To Follow

- Seed SQL lives in `database/seed.sql`.
- Seed strategy expectations live in:
  - `database/AGENTS.md`
  - `database/MIGRATIONS.md`
  - `docs/implementation-checklist.md`
- The app currently uses:
  - auth
  - organizations
  - communities
  - units
  - households
  - residents
  - organization/community/resident invite flows
- Do not assume every schema table needs data just because it exists.

## Workflow

1. Inspect what the app actually uses.
   Check:
   - `apps/api/internal/modules/...`
   - `apps/web/src/routes/...`
   - any relevant roadmap or checklist docs

2. Decide the seeded scope.
   Default to:
   - seed the tables needed by implemented flows
   - skip future-facing unused tables
   - expand only when the application starts depending on those tables

3. Keep the dataset medium-sized and realistic.
   - include at least 20 rows overall for each actively seeded entity
   - vary names, statuses, roles, and lifecycle states
   - keep tenant relationships internally consistent

4. Keep local testing easy.
   - preserve one easy-to-recognize primary tenancy
   - keep example login accounts documented in `database/seed-credentials.md`
   - prefer deterministic IDs or generation patterns where repeatability helps

5. Update adjacent docs in the same change.
   Usually:
   - `database/seed-credentials.md`
   - `database/MIGRATIONS.md` if the strategy changed
   - `docs/implementation-checklist.md` if seed scope or progress changed

## Strong Defaults

- Prefer deterministic seed generation over huge hand-written value lists when it improves maintainability.
- Prefer realistic variation over perfectly uniform sample rows.
- Seed statuses should exercise real workflows: active, invited, suspended, pending, moved out, accepted.
- Do not populate invoicing, visitors, maintenance, notifications, audit logs, or outbox data until those flows are actually implemented, unless a task explicitly needs them.

## Verification

Use the repo's existing commands:

```sh
make db-init-docker
make db-seed-docker
```

Or locally:

```sh
make db-init
make db-seed
```

## Done Criteria

- `database/seed.sql` matches the current implemented product surface
- each actively seeded entity has at least 20 rows overall
- seed data is realistic and varied
- unused future-facing tables are left unseeded unless intentionally needed
- seed credentials and strategy docs are updated if necessary
