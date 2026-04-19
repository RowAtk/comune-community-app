---
name: generate-svelte-crud-page
description: Generate or extend a CRUD page in apps/web for this repository's SvelteKit frontend when building organization or community resource pages that use server loads, form actions, shared API helpers, and existing UI components.
---

# Generate Svelte CRUD Page

Use this skill when creating a new SvelteKit CRUD page or extending an existing resource page in `apps/web`.

## Goal

Build pages that match the current frontend pattern:

- server-first data loading with `+page.server.ts`
- form actions for create flows
- shared API calls through `$lib/server/api`
- typed payloads from `$lib/api/types`
- reusable UI via existing components and `ResourceCollection`

## Repo Conventions To Follow

- Authenticated app routes live under `src/routes/(app)/...`.
- Tenant hierarchy is encoded in the URL:
  - organization
  - community
  - resource
- Most current CRUD pages pair:
  - `+page.server.ts` for load/actions
  - `+page.svelte` for presentation
- Shared helpers:
  - `$lib/server/api`
  - `$lib/api/http`
  - `$lib/api/types`
- Existing reusable UI:
  - `$lib/features/community-resources/resource-collection/*`
  - `$lib/components/ui/*`

## Workflow

1. Find the closest existing page.
   Best references:
   - `src/routes/(app)/organizations/[organizationId]/communities/[communityId]/units`
   - `src/routes/(app)/organizations/[organizationId]/communities/[communityId]/households`
   - `src/routes/(app)/organizations/[organizationId]/communities/[communityId]/residents`

2. Start with the server file.
   In `+page.server.ts`:
   - fetch route-owned data with `apiServerRequest`
   - return stable fallback shapes on error
   - implement `actions.create`
   - validate form inputs server-side
   - preserve submitted values with `fail(...)`

3. Build the page UI in `+page.svelte`.
   - render the resource list
   - use `ResourceCollection` when the feature is list-plus-create
   - use shared UI primitives for labels, inputs, buttons, badges, cards
   - include useful empty and error states

4. Update API types if the backend contract changed.
   Usual file:
   - `src/lib/api/types.ts`

5. Add small helpers only if reused or meaningfully clarifying.

## Design Rules

- Prefer server data over client-managed fetch state.
- Keep tenant IDs sourced from route params, not local stores.
- Keep form error recovery smooth by preserving values.
- Do not add heavy client-side state libraries.
- Avoid optimistic UI for permissions, occupancy, invitation, billing, or security-sensitive flows.

## Typical Server Pattern

Use the existing shape:

- `load` fetches resource data and returns a fallback empty collection plus `apiError` on failure
- `actions.create`:
  - reads `FormData`
  - trims and normalizes values
  - validates required fields
  - posts JSON to the API
  - returns `createSuccess` or `fail(...)` with `createError` and `values`

## Typical UI Pattern

- derive route base paths from `$app/state`
- group or label resources in the page rather than adding premature shared abstractions
- map status fields to badge variants
- explain the operational meaning of the resource in the page copy

## Verification Commands

```sh
cd apps/web && pnpm check
cd apps/web && pnpm lint
```

If the page depends on new backend fields, also run:

```sh
cd apps/api && go test ./...
```

## Done Criteria

- The route follows existing server-first CRUD patterns.
- Form failures preserve user input.
- API calls use shared helpers.
- Types and UI match the backend contract.
- `pnpm check` passes.
