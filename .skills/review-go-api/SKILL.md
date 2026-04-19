---
name: review-go-api
description: Review Go API changes in apps/api for correctness, maintainability, tenant safety, HTTP contract consistency, service and repository boundaries, and realistic production risks in this repository.
---

# Review Go API

Use this skill for code review of Go backend changes in `apps/api`.

## Review Priorities

Focus on findings, not praise. Prioritize:

1. correctness bugs
2. tenant-safety regressions
3. broken API contracts
4. transaction and persistence issues
5. maintainability problems that will slow future features

## Repo-Specific Standards

- Handlers should stay thin.
- Services should own validation, defaults, and orchestration.
- Repositories should own SQL and row scanning.
- Errors should map into `apperror` and `httpx` consistently.
- Route registration should happen in module handlers.
- Top-level server wiring lives in `apps/api/internal/platform/server/server.go`.

## What To Look For

### HTTP Layer

- inconsistent status codes
- body parsing without validation
- route params not copied into trusted input fields
- direct business logic in handlers
- protected routes missing auth middleware

### Service Layer

- missing sanitization of strings and statuses
- missing validation for required IDs or business invariants
- not-found, unique, or validation failures mapped to generic internals
- business operations that should be atomic but are not transactional

### Repository Layer

- SQL missing tenant predicates
- updates/deletes missing `deleted_at IS NULL` where soft-delete is expected
- row scanning order mismatches
- uniqueness/constraint failures not handled upstream
- hidden N+1 behavior or repeated round-trips

### Platform And Ops

- request/response logging that leaks too much data
- tracing or middleware ordering issues
- insecure cookie/session defaults that should be environment-aware
- config fallbacks that are fine for dev but dangerous if shipped unchanged

## Findings Format

Report findings with:

- severity
- file reference
- concise explanation of the risk
- suggested fix direction

If no findings, state that and mention remaining gaps, such as missing tests or unreviewed runtime behavior.

## Good Commands

```sh
rg -n "httpx|apperror|QueryRow|Exec\\(|RunInTx|RequireAuth|SetCookie|Secure:" apps/api
cd apps/api && go test ./...
```

## Review Checklist

- Are handlers, services, and repositories still cleanly separated?
- Are HTTP responses and error codes consistent with neighboring modules?
- Are multi-write operations transactional?
- Are logs and cookies safe enough for the current environment assumptions?
- Does the change match existing module and route conventions?
