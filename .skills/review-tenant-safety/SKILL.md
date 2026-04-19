---
name: review-tenant-safety
description: Review code changes in this repository for multi-tenant safety, focusing on organization and community scoping, authorization gaps, cross-tenant data leakage, unsafe logging, and consistency between Go API, PostgreSQL queries, and Svelte routes.
---

# Review Tenant Safety

Use this skill when reviewing backend, frontend, or schema changes that could break tenant isolation.

## Primary Review Lens

Find concrete risks first. This repo is organization-first multi-tenant SaaS, so the main question is:

Can one organization or community read, mutate, infer, or affect another tenant's data?

## What To Inspect

### Backend

- Missing `organization_id` predicates
- Missing `community_id` predicates for community-scoped resources
- Authorization checks that rely on UI assumptions
- Route params ignored in favor of body fields
- invite/session/token flows that can be replayed or guessed
- logs that include tokens, response bodies, or sensitive personal data

### SQL

- tables missing tenant ownership columns
- uniqueness that should be scoped by organization or community but is global
- joins that can cross tenants because the predicates are incomplete
- soft-delete predicates missing from reads or updates

### Frontend

- routes building requests with stale or user-editable tenant IDs
- pages assuming access based on navigation alone
- forms that allow hidden tenant identifiers to drift from route params
- exposing data from parent loads across the wrong organization/community context

## Repo-Specific Hotspots

Check these patterns carefully:

- `apps/api/internal/modules/**/repository.go`
- `apps/api/internal/modules/**/service.go`
- `apps/api/internal/modules/**/http.go`
- `apps/api/internal/platform/server/server.go`
- `apps/web/src/routes/**/+page.server.ts`
- `apps/web/src/routes/**/+layout.server.ts`
- `database/*.sql`

## Findings Format

Present findings ordered by severity with file references.

For each finding, state:

- what is unsafe
- how cross-tenant impact could happen
- the likely fix direction

If there are no findings, say that explicitly and call out residual risks or unreviewed areas.

## High-Value Checks

- Query reads by `id` only instead of `(organization_id, id)` or `(organization_id, community_id, id)`
- Create/update handlers trust JSON tenant IDs more than path/auth context
- Service methods skip membership/role checks on mutating actions
- Invitation acceptance is not bound to intended user/email/expiry
- Global uniqueness accidentally leaks tenant existence
- Logging captures raw response bodies containing sensitive records

## Useful Commands

```sh
rg -n "organization_id|community_id|PathValue|CurrentAuthResult|RequireAuth|token|set-cookie|response_body" apps/api apps/web database
cd apps/api && go test ./...
cd apps/web && pnpm check
```

## Done Criteria

- Findings are concrete and repo-specific.
- Cross-tenant impact is explained clearly.
- False positives are minimized by checking surrounding code before reporting.
