# AGENTS.md

## Scope

These rules apply to `apps/web` and all Svelte/SvelteKit frontend code beneath it.

The web app should feel fast, calm, and operationally trustworthy for property managers, staff, and residents. Prefer server-first simplicity over front-end cleverness.

## Frontend Design Principles

- Use `docs/ui-ux/ux-design-guide.md` as the default decision framework for hierarchy, grouping, pacing, and primary-action emphasis.
- Use `docs/ui-ux/layout-patterns.md` for recurring page arrangement and screen composition decisions.
- Use `docs/ui-ux/visual-language.md` for visual language, theme tokens, and anti-generic UI rules.
- Use `docs/ui-ux/data-dense-screens.md` when building tables, long lists, list-detail views, logs, or other high-density operational screens.
- Use SvelteKit’s server capabilities first. Reach for `load` functions and form actions before adding client-only data orchestration.
- Keep route data ownership close to the route that uses it.
- Treat the API as the source of truth for permissions, statuses, and tenant-scoped data.
- Build UI that helps operators move confidently through real workflows, not just demo screens.
- Favor composable feature components over sprawling page files, but do not extract components too early.

## Routing And Data Loading

- Prefer `+page.server.ts` and `+layout.server.ts` for authenticated data fetching and mutations.
- Use shared API helpers for consistent envelope parsing, cookie forwarding, and error handling.
- Keep redirects and authentication checks in server load paths.
- Page loads should return stable shapes so Svelte components stay simple.
- Avoid duplicate fetching across nested layouts and pages when parent data can be reused cleanly.
- Use URL structure to reinforce the tenant hierarchy: organization first, then community, then module resource.

## Forms And Mutations

- Prefer SvelteKit form actions for CRUD workflows.
- Validate user input on the server and return actionable form errors.
- Preserve submitted values on failure so operators do not retype operational forms.
- Keep success messaging concise and task-oriented.
- Avoid optimistic UI for workflows that affect permissions, occupancy, invitations, billing, or security-sensitive data unless conflict handling is fully designed.

## State Management

- Prefer plain props, route data, and derived state over custom client-side state systems.
- Introduce shared client state only when multiple distant parts of the UI truly need synchronized interactive state.
- Keep ephemeral UI state local to the component that owns it.
- Do not mirror server truth into long-lived client stores without a clear invalidation strategy.

## API Integration

- Continue using shared API request utilities for browser and server calls.
- Normalize backend errors into user-facing messages without leaking internal details.
- Do not hardcode API base URLs inside feature code.
- Keep TypeScript API types aligned with backend payloads.
- When backend contracts evolve, update the frontend types and consuming routes in the same change when possible.

## Component Rules

- Prefer small, purposeful components with clear inputs.
- Shared UI primitives should stay generic; domain-specific behavior belongs in feature components or routes.
- Reuse existing UI building blocks before introducing new patterns.
- Keep accessibility intact: labels, button semantics, focus states, keyboard navigation, and readable contrast are required.
- Empty states, loading states, and error states are part of the feature, not polish to defer indefinitely.

## Svelte And TypeScript Conventions

- Use Svelte 5 idioms already present in the codebase when they improve clarity.
- Keep component scripts typed.
- Prefer derived values and straightforward functions over heavy abstraction.
- Avoid unnecessary reactive complexity.
- Keep utility functions close to the feature unless they are genuinely shared.

## Visual And Product Standards

- Preserve the existing application visual language unless there is a deliberate design update.
- Keep the app themeable. Prefer semantic theme tokens from `apps/web/src/lib/theme/palette.css` over hardcoded route-level palette values so future palette changes remain cheap.
- Interfaces should feel professional, not generic admin-template clutter.
- Design for dense operational use: fast scanning, clear hierarchy, and obvious next actions.
- Statuses, warnings, and badges should map to real business meaning consistently.
- Never use raw database IDs as the primary thing a user sees on a screen. Prefer names, emails, usernames, slugs, unit numbers, household names, or other business-friendly identifiers.
- Jamaica-specific context matters. Date, time, phone, address, and terminology choices should fit local operations when product requirements call for it.

## Performance And UX

- Keep pages fast by minimizing duplicate fetches and unnecessary client-side JavaScript.
- Prefer progressive enhancement over JS-only flows.
- Split complex pages into smaller sections when it improves scanability and maintenance.
- Large collections should gain pagination, filtering, or grouping before they become unwieldy.
- Build screens with the expectation that a single community may eventually have hundreds or thousands of residents, households, visitors, invoices, notifications, or audit records.
- Do not rely on “load everything and filter in the browser” once a resource is plausibly high-volume; move that work server-side.
- Prefer human-friendly summaries and drill-down views over rendering huge dense datasets on one page.

## Testing Expectations

- Add tests for critical UI workflows once frontend test infrastructure is introduced in this app.
- Prioritize coverage for auth redirects, form validation, invitation flows, and tenant-aware navigation.
- For now, verify route loads, form actions, and error handling manually when changing user flows.

## Change Checklist

- Is the route server-first unless there is a strong reason not to be?
- Are tenant and resource IDs sourced from route params instead of ad hoc client state?
- Does the UI preserve backend truth for permissions and statuses?
- Are form failures recoverable without forcing the user to start over?
- Does the page include useful loading, empty, and error behavior?
