# AGENTS.md

## Scope

These rules apply to everything under `database/`.

The database layer is a first-class part of the product design. Changes here should preserve tenant safety, operational clarity, realistic local testing, and a clean path to production migrations.

## Schema Rules

- Treat `organizations` as the tenant boundary and preserve explicit organization/community scoping in schema design.
- Prefer additive changes and safe rollout paths over destructive rewrites.
- Keep constraints, indexes, and status fields explicit so business rules are enforced close to the data.
- Update both the schema snapshot and any relevant migrations or notes when the intended database shape changes.
- Keep `gated-community-schema.sql` suitable for fresh bootstrap and keep numbered migrations suitable for evolving existing environments.

## Migration Rules

- New schema changes should land as numbered migrations under `database/migrations/`.
- Migrations should be ordered, reviewable, and safe to run once.
- If a migration changes a business rule or table contract, update schema notes in the same change.
- Avoid hidden data rewrites. If backfills or compatibility steps are needed, make them explicit.

## Seed Data Rules

- When schema changes or new features affect developer workflows, update `seed.sql` in the same change unless there is a deliberate reason not to.
- Seed data should represent a medium-sized, realistic testing dataset rather than the smallest possible happy path.
- Prefer seeded data that exercises multiple roles, multiple units, multiple residents, lifecycle states, and important operational paths such as invites, invoicing, payments, visitors, and maintenance.
- Keep seeded records tenant-safe and internally consistent so local development can validate real workflows without manual cleanup.
- If a feature is not yet fully represented in seed data, leave a clear path to add it soon rather than letting the seed script drift behind the schema.
- Do not populate future-facing tables just because they exist in the schema. Prefer seeding the tables that are actually used by current application flows, then expand the dataset when those modules are implemented.

## Review Checklist

- Does the schema change preserve tenant scoping and auditability?
- Is there a migration for the change if existing environments need to evolve?
- Does `seed.sql` still apply cleanly after this change?
- Does the seeded dataset still provide enough realistic coverage for local testing?
- If this change adds or materially changes database governance docs, does `docs/primer-prompt.md` still reference the right sources and workflow expectations?
