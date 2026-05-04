# Theme Tokens

The current web palette lives in [palette.css](/home/rowana/projects/gated-community/comune/apps/web/src/lib/theme/palette.css).

Use this file as the source of truth for:

- the current named palette
- semantic theme tokens such as `--color-app-bg`, `--color-app-text`, and `--color-app-accent`
- future theme variants applied through `body[data-theme="..."]`

This app uses Tailwind v4's CSS-based `@theme` tokens rather than a legacy `tailwind.config` color map, so palette work should usually happen here first.
