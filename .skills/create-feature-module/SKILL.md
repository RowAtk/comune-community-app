---
name: create-feature-module
description: Create a new backend feature module for this repository's Go API when the task is to add a new domain area or resource under apps/api using the existing handler-service-repository-models pattern, tenant-safe routing, and PostgreSQL-backed persistence.
---

# Create Feature Module

Use this skill when adding a new API module or submodule in `apps/api`, especially for new community features such as visitors, maintenance, billing, notifications, or security records.

## Goal

Create the smallest production-credible module that matches current repo conventions:

- vertical slice under `apps/api/internal/modules/...`
- `models.go`, `service.go`, `repository.go`, `http.go`
- thin HTTP handlers
- service-owned validation and orchestration
- repository-owned SQL
- tenant-safe route and query scoping

## Repo Conventions To Follow

- Top-level tenant boundary is `organization_id`.
- Community resources should usually carry both `organization_id` and `community_id`.
- Use `httpx.Adapt`, `httpx.Chain`, `httpx.ResponseEnvelope`, and `apperror`.
- Route params currently use `organizationID` and `id` for community ID in nested community routes.
- Services sanitize input, enforce invariants, and map storage errors into typed app errors.
- Repositories use explicit SQL with `pgx` and `pgxpool`.

## Workflow

1. Inspect adjacent modules before adding files.
   Good references:
   - `apps/api/internal/modules/organizations`
   - `apps/api/internal/modules/communities`
   - `apps/api/internal/modules/communities/units`
   - `apps/api/internal/modules/communities/residents`

2. Decide the module placement.
   - Organization-scoped feature: `apps/api/internal/modules/<feature>`
   - Community-scoped feature: `apps/api/internal/modules/communities/<feature>`

3. Define the resource contract in `models.go`.
   Include:
   - API-facing resource struct
   - create/update input structs
   - timestamps and soft-delete fields when auditability matters

4. Implement the service in `service.go`.
   Include:
   - validation errors and internal errors as `apperror.Error`
   - sanitize helpers for create/update inputs
   - business defaults such as uppercase statuses
   - domain invariants and relationship checks

5. Implement SQL in `repository.go`.
   Include:
   - explicit column lists
   - tenant-scoped `WHERE` clauses
   - `RETURNING` for created/updated resources
   - transaction usage when multiple writes must succeed together

6. Implement handlers in `http.go`.
   Include:
   - route registration
   - JSON decoding
   - path-param injection into trusted fields
   - envelope writes via `httpx.WriteJSON`

7. Wire the module into the parent module or server.
   Common places:
   - parent module `Handler.Register(...)`
   - parent service construction if sub-services are composed
   - `apps/api/internal/platform/server/server.go` if it is a new top-level module

8. Verify the change.

## Design Rules

- Keep handlers boring.
- Never let the request body be the only source of tenant IDs.
- Translate unique constraint and not-found failures into stable error codes.
- Prefer explicit statuses over booleans for lifecycle-heavy resources.
- Use soft delete when the resource is operationally important.
- Keep the first slice thin; avoid over-engineering cross-module abstractions.

## Suggested Route Shapes

Follow existing patterns:

- Organization-scoped:
  - `POST /v1/organizations/{id}/<resources>`
  - `GET /v1/organizations/{id}/<resources>`
  - `GET /v1/organizations/{id}/<resources>/{resourceID}`

- Community-scoped:
  - `POST /v1/organizations/{organizationID}/communities/{id}/<resources>`
  - `GET /v1/organizations/{organizationID}/communities/{id}/<resources>`
  - `GET /v1/organizations/{organizationID}/communities/{id}/<resources>/{resourceID}`

Match established naming in nearby modules unless you are deliberately cleaning up a local inconsistency.

## Verification Commands

Run the smallest relevant checks after implementation:

```sh
cd apps/api && go test ./...
```

If the feature includes a schema change, also ensure the migration SQL is present in `database/` and consistent with existing schema notes.

## Done Criteria

- Files are placed in the right module path.
- Routes are registered and reachable from the current server setup.
- Queries are tenant-safe.
- Errors use `apperror` with stable codes.
- The module compiles with `go test ./...`.
