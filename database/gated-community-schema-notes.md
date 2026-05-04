# Gated Community SaaS Schema Notes

## Overview

This schema is built around an **organization-first multi-tenant model**.

Hierarchy:

- **organizations** = the real SaaS customer and billing owner
- **communities** = gated communities managed by that organization
- **users** = global login identities
- **organization_users** = org-level access
- **community_users** = community-level access

This supports the likely commercial model where a property management organization or homeowners association is the paying customer, and that organization can manage one or many communities.

---

## Cross-cutting design decisions

### Why `organization_id` exists on many child tables
Even when a table already references `community_id`, `organization_id` is still stored explicitly for:

- easier filtering in repositories
- cleaner authorization checks
- simpler row-level security policies
- better indexing for tenant-scoped queries
- reduced joins in common queries

This is intentional denormalization for operational clarity.

### Why `users` and `residents` are separate
A user is an authentication identity.  
A resident is a business/domain entity.

A resident may not have an account yet, and a user may not be a resident.

Examples:
- resident child/dependent without a login
- security guard with a login but not a resident
- resident invited now, activates account later

### Why UUIDs are used
UUIDs are used for nearly all primary keys because they are safer for APIs and better for distributed systems than integer IDs.

### Why statuses use `VARCHAR + CHECK`
For an MVP and early production system, `VARCHAR + CHECK` is usually easier to evolve than PostgreSQL enums. It keeps constraints strong without making future value changes awkward.

### Why there are two billing layers
There are two different billing concerns:

- **subscriptions**: the organization pays your SaaS platform
- **invoicing**: residents/units pay community-related fees inside the application

They are intentionally separate.

### Why `settings` is JSONB
`settings` is used for flexible configuration that does not yet justify dedicated columns or dedicated tables. It should not become a dumping ground for core relational data.

---

## Table-by-table explanation

## 1. `organizations`

Top-level SaaS tenant and billing owner.

Columns:
- `id`: primary key
- `name`: display name of the customer account
- `slug`: unique routing-friendly identifier
- `legal_name`: optional formal/legal business name
- `billing_email`: organization billing contact
- `phone`: general contact number
- `country_code`: defaulted to `JM`
- `timezone`: org default timezone
- `status`: lifecycle of the organization account
- `settings`: organization-wide config blob
- `created_at`, `updated_at`: audit timestamps
- `deleted_at`: soft delete marker

Important constraints:
- `slug` is globally unique
- `status` is constrained to valid organization states

---

## 2. `subscriptions`

Represents the organization's SaaS subscription to your platform.

Columns:
- `id`: primary key
- `organization_id`: owning organization
- `plan_code`: internal plan identifier
- `status`: subscription lifecycle
- `billing_cycle`: monthly/quarterly/yearly
- `current_period_start`, `current_period_end`: current billing window
- `cancel_at_period_end`: marks non-renewing state
- `external_customer_id`: external provider customer ID
- `external_subscription_id`: external provider subscription ID
- `metadata`: provider-specific data
- `created_at`, `updated_at`: audit timestamps

Notes:
- This table is for **your SaaS billing**, not resident community invoices.

---

## 3. `communities`

Represents communities/properties managed by an organization.

Columns:
- `id`: primary key
- `organization_id`: owning organization
- `name`: display name
- `slug`: organization-scoped slug
- `address`: optional address
- `timezone`: community-specific timezone
- `status`: lifecycle state
- `settings`: community-specific config
- `created_at`, `updated_at`, `deleted_at`: audit/soft delete fields

Constraints:
- `(organization_id, slug)` unique
- `(organization_id, name)` unique

This means an organization cannot create duplicate community names/slugs, but different organizations may use the same names.

---

## 4. `users`

Global login identities.

Columns:
- `id`: primary key
- `email`: login email
- `password_hash`: password hash
- `first_name`, `last_name`, `phone`: profile fields
- `is_active`: auth enable/disable flag
- `last_login_at`: last successful login
- `created_at`, `updated_at`, `deleted_at`: audit/soft delete fields

Notes:
- `email` is globally unique
- This table is intentionally not resident-specific

---

## 5. `organization_users`

Organization-level membership and authorization.

Columns:
- `organization_id`: organization being joined
- `user_id`: member user
- `role`: organization-level role
- `status`: membership state
- `invited_by`: inviter user
- `joined_at`: join timestamp
- `created_at`, `updated_at`: audit fields

Common roles:
- `OWNER`
- `ORG_ADMIN`
- `BILLING_ADMIN`
- `OPERATIONS`
- `SUPPORT`

Why separate this from community memberships:
- an org admin or billing admin may not belong to every community directly
- organization permissions and community permissions are different scopes

---

## 6. `community_users`

Community-level membership and authorization.

Columns:
- `organization_id`: tenant owner
- `community_id`: target community
- `user_id`: member user
- `role`: community-level role
- `status`: membership state
- `invited_by`: inviter
- `joined_at`: join timestamp
- `created_at`, `updated_at`: audit fields

Common roles:
- `COMMUNITY_ADMIN`
- `RESIDENT`
- `SECURITY`
- `MANAGER`
- `BOARD_MEMBER`
- `MAINTENANCE_STAFF`

Notes:
- Primary key is `(community_id, user_id)`
- `organization_id` is repeated deliberately for tenant scoping

---

## 7. Invitation tables

### `organization_invitations`
Invites someone into an organization.

Columns:
- `id`
- `organization_id`
- `email`
- `role`
- `token`
- `expires_at`
- `accepted_at`
- `invited_by`
- `created_at`

### `community_invitations`
Invites someone into a specific community.

Columns:
- `id`
- `organization_id`
- `community_id`
- `email`
- `role`
- `token`
- `expires_at`
- `accepted_at`
- `invited_by`
- `created_at`

Notes:
- `token` should be used in invite acceptance flows
- `accepted_at` lets you retain invite history rather than deleting accepted invites

---

## 8. `units`

Represents physical units/homes/apartments inside a community.

Columns:
- `id`
- `organization_id`
- `community_id`
- `unit_number`
- `block_floor`
- `unit_type`
- `status`
- `created_at`, `updated_at`, `deleted_at`

Constraints:
- `(community_id, unit_number)` unique

Notes:
- `unit_number` uniqueness is scoped to a community, not globally

---

## 9. `households`

Groups multiple residents within a unit.

Columns:
- `id`
- `organization_id`
- `community_id`
- `unit_id`
- `name`
- `created_at`, `updated_at`

Purpose:
- allows a family or shared household grouping
- supports a single primary contact constraint at household level

---

## 10. `residents`

Business/domain representation of people connected to a unit.

Columns:
- `id`
- `organization_id`
- `community_id`
- `unit_id`
- `user_id`: optional login linkage
- `household_id`: optional grouping
- `first_name`, `last_name`, `email`, `phone`
- `resident_type`
- `status`
- `is_primary_contact`
- `move_in_date`, `move_out_date`
- `created_at`, `updated_at`, `deleted_at`

Important notes:
- `user_id` is nullable because not every resident needs an account
- `is_primary_contact` is constrained so only one active primary contact exists per household
- move dates are validated so move-out cannot be before move-in

Typical `resident_type` values:
- `OWNER`
- `TENANT`
- `DEPENDENT`
- `OCCUPANT`

---

## 11. `vendors`

Known service providers for the community.

Columns:
- `id`
- `organization_id`
- `community_id`
- `company_name`
- `category`
- `contact_name`
- `email`
- `phone`
- `address`
- `status`
- `rating`
- `notes`
- `created_at`, `updated_at`

Notes:
- Useful for maintenance assignment and vendor access tracking
- `rating` is constrained from 0 to 5

---

## 12. `visitors`

Visitor invitation/pass records for gate access.

Columns:
- `id`
- `organization_id`
- `community_id`
- `unit_id`
- `name`
- `phone`
- `vehicle_number`
- `expected_arrival`
- `invite_code`
- `purpose`
- `status`
- `qr_token`
- `created_by`
- `created_at`, `updated_at`
- `expires_at`

Current modeling choice:
- this table combines visitor info and pass/invite info
- that is okay for MVP

Later evolution:
- `visitors`
- `visitor_passes`
- `visitor_access_events`

Why both `invite_code` and `qr_token`:
- `invite_code` is easier for manual guard lookup
- `qr_token` supports scan-based validation

---

## 13. `security_logs`

Tracks entry/exit records or security-related access logging.

Columns:
- `id`
- `organization_id`
- `community_id`
- `person_type`
- `visitor_id`
- `resident_id`
- `vendor_id`
- `unit_id`
- `time_in`
- `time_out`
- `guard_id`
- `remarks`

Notes:
- `person_type` indicates which linked entity is expected
- `time_out` must not be earlier than `time_in`

Later evolution:
- split access events from incident/security reporting

---

## 14. `maintenance_requests`

Represents maintenance tickets/issues.

Columns:
- `id`
- `organization_id`
- `community_id`
- `unit_id`
- `resident_id`
- `vendor_id`
- `assigned_to_user_id`
- `title`
- `description`
- `category`
- `status`
- `priority`
- `scheduled_date`
- `attachments`
- `created_at`, `updated_at`

Notes:
- `attachments` is a JSONB array for MVP convenience
- later, this may become dedicated attachment tables
- supports both internal assignment and external vendor assignment

---

## 15. Operational invoicing tables

These tables should be treated as the core of a standalone-oriented `invoicing` module inside the monolith. They should remain separate from `maintenance_requests`, even when the first charge type is maintenance-related.

### `invoice_plans`
Represents recurring monthly invoicing rules configured by community admins.

Columns:
- `id`
- `organization_id`
- `community_id`
- `name`
- `plan_type`
- `status`
- `issue_day_of_month`
- `due_day_of_month`
- `default_amount`
- `starts_on`
- `ends_on`
- `description`
- `settings`
- `created_by`
- `created_at`, `updated_at`, `deleted_at`

Notes:
- This stores the recurring invoicing rule, not the generated invoice
- v1 should treat this as a monthly maintenance billing plan
- `issue_day_of_month` and `due_day_of_month` are intentionally constrained to `1..28`
- `default_amount` is the fallback amount unless a unit override exists

### `invoice_plan_unit_overrides`
Represents per-unit pricing overrides for a recurring invoicing plan.

Columns:
- `invoice_plan_id`
- `organization_id`
- `community_id`
- `unit_id`
- `amount`
- `created_at`, `updated_at`

Notes:
- Keeps per-unit pricing relational and queryable
- One override exists per plan and unit at most
- If no override exists, generated invoices should use the plan `default_amount`

### `invoices`
Represents actual invoicing records and charges to units/residents within the community.

Columns:
- `id`
- `organization_id`
- `community_id`
- `unit_id`
- `invoice_plan_id`
- `source`
- `amount`
- `paid_amount`
- `issued_on`
- `due_date`
- `status`
- `billing_period`
- `description`
- `created_at`, `updated_at`

Notes:
- `amount` and `paid_amount` are validated
- `billing_period` is a simple label like `2026-04`
- `source` distinguishes manual invoices from scheduled invoices
- scheduled invoices can reference `invoice_plan_id`
- a uniqueness rule should prevent duplicate scheduled invoices for the same plan, unit, and billing period
- these records should remain owned by the invoicing domain, not by maintenance request workflows

### `payments`
Represents payments recorded against invoices.

Columns:
- `id`
- `organization_id`
- `community_id`
- `invoice_id`
- `amount`
- `status`
- `payment_method`
- `transaction_id`
- `submitted_by_user_id`
- `recorded_by_user_id`
- `notes`
- `evidence`
- `reviewed_at`
- `payment_date`
- `created_at`

Notes:
- Supports both resident proof-of-payment submission and admin-recorded offline payments
- Pending or rejected records should not affect invoice balances until approved or recorded
- payments belong to the invoicing domain even when the charge being paid is a maintenance fee
- later, for stronger accounting flexibility, use:
  - `invoice_items`
  - `payment_allocations`

---

## 16. Communication tables

### `announcements`
Community-wide or segment-based notices.

Columns:
- `id`
- `organization_id`
- `community_id`
- `author_id`
- `title`
- `content`
- `audience_type`
- `priority`
- `expires_at`
- `created_at`

### `notifications`
User-targeted app notifications.

Columns:
- `id`
- `organization_id`
- `community_id`
- `user_id`
- `title`
- `message`
- `is_read`
- `link_to_resource`
- `created_at`

Notes:
- `announcements` are broad communications
- `notifications` are direct per-user records

---

## 17. `audit_logs`

Append-only history of important changes/actions.

Columns:
- `id`
- `organization_id`
- `actor_user_id`
- `entity_type`
- `entity_id`
- `action`
- `old_values`
- `new_values`
- `created_at`

Purpose:
- supports accountability
- useful for admin actions, financial actions, access changes, etc.

Recommended uses:
- resident updated
- invoice created
- payment posted
- visitor approved
- role changed

---

## 18. `outbox_events`

Supports reliable asynchronous event publishing.

Columns:
- `id`
- `organization_id`
- `event_type`
- `aggregate_type`
- `aggregate_id`
- `payload`
- `status`
- `available_at`
- `processed_at`
- `created_at`

Purpose:
- store events in the same transaction as business data
- background workers later process these events safely

Examples:
- send notification email
- push websocket event
- publish webhook
- queue QR generation

---

## Row-Level Security note

The schema includes an example RLS policy for `announcements`.

RLS is optional at first. A common rollout is:

1. enforce tenant/community scoping in Go repositories and services
2. later add RLS as defense in depth

If used, your app must set session or transaction-local values such as:
- `app.current_organization_id`
- `app.current_community_id`

before querying protected tables.

---

## Migration note

Even if you postpone formal migrations until the MVP stabilizes, design the schema *as if* it will be migrated cleanly later:
- explicit constraint names
- additive changes first
- avoid destructive rewrites
- keep SQL as the source of truth

When you are ready, this schema can be split into numbered migration files cleanly.
