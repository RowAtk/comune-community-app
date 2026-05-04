# Implementation Checklist

Use this file as the working checklist for current repository progress. Update checkboxes when work is completed so the checklist stays accurate.

## Current Completed Foundation

- [x] Authentication: signup and login
- [x] Organization onboarding and organization membership invites
- [x] Community creation and community membership invites
- [x] Unit management
- [x] Household management
- [x] Resident management
- [x] Resident invitation and claim flow
- [x] Basic app shell and dashboard structure
- [x] Invoicing positioned as a standalone-oriented module boundary in docs and schema guidance
- [x] Numbered migration workflow added for ongoing schema changes
- [x] Primer prompt separated into its own document

## Active Documentation And Workflow Tasks

- [x] Create a dedicated implementation checklist separate from the product roadmap
- [x] Add governance rules for keeping the primer prompt updated as governance docs change
- [x] Add governance rules for keeping seed data aligned with schema and active features
- [x] Add a database-specific AGENTS file

## Seed Data Strategy

- [x] Seed only the tables that are actively used by the current application flows
- [x] Keep seeded data medium-sized and realistic rather than minimal or uniform
- [x] Include at least 20 rows overall for each actively seeded entity
- [x] Keep the first seeded tenancy easy to use for manual local testing
- [x] Avoid populating unused future-facing tables until the application actually depends on them

## Next Feature Track

- [x] Create the `apps/api/internal/modules/invoicing` module
- [x] Add invoicing routes in the API and wire them into the server
- [x] Add dedicated invoicing routes in the web app under the organization/community hierarchy
- [x] Implement recurring invoice plan management
- [x] Implement manual invoice creation
- [x] Let community admins see overdue households for maintenance invoices
- [x] Let residents see the next maintenance due date and highlight overdue invoices
- [ ] Implement proof-of-payment submission and admin review flow
- [x] Expand seed data to include invoicing tables once invoicing flows are actually wired into the app
