# Reminder: verify mobile top nav sidebar

Status: needs manual check (not yet self-tested per frontend rules — no
dev server/browser was spun up for this change).

## What changed

Top nav tabs no longer overflow on mobile. On narrow/portrait screens the
`TopbarNav` tab row is replaced by a hamburger button that opens a sliding
left-side drawer with the same nav links.

Files: `pkg/client/ArtelUI/src/segments/Topbar/Topbar.tsx` and new
components under `pkg/client/ArtelUI/src/segments/Topbar/components/`
(`TopbarMobileTrigger`, `TopbarMobileDrawer`, `TopbarDrawerCloseButton`,
`TopbarHamburgerIcon`, `TopbarCloseIcon`), plus
`pkg/client/ArtelUI/src/app/hooks/useIsMobileNav.ts`.

## What to check

- [ ] Shrink browser to phone width (portrait) — tabs disappear, hamburger
      icon appears next to the brand.
- [ ] Click hamburger — drawer slides in from the left with all nav links
      (Vaults/Notes/Connections/Toolbox/Tracts), backdrop dims the rest of
      the page.
- [ ] Clicking a link navigates and closes the drawer; clicking the
      backdrop or the close (chevron) button also closes it without
      navigating.
- [ ] Rotate to landscape: above ~767px wide it should switch back to the
      normal desktop tab row; below that width (narrow landscape) it
      should stay in mobile/hamburger mode.
- [ ] Desktop width unaffected — tabs render as before.
