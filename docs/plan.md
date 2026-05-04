# Product Plan

This product is an organization-first gated community management SaaS for Jamaica. The core tenant boundary is `organizations`, with `communities` operating underneath them. The near-term goal is to ship the smallest production-credible modules in a safe order while preserving tenant scoping, auditability, and clear domain boundaries across residents, invoicing, visitors, maintenance, security, and communications.

## Current Foundation

The current repo baseline is strong enough to support the next operational module rather than restarting core setup work.

- Authentication: signup and login
- Organization onboarding
- Community creation
- Unit management
- Household management
- Resident management
- Resident invitation and claim flow
- Basic app shell and dashboard structure

This foundation is enough to begin operational invoicing because invoice and payment workflows depend on existing organization, community, unit, resident, and membership context that is already present or far enough along to build on safely.

## Recommended Feature Priority

1. **Maintenance payments, manual payment recording, and invoice generation**
   Unlock the next layer of operational value for community admins without waiting on external payment gateways. This slice should be implemented as a standalone-oriented invoicing module with recurring monthly invoice generation from community-admin plans, invoice list and detail views, manual invoice creation for ad hoc charges, manual payment entry, payment history, balance tracking, and invoice status updates. It should support community default pricing with per-unit overrides, separate issue and due days, resident proof-of-payment submission for admin review, and admin-recorded offline payments such as cash, cheque, or bank transfer. This remains an operational invoicing feature, not a payment gateway feature.

2. **Billing audit trail and accounting visibility**
   Treat this as part of the invoicing module or the immediate follow-up slice. Invoice creation, payment posting, approval, reversal or adjustment actions, and status changes should leave an audit-friendly history. Prefer a simple historical ledger view over premature accounting-engine complexity.

3. **Visitor access and guard verification**
   This is a high-fit module for gated communities and should be delivered as one coherent operational slice: visitor pass creation, lookup by invite code, guard verification, and entry or exit recording. Organization and community scoping plus guard-role authorization must remain strict.

4. **Security logs and operational gate history**
   Extend visitor verification into searchable entry and exit history for community admins and security staff. Keep person types explicit, especially visitor, resident, and vendor, so the module can evolve cleanly.

5. **Maintenance requests**
   Add resident and admin maintenance ticket workflows after billing and access-control flows are stable. Include statuses, priorities, assignment, and vendor linkage, while avoiding unnecessary complexity in scheduling and attachment handling for the MVP.

6. **Announcements and notifications**
   Add broad communications after the core operational workflows exist. Keep community-wide announcements separate from direct user notifications and prioritize useful delivery over deep automation.

7. **Advanced permissions and household or community rules**
   This includes ideas already noted for future work, such as household admins and community-rule-driven action limits. Leave this until after the operational modules above because it depends on stronger policy, lifecycle, and audit foundations.

## Next Feature: Maintenance Payments And Invoicing

**Goal**

Let community admins manage maintenance invoicing without online payment integration.

**In scope**

- Create recurring community invoicing plans for maintenance fees and dues
- Generate invoice records monthly from active invoicing plans
- Support a community default amount with per-unit override pricing
- Support separate invoice issue day and payment due day
- Create invoices manually for units and residents when admins need ad hoc charges
- View invoice lists and invoice detail
- Record manual payments
- Track `amount`, `paid_amount`, `due_date`, `billing_period`, and `status`
- Keep historical payment records
- Allow resident proof-of-payment submission
- Allow admin approval, rejection, and recording of offline payments
- Surface payment method and reference or transaction details when relevant
- Add audit-friendly history for invoicing actions

**Out of scope**

- Card processing
- Automated gateway webhooks
- SaaS subscription billing for the platform itself
- A full accounting engine or multi-allocation ledger design beyond MVP needs

**Suggested MVP lifecycle**

- Invoicing plan created or updated
- Scheduled or admin-triggered generation creates monthly invoices
- Invoice created
- Invoice sent or published
- Payment submitted or recorded
- Payment approved if review is required
- Invoice becomes partial or paid
- Invoicing history remains visible and sufficiently immutable for auditability

**Data alignment**

Use the existing schema direction around the invoicing module:

- `invoice_plans`
- `invoice_plan_unit_overrides`
- `invoices`
- `payments`
- `audit_logs`
- organization, community, and unit scoping on all invoicing records

## Delivery Guidance

- Ship `invoicing` as its own top-level module rather than mixing it into residents, communities, or maintenance requests
- Keep recurring plan configuration separate from generated invoices so schedule rules and invoicing records evolve independently
- Preserve organization and community scoping on every invoicing query and mutation
- Keep invoicing generic enough to support future non-maintenance charge types even if v1 uses `MAINTENANCE`
- Treat recurring plans, invoice generation, invoice records, and payment review as invoicing-owned capabilities
- Prefer explicit statuses over booleans
- Favor append-only history for important money actions
- Default Jamaica assumptions where regional defaults matter
- Keep the first release production-credible rather than feature-complete
- Keep the module boundary strong enough that invoicing can later become a separate application or API without reworking unrelated domains

## Near-Term Success Criteria

- An org admin can create and manage communities, units, households, and residents
- A community admin can issue invoices
- A community admin can configure recurring monthly maintenance invoicing with default and per-unit pricing
- A resident or admin can record payment evidence or manual payments
- Admins can see invoice balances and payment history
- Invoicing actions are tenant-safe and audit-friendly
- The next operational modules can build on the same tenant and data-model foundations

The active implementation tracker for completed and pending work lives in [implementation-checklist.md](/home/rowana/projects/gated-community/comune/docs/implementation-checklist.md).
