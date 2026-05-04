# Layout Patterns

Use this file as the practical layout companion to [ux-design-guide.md](/home/rowana/projects/gated-community/comune/docs/ui-ux/ux-design-guide.md).

The UX guide covers principles. This file turns those principles into repeatable page structures for Comune.

## Purpose

Comune is an operational product, not a marketing site. Most screens should optimize for:

- fast scanning
- low cognitive load
- obvious next actions
- strong grouping
- stable structure across modules

This guidance is shaped by:

- Laws of UX on Hick's Law, chunking, common region, proximity, and aesthetic-usability
- NN/g guidance on visual hierarchy, balance, contrast, and minimalist design
- practitioner discussion from Reddit about what makes interfaces feel generic or "AI-generated"

## Default Page Order

For most authenticated product screens, prefer this order:

1. Context header
2. Primary actions
3. Status or summary row
4. Main operational content
5. Secondary tools or configuration

This keeps the first screenful answerable:

- Where am I?
- What can I do next?
- What needs attention?

## Header Pattern

Every major page should begin with:

- a small context label or kicker
- one clear page title
- one short explanation sentence
- one compact action group

Rules:

- Keep the description to one thought.
- Do not stack multiple competing messages in the header.
- Primary action first, secondary action second.
- If there is no meaningful action, do not invent one just to fill the space.

## Summary Row Pattern

Use summary cards only when they help someone decide what to do next.

Good summary metrics:

- overdue households
- unpaid invoices
- next maintenance due date
- active communities
- residents awaiting claim or approval

Avoid summary rows when:

- the numbers do not change the next action
- the cards only restate visible list data
- they exist only to make the page feel "full"

Rules:

- Prefer 2-4 cards, not 5-8.
- The most urgent card should be visually strongest.
- At least one card should usually answer "what needs attention?"

## Main Content Pattern

Choose one dominant content mode per page:

- list
- grouped list
- form
- detail view
- split view

Do not make a page equally about all five.

Guidance:

- Lists are for scanning many records quickly.
- Grouped lists are for related operational objects, such as residents by household.
- Detail views are for one record and its related actions.
- Split views are useful when one side is clearly secondary, such as settings beside a primary list.

## List Pattern

Use lists for operational resources like units, households, residents, and invoices.

Rules:

- Each row or card must expose one human-friendly identifier first.
- Show the most decision-relevant metadata next.
- Status should be visible without opening the detail view.
- The row action should be obvious and placed consistently.

For high-volume lists:

- add grouping, filtering, or pagination before the page becomes noisy
- do not rely on a long undifferentiated stack of cards forever

## Grouping Pattern

Prefer grouping when the user naturally thinks in clusters:

- residents by household
- households by unit
- invoices by status or billing period
- communities by organization

Rules:

- Group labels must be meaningful, not decorative.
- Keep groups collapsible only if there are enough items to justify it.
- Use spacing and boundary, not just color, to signal grouping.

## Form Pattern

Operational forms should feel short even when they are not small.

Rules:

- Group fields into small clusters.
- Put the most important fields first.
- Do not mix unrelated settings in one uninterrupted block.
- Use helpful defaults where safe.
- Error messages should appear close to the problem and preserve entered values.

Preferred order:

1. core identity fields
2. scope or assignment fields
3. schedule or money fields
4. optional notes
5. submission actions

## Empty States

Empty states should do one job: help the user start.

Rules:

- State what is missing.
- Explain why it matters.
- Offer one clear next action.

Avoid:

- vague "nothing here yet" copy
- multiple competing calls to action
- decorative filler illustrations without operational value

## Error States

Errors should reduce confusion, not merely announce failure.

Rules:

- describe the failed action
- keep blame out of the language
- preserve progress whenever possible
- present the next safe action

## Arrangement Rules That Fight The "AI-Generated" Feel

These came up repeatedly in practitioner discussions and match the product needs here.

Avoid:

- generic three-column layouts everywhere
- oversized hero space without operational value
- identical card treatment on every section regardless of importance
- pages that look balanced but do not have a clear dominant action
- sections added only because "dashboards usually have them"

Prefer:

- one strong focal area per screen
- asymmetry when it improves hierarchy
- domain-specific labels and structure
- visible urgency where urgency is real
- modules that feel designed for the task, not copied from a template

## Comune-Specific Defaults

- Community routes should feel like workspaces, not landing pages.
- Residents, households, units, and invoicing should share a recognizable shell but not identical page compositions.
- Maintenance and invoicing pages should surface urgency early.
- Setup flows should feel guided and short, not like long administrative forms.
- Jamaica context should influence date, currency, and address presentation where relevant.

## Sources

- Laws of UX: https://lawsofux.com/laws/
- Laws of UX, Aesthetic-Usability Effect: https://lawsofux.com/aesthetic-usability-effect/
- NN/g visual-design principles poster: https://media.nngroup.com/media/articles/attachments/Visual-Design-Principles-Poster.pdf
- NN/g usability heuristics summary: https://media.nngroup.com/media/articles/attachments/Heuristic_Summary1-compressed.pdf
- Reddit discussion on generic AI-looking websites: https://www.reddit.com/r/SaasDevelopers/comments/1r5zwro/aigenerated_websites_always_look_generic_how_do/
- Reddit discussion on what makes sites feel AI-generated: https://www.reddit.com/r/website/comments/1rfkiqm/advice_whats_meant_by_looks_ai/
