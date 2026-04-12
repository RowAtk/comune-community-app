# Seed User Credentials

These credentials match the records inserted by [seed.sql](/home/rowana/projects/gated-community/comune/database/seed.sql).

## Seeded tenancy

- Organization: `Palm View Property Management`
- Community: `Palm View Estate`

## Local development users

All seeded users share the same password:

- Password: `Comune123!`

User accounts:

- Owner: `owner@palmview.local`
- Org admin / community admin: `admin@palmview.local`
- Resident: `resident@palmview.local`
- Security: `security@palmview.local`

## How to apply the seed

1. Apply the schema: `make db-init`
2. Apply the seed: `make db-seed`

If you are using the checked-in Docker database defaults directly, set:

```sh
export DATABASE_URL=postgres://user:pass@localhost:5432/comune_dev?sslmode=disable
```
