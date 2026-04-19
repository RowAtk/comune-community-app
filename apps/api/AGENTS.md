# AGENTS.md

## Scope

These rules apply to `apps/api` and all Go backend code beneath it.

This API is the source of truth for tenant safety, authorization, business rules, and data integrity. Favor correctness and explicitness over convenience.

## Backend Design Principles

- Keep handlers thin, services authoritative, and repositories focused on persistence.
- Put business validation in services, not only in HTTP handlers.
- Keep transport concerns out of repositories.
- Prefer domain modules under `internal/modules/...` over generic utility packages.
- Use `internal/platform/...` only for true platform concerns such as config, database helpers, HTTP adapters, telemetry, or reusable error handling.
- Avoid circular dependencies between modules. If two modules are pulling on each other, extract a clearer boundary instead of layering hacks on top.

## Recommended Module Shape

For most new modules, follow the current structure:

- `models.go` for request/response/domain structs
- `service.go` for validation, orchestration, and business rules
- `repository.go` for SQL access
- `http.go` for route registration and HTTP handling

This is the default, not a straitjacket. Add files when a module grows, but keep the separation of concerns intact.

## HTTP And Transport Rules

- Register routes centrally within the module handler.
- Continue using the shared `httpx` adapter/envelope/error flow unless there is a compelling repo-wide reason to change it.
- Parse and validate request bodies explicitly.
- Derive tenant and resource IDs from trusted route params and authenticated context, not from untrusted JSON bodies alone.
- Return consistent status codes and machine-readable error codes.
- Do not let handlers accumulate business branching or direct SQL logic.
- Keep authentication and authorization enforcement server-side for every protected route.

## Service Layer Rules

- Services own input sanitization, invariant enforcement, and orchestration across repositories/modules.
- Trim strings, normalize statuses, and apply defaults before persistence.
- Prefer explicit validation errors over vague internal failures.
- Wrap lower-level failures with operation context, but keep user-facing messages safe and stable.
- Use app-level typed errors for validation, conflict, not-found, unauthorized, and internal failures.
- Keep service methods aligned with business actions: `Create`, `List`, `GetByID`, `Update`, `Delete`, `AcceptInvitation`, etc.

## Repository And SQL Rules

- Repositories should contain SQL and row-mapping, not policy decisions.
- Every query for operational data must be tenant-scoped. For community resources, scope by both `organization_id` and `community_id` where applicable.
- Prefer explicit `INSERT`, `SELECT`, `UPDATE`, and `RETURNING` column lists.
- Use `COALESCE`, `NULLIF`, and SQL constraints intentionally, not as band-aids for unclear semantics.
- Keep query text readable and local to the function unless shared reuse is clear and stable.
- Use transactions when a business operation must succeed atomically across multiple statements.
- Detect unique and constraint violations and translate them into domain errors.
- Avoid hidden N+1 patterns. Join or batch where the business view naturally needs related data.

## Data Modeling Rules

- IDs should remain opaque to clients.
- Preserve explicit status fields and constrained enums for lifecycle-heavy entities.
- Keep audit columns meaningful: `created_at`, `updated_at`, `deleted_at`, `accepted_at`, `joined_at`, and similar timestamps should reflect real business events.
- Default to soft delete for business entities unless the data is truly disposable.
- JSON fields are for flexible settings and metadata, not for hiding first-class domain relationships.

## Auth, Permissions, And Tenant Safety

- Assume hostile input. Validate authentication, organization membership, and community-level permissions in backend code.
- Authorization checks belong near the business action, not only in middleware.
- Never expose or mutate data across organizations due to missing predicates.
- Invitations, sessions, membership changes, and future visitor/security actions must be designed with replay, expiry, and misuse in mind.
- Cookie and session security settings should be environment-aware before production deployment.

## Errors And Logging

- Keep error codes stable once consumed by the frontend.
- Internal errors may wrap detailed causes; external messages should stay safe and user-appropriate.
- Log enough structured context to trace failures by tenant, module, route, and action.
- Do not log secrets, passwords, tokens, or raw session material.

## Concurrency, Context, And Performance

- Accept and pass `context.Context` through I/O boundaries.
- Respect request cancellation and timeouts.
- Prefer straightforward performance wins: correct indexes, bounded queries, and avoiding repeated round-trips.
- Do not optimize with caching until the ownership, invalidation, and tenant boundaries are clear.

## Testing Expectations

- Add table-driven tests for validation-heavy service logic when possible.
- Add repository tests for tricky SQL behavior, unique constraints, soft-delete handling, and tenant scoping.
- Prioritize tests for authorization paths, invitation acceptance, membership changes, and cross-tenant isolation.
- Bug fixes should include a regression test if the surrounding area has test scaffolding.

## Change Checklist

- Does every query include the correct tenant scope?
- Are handler, service, and repository responsibilities still clean?
- Are errors mapped to stable codes and correct HTTP statuses?
- Are writes atomic when the business action spans multiple records?
- Does the change preserve auditability and safe defaults?
