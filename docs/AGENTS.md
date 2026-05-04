# AGENTS.md

## Purpose

This repository powers a multi-tenant gated community management SaaS for Jamaica. It supports organizations that manage one or more communities and operational modules such as residents, visitors, billing, maintenance, security, and notifications.

These governance rules exist to keep the codebase:

- safe for multi-tenant production use
- fast to evolve as a solo founder-led product
- easy to onboard future engineers into without major rewrites

When making changes anywhere in this repo, follow this file first, then any more specific `AGENTS.md` deeper in the tree.

## Product And Domain Rules

- Treat `organizations` as the true SaaS tenant boundary. Communities are children of organizations, not standalone tenants.
- Preserve tenant scoping in every layer: database, API, background work, and UI navigation.
- Jamaica is the default operating context unless requirements explicitly say otherwise.
- Prefer `country_code = 'JM'` and `timezone = 'America/Jamaica'` defaults where the product needs regional assumptions.
- Design modules so they can grow independently: residents, invoicing, visitors, maintenance, security, and notifications should not become tightly coupled blobs.
- If a module is a plausible future standalone product, keep its contracts and ownership boundaries clean inside the monolith from the start rather than scattering its logic across unrelated domains.
- Model household, resident, unit, and membership data as operational facts. Avoid burying core domain behavior in generic JSON settings.
- Use soft deletion for business records when auditability matters. Hard deletes should be rare and intentional.
- Make lifecycle states explicit. Use constrained statuses rather than ambiguous booleans for business state.

## Architecture Rules

- Keep the monorepo modular. Each app should own its runtime concerns and share only stable, intentional contracts.
- Prefer vertical slices by domain/module over horizontal dumping grounds.
- New backend features should usually land as a module with clear models, service logic, repository access, and transport handlers.
- New frontend features should usually land as route-owned server loads/actions plus focused reusable UI where repetition is real.
- For invoicing work specifically, use `invoicing` as the implementation boundary rather than burying invoice logic inside communities, residents, or maintenance features.
- Keep the application aligned with modern, well-supported web development standards so routine work does not accumulate avoidable technical debt.
- Avoid premature shared packages. Duplicate a small amount first, then extract only once the shape is proven.
- Favor boring, explicit code over clever abstractions.
- Every cross-layer contract should be easy to trace from UI -> API -> database.

## Scale Rules

- Design every new feature so it remains reasonable at thousands of users overall and thousands of records within a single community.
- Do not assume lists will stay small. Plan query patterns, indexes, pagination, filtering, and sorting before large datasets become a production problem.
- Avoid N+1 data access in APIs and avoid duplicate fetch patterns in the web app when related data can be loaded in one bounded pass.
- Prefer stable, indexed lookups and bounded result sets over convenience queries that scan or return everything.
- Features should degrade gracefully under growth: admin screens, dashboards, invite flows, resident lists, household lists, invoicing records, and notifications must all be able to handle materially larger tenant datasets than exist today.
- If a first version ships with a scale constraint for speed, document that constraint explicitly in code or adjacent docs rather than leaving it as a hidden assumption.

## Multi-Tenancy And Security

- Never query or mutate tenant data without organization scoping.
- Community-scoped records must carry both `organization_id` and `community_id` when that is the established data model.
- Authorization decisions must be enforced on the server, never only in the UI.
- Do not trust client-supplied tenant identifiers without validating the authenticated user has access.
- Treat invitation tokens, sessions, and future payment/security workflows as sensitive data. Log metadata, not secrets.
- Default to least privilege for new roles and permissions.
- Any feature that affects residents, visitors, security incidents, access control, or invoicing must leave an audit-friendly trail.

## Database And Schema Discipline

- PostgreSQL is a core part of the application design, not a passive storage layer.
- Prefer explicit schema changes in `database/` with clear names and safe execution order.
- Add constraints for invariants that must never be violated in production.
- Add indexes for real access paths, especially tenant-scoped queries and time-ordered listings.
- Avoid nullable columns when a default or explicit state is better.
- Keep SQL readable and reviewable. Favor explicit column lists over `SELECT *`.
- Schema changes must consider backfill, rollout safety, and impact on existing seed data.
- When schema changes or new product capabilities affect seeded workflows, update `database/seed.sql` in the same change so local development keeps working without manual repair.
- Seed data should remain useful for realistic testing, not just smoke tests. Prefer a medium-sized, operationally believable dataset that exercises tenant scoping, roles, lifecycle states, and high-value module flows.
- If a change alters a business rule, update schema notes or adjacent documentation when the rule is subtle.

## API Contract Rules

- Keep APIs resource-oriented and predictable.
- Maintain consistent response envelopes and error shapes within each app.
- Validate input at the edge and again in domain/service logic where invariants matter.
- Prefer additive API evolution. Avoid breaking existing routes or payloads without a deliberate migration plan.
- Return domain-meaningful error codes so the frontend can react without brittle string matching.
- Pagination, filtering, and sorting should be introduced before large lists become operational pain.

## Quality And Maintainability

- Optimize for code that one person can safely change at 11 PM and a team can still understand six months later.
- Keep functions small enough to scan quickly.
- Prefer explicit names that match the business language used by property managers and community admins.
- Keep dependency technical debt low: prefer actively maintained dependencies, upgrade on a steady cadence, and avoid introducing packages that duplicate platform capabilities or add long-term maintenance drag without clear product value.
- Do not expose raw database IDs in user-facing UI copy. IDs are acceptable in routes, API contracts, logs, and internal tooling, but screens should prefer names, labels, slugs, unit numbers, emails, or other human-meaningful identifiers.
- Write comments only when they explain non-obvious intent, policy, or tradeoffs.
- Remove dead code, stale TODOs, and abandoned experiments as part of normal maintenance.
- Avoid broad refactors during feature work unless they directly reduce risk or unblock the feature.

## Testing And Verification

- New business logic should be covered by automated tests once test infrastructure exists in that area.
- Prioritize tests around tenant isolation, authorization, lifecycle transitions, money movement, and invite/access flows.
- When tests do not yet exist for an area, verify changes with the smallest realistic end-to-end path and document what was checked.
- Bug fixes should add a regression test whenever practical.

## Observability And Operations

- Log with enough structure to debug production issues by tenant, route, and operation.
- Do not leak secrets, tokens, passwords, or personally sensitive data into logs.
- Prefer metrics and tracing around external boundaries, slow queries, and critical workflows.
- Make failure modes obvious. Silent partial success is usually worse than a clear error.

## Delivery Rules

- Ship the smallest production-credible slice first, then harden.
- Preserve a clean path for future extraction of shared policies, permissions, and domain services.
- If a shortcut is taken for speed, keep it local, document the constraint, and avoid turning it into a hidden platform standard.
- Any new module should start with clear ownership boundaries, tenant-safe data access, and a path to auditability.

## Documentation Expectations

- Update local documentation when architectural decisions, workflows, or domain rules materially change.
- Keep setup docs, schema notes, and environment assumptions aligned with the code.
- Keep `docs/primer-prompt.md` aligned with the current governance and workflow docs. When new governance markdown files are added or old ones become irrelevant, update the primer references and guidance in the same change.
- Keep `docs/implementation-checklist.md` current. When a tracked task is completed or intentionally deferred, update the checklist in the same change rather than letting it drift.
- Treat `docs/ui-ux/ux-design-guide.md` as an active governance document for future UI work. If the design system or workflow changes materially, update the guide references and any theme docs in the same change.
- Governance files should stay practical. If a rule no longer helps this repo move faster or safer, update it.
