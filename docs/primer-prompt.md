# Primer Prompt For Future Sessions

```text
Use the repository docs as the primary planning context for this product.

Primary references:
- docs/AGENTS.md
- docs/features.md
- docs/plan.md
- docs/implementation-checklist.md
- docs/primer-prompt.md
- docs/ui-ux/ux-design-guide.md
- docs/ui-ux/layout-patterns.md
- docs/ui-ux/visual-language.md
- docs/ui-ux/data-dense-screens.md
- notes.md
- database/AGENTS.md
- database/MIGRATIONS.md
- database/gated-community-schema-notes.md

When helping with this repo:
- Treat organizations as the true SaaS tenant boundary and communities as children of organizations.
- Preserve tenant scoping in database, API, and UI decisions.
- Assume Jamaica defaults unless a requirement says otherwise.
- Keep modules separated by domain: residents, invoicing, visitors, maintenance, security, and notifications.
- Treat `invoicing` as its own module boundary, even when the first charge type is maintenance-related.
- Use `docs/ui-ux/ux-design-guide.md` as the default UX decision framework for layout, hierarchy, grouping, and action design.
- Use `docs/ui-ux/layout-patterns.md` for page arrangement and workflow structure.
- Use `docs/ui-ux/visual-language.md` for visual hierarchy, themeability, and anti-generic UI guidance.
- Use `docs/ui-ux/data-dense-screens.md` for tables, list-detail flows, logs, audit views, and other operationally dense interfaces.
- Keep the web app themeable by using shared semantic theme tokens from `apps/web/src/lib/theme/palette.css` instead of scattering hardcoded palette values through routes and components.
- Prefer explicit lifecycle statuses, auditability, and production-credible MVP slices.
- Use the schema notes as the source of truth for entity relationships and likely table structure.
- Use numbered migrations for schema evolution and keep the schema snapshot suitable for fresh bootstrap.
- When schema or workflow changes are introduced, update migrations and keep `database/seed.sql` useful for medium-sized realistic testing.
- When new UX governance or theme reference docs are added, update this primer so future sessions keep using the right design sources.
- Before proposing new work, check what is already implemented and avoid re-planning completed foundations.

Current roadmap context:
- Foundational auth, organization/community setup, units, households, residents, and resident invite flows are already in place or far enough along to build on.
- The next priority feature is maintenance payments and invoicing without payment gateway integration.
- For that invoicing phase, prioritize recurring monthly invoice plans, manual invoice creation, manual payment entry, resident proof-of-payment submission, invoice generation, payment history, invoice balance tracking, audit-friendly records, overdue-household visibility for community admins, and resident visibility into the next maintenance due date plus overdue invoices.

When responding:
- Ground recommendations in the existing docs and schema notes.
- Suggest the smallest production-credible next slice first.
- Call out assumptions clearly.
- Keep plans modular, tenant-safe, and implementation-ready.
```
