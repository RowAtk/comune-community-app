# Data-Dense Screens

Use this file for tables, long lists, audit-style screens, and list-detail workflows.

This is the Comune-specific guide for making operational screens feel fast, readable, and human-designed instead of heavy, generic, or AI-generated.

## Purpose

Comune will increasingly include:

- resident directories
- household and unit lists
- invoice records
- payment history
- security or visitor logs
- maintenance queues

These screens will often contain more information than setup or dashboard pages. The goal is not minimalism at all costs. The goal is controlled density.

## Core Principles

Dense screens should optimize for:

- scanability
- prioritization
- predictable row structure
- low memory burden
- fast drill-down

Use these ideas especially:

- Hick's Law: do not overwhelm with too many equal choices at once
- chunking: break data into meaningful groups
- law of proximity and common region: visually group related values
- recognition over recall: make filters, statuses, and row meaning visible
- aesthetic-usability: dense screens still need polish to feel usable

## When To Use A Table

Use a table when users need to:

- compare many rows across shared attributes
- scan for changes in status
- sort or filter repeatedly
- review structured operational history

Good table use cases:

- invoices
- payments
- visitor logs
- resident membership rosters
- maintenance queue backlogs

Do not use a table when:

- each item needs a lot of narrative context
- rows are highly irregular
- mobile use would collapse the table into unreadable noise

## When To Use A Card List Instead

Use cards or grouped list items when:

- the record has one dominant identity and 2-4 supporting facts
- actions matter more than side-by-side comparison
- mobile readability is more important than high-volume comparison

Good card-list use cases:

- households
- residents grouped by household
- small invoice summaries
- setup-phase resources

## Table Structure Rules

Every table should answer these quickly:

- what is each row?
- why does it matter?
- what needs attention?
- what can I do with it?

Preferred column order:

1. human-friendly identifier
2. operational context
3. current status
4. time or money signals
5. row action

Avoid:

- starting with internal IDs
- burying the status near the far edge
- forcing users to zig-zag to understand row meaning

## Row Design

Rows should be scannable in under a second.

Rules:

- Make the first cell the anchor.
- Use a clear visual difference between primary and secondary text.
- Keep repeated metadata styles consistent.
- Let status badges do real semantic work.
- Keep row actions in a stable place.

If a row needs too much text to explain itself, it probably wants:

- grouping
- an expanded detail drawer
- or a separate detail page

## Column Rules

Every column must earn its space.

Keep:

- identifiers
- status
- urgency signals
- dates that affect action
- money values when relevant

Question heavily:

- decorative icons with no meaning
- duplicate columns
- low-value timestamps
- verbose descriptions repeated on every row

## Sorting And Filtering

Dense screens become human-friendly when they narrow quickly.

Rules:

- Put the highest-value filters first.
- Default sort should match the operational job.
- Status filters should be easy to reach.
- Date filters matter for histories and logs.
- Search should work on human terms, not just exact IDs.

Examples:

- invoices: status, billing period, due date
- residents: name, household, unit, status
- logs: date range, event type, person type

## Priority Cues

Not every row should shout.

Use stronger emphasis only for:

- overdue
- blocked
- failed
- pending approval
- unusually recent or urgent activity

Methods:

- status color
- stronger text weight
- pinned first column cues
- row highlight bands for truly urgent states

Do not:

- use warning styling on half the table
- apply equal emphasis everywhere
- encode importance using color alone

## Grouped Dense Screens

Sometimes grouped sections beat one giant table.

Use grouped dense layouts for:

- residents by household
- invoices by status
- maintenance items by queue stage
- logs grouped by day

Rules:

- group labels must carry real meaning
- counts are useful when groups are large
- keep the group style consistent across modules

## List-Detail Pattern

Use list-detail when operators need both scanning and context.

Pattern:

- left or top side for list
- right or lower side for selected detail
- stable selection state
- detail actions close to the details, not buried in the list

Use this when:

- opening a full page for every row would be too slow
- users need to review several records in sequence

Avoid list-detail when:

- the detail panel becomes a second full application
- the list loses enough space that scanning suffers badly

## Mobile Behavior

Dense screens must not simply shrink.

Rules:

- collapse less-important columns first
- keep the row anchor visible
- convert wide rows into stacked cells when needed
- preserve status and next action above the fold

If a table cannot adapt well on mobile:

- switch to grouped cards
- or provide a narrower default dataset for smaller screens

## Empty, Loading, And Error States

Dense screens need operational states too.

Rules:

- empty states should suggest how to populate or broaden results
- loading states should preserve layout expectation
- errors should explain whether the list failed, the filters failed, or a row action failed

## What Makes Dense Screens Feel AI-Generated

- every module uses the same generic table with no task-specific tuning
- too many columns with equal visual weight
- no clear row anchor
- filters that exist but do not reflect the actual job
- actions hidden in vague kebabs for everything
- no distinction between urgent and routine rows

## What Makes Them Feel Human-Designed

- the table structure reflects the operator's real decisions
- the first column is meaningful and stable
- sorting defaults feel intentional
- urgent rows rise naturally
- supporting details are present but quiet
- drill-down feels fast and obvious

## Comune Defaults

- tenant and community context should be obvious in dense admin screens when users can cross between scopes
- money screens should surface overdue and outstanding amounts early
- people screens should prefer names, units, and household context over IDs
- security and visitor logs should prioritize timestamp, person type, and event outcome
- maintenance queues should prioritize status, priority, and assignment

## Sources

- Laws of UX: https://lawsofux.com/laws/
- NN/g usability heuristics summary: https://media.nngroup.com/media/articles/attachments/Heuristic_Summary1-compressed.pdf
- NN/g visual-design principles poster: https://media.nngroup.com/media/articles/attachments/Visual-Design-Principles-Poster.pdf
- NN/g information foraging reference: https://media.nngroup.com/media/articles/attachments/InformationForaging_SizeLetter.pdf
- Smashing Magazine, Building Better UI Designs With Layout Grids: https://shop.smashingmagazine.com/2017/12/building-better-ui-designs-layout-grids/
- Reddit discussion on what makes sites feel AI-generated: https://www.reddit.com/r/website/comments/1rfkiqm/advice_whats_meant_by_looks_ai/
- Reddit discussion on generic AI-looking websites: https://www.reddit.com/r/SaasDevelopers/comments/1r5zwro/aigenerated_websites_always_look_generic_how_do/
