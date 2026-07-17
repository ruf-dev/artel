import stylelint from "stylelint"

const ruleName = "artel/dialog-scrollable"

const messages = stylelint.utils.ruleMessages(ruleName, {
    expected: (detail) => `Dialog container has no way to scroll long content — ${detail}. Either (a) set both `
        + "max-height (var(--dialog-max-height) or var(--dialog-max-height-lg)) and overflow-y: auto directly on "
        + "the dialog's top-level *Container rule (see WebhookDetailsDialog.module.css), or (b) give the "
        + "*Container rule max-height plus display: flex; flex-direction: column, and give a nested content "
        + "region (between the pinned header and footer) its own overflow-y: auto (see "
        + "UserSubscriptionDialog.module.css's .ModalBody).",
})

const meta = {
    url: "https://github.com/Red-Sock/artel/blob/master/pkg/client/ArtelUI/stylelint-rules/dialog-scrollable.js",
}

function isOverflowY(decl) {
    return decl.prop === "overflow-y" || decl.prop === "overflow"
}

function directDeclsOf(rule) {
    return rule.nodes.filter((node) => node.type === "decl")
}

/** Every rule in the file, top-level or nested at any depth — some older dialog CSS
 * files predate this codebase's native-CSS-nesting convention and declare each class
 * as a flat top-level sibling rule instead of nesting it under the container. */
function allRules(root) {
    const rules = []
    function walk(node) {
        for (const child of node.nodes ?? []) {
            if (child.type === "rule") {
                rules.push(child)
                walk(child)
            }
        }
    }
    walk(root)
    return rules
}

/**
 * A handful of dialogs (FastSetupDialog, CreateKeyDialog, CreateVaultDialog) render
 * their own full-viewport backdrop instead of relying on the shared overlay in
 * pages/segments/Dialog.tsx — `position: fixed` + `inset: 0` (or all four offsets) —
 * with the actual card as a nested rule carrying its own `max-width`. Capping height
 * on that backdrop is wrong (an over-constrained fixed box with top+bottom both set
 * either gets clamped and ignores one offset, or is simply pointless since the
 * backdrop is deliberately viewport-sized); the card is the thing that needs to cap
 * its own height, so find it instead.
 */
function findScrollTarget(containerRule, root) {
    const decls = directDeclsOf(containerRule)
    const isFixed = decls.some((decl) => decl.prop === "position" && decl.value.trim() === "fixed")
    const isFullViewport = decls.some((decl) => decl.prop === "inset")
        || ["top", "right", "bottom", "left"].every(
            (side) => decls.some((decl) => decl.prop === side && decl.value.trim() === "0")
        )

    if (!isFixed || !isFullViewport) return containerRule

    // The card may be nested under the backdrop (native CSS nesting) or, in older
    // files that predate that convention, declared as a flat sibling rule instead —
    // search the whole file for whichever rule actually carries the card's max-width.
    const card = allRules(root).find(
        (node) => node !== containerRule && directDeclsOf(node).some((decl) => decl.prop === "max-width")
    )
    return card ?? containerRule
}

/**
 * Every dialog's outer shell is the first top-level class rule declared in its
 * module.css file (see the "Component Structure" rule in ArtelUI/CLAUDE.md: the
 * top-level container must be named `{ComponentName}Container`). A dialog can grow
 * arbitrarily tall (more form fields, a longer list, a wider viewport-relative admin
 * screen shrunk into a phone), so it must always cap its own height and scroll
 * internally rather than relying on content happening to fit. Three shapes are all
 * fine, and this rule accepts any of them:
 *   (A) the container itself has both max-height and overflow-y (whole shell scrolls)
 *   (B) the container has max-height + flex column, and *some* rule in the file has
 *       overflow-y (a pinned header/footer with a flexed, scrollable body between them
 *       — the body itself doesn't need its own max-height, it grows into the budget
 *       the container's max-height leaves it via flex)
 *   (C) some rule anywhere in the file self-caps with its own max-height *and*
 *       overflow-y together (a bounded inner list/section keeps the dialog's total
 *       height sane even though the outer container declares no cap of its own)
 */
const dialogScrollable = (enabled) => {
    return (root, result) => {
        if (!enabled) return

        const firstRule = root.nodes.find((node) => node.type === "rule")
        if (!firstRule) return
        const containerRule = findScrollTarget(firstRule, root)

        const rules = allRules(root)
        const selfCapped = rules.some((rule) => {
            const decls = directDeclsOf(rule)
            return decls.some((decl) => decl.prop === "max-height") && decls.some(isOverflowY)
        })
        if (selfCapped) return // pattern A (when it's the container itself) or pattern C

        const containerDecls = directDeclsOf(containerRule)
        const hasMaxHeight = containerDecls.some((decl) => decl.prop === "max-height")
        const hasFlexColumn = containerDecls.some((decl) => decl.prop === "display" && /flex/.test(decl.value))
            && containerDecls.some((decl) => decl.prop === "flex-direction" && decl.value.trim() === "column")
        const anyOverflowY = rules.some((rule) => directDeclsOf(rule).some(isOverflowY))

        if (hasMaxHeight && hasFlexColumn && anyOverflowY) return // pattern B

        stylelint.utils.report({
            message: messages.expected(
                !hasMaxHeight
                    ? 'missing "max-height" on the container'
                    : !anyOverflowY
                        ? 'missing "overflow-y: auto" anywhere in the file'
                        : 'container has max-height but is missing "display: flex; flex-direction: column" '
                            + "to pair with the scrollable region found elsewhere in the file"
            ),
            node: containerRule,
            result,
            ruleName,
        })
    }
}

dialogScrollable.ruleName = ruleName
dialogScrollable.messages = messages
dialogScrollable.meta = meta

export default stylelint.createPlugin(ruleName, dialogScrollable)
