# Seed User Credentials

These credentials match the records inserted by [seed.sql](/home/rowana/projects/gated-community/comune/database/seed.sql).

## Seeded tenancy

- Organizations: `20`
- Communities: `20`
- Seed shape: medium-sized dataset with varied operational records across the currently implemented auth, membership, community, unit, household, resident, invite, and invoicing flows

The first seeded tenancy remains:

- Organization: `Palm View Property Management`
- Community: `Palm View Estate`

## Local development users

All seeded users share the same password:

- Password: `Comune123!`

Example accounts for the first seeded organization:

- Owner: `owner@palmview.local`
- Org admin / community admin: `admin@palmview.local`
- Resident: `resident@palmview.local`
- Security: `security@palmview.local`
- Operations / manager: `ops@palmview.local`

Other seeded organizations follow the same role-based pattern with their own domain, for example:

- `owner@harbourcrest.local`
- `admin@coralgardens.local`
- `resident@sunrisemeadows.local`

This seed now populates the active invoicing tables as part of the next feature slice:

- `invoice_plans`
- `invoice_plan_unit_overrides`
- `invoices`

It still intentionally does not populate future-facing tables such as visitors, maintenance requests, notifications, audit logs, outbox events, or payment-proof workflows until those flows are wired into the application.

## How to apply the seed

1. Apply the schema: `make db-init`
2. Apply the seed: `make db-seed`

If you are using the checked-in Docker database defaults directly, set:

```sh
export DATABASE_URL=postgres://user:pass@localhost:5432/comune_dev?sslmode=disable
```
