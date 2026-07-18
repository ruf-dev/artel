// CodeMirror setup for the script step's editable code body: a theme mapping the editor's
// chrome onto this project's existing CSS custom properties (so it tracks light/dark and the
// coral accent automatically, no JS theme branching needed), a syntax highlight style using
// the --syntax-* tokens, and a completion source that layers the step's declared param names
// on top of CodeMirror's built-in JS keyword/snippet completions.

import {javascriptLanguage} from "@codemirror/lang-javascript"
import {HighlightStyle} from "@codemirror/language"
import {EditorView} from "@codemirror/view"
import {tags as t} from "@lezer/highlight"
import type {Completion, CompletionContext, CompletionResult} from "@codemirror/autocomplete"

import {ScriptParam} from "@/processes/Tracts.ts"

export const scriptEditorTheme = EditorView.theme({
    "&": {
        color: "var(--color-fg-primary)",
        backgroundColor: "var(--color-bg-input)",
        fontFamily: "var(--font-mono)",
        fontSize: "0.6875rem",
        borderRadius: "0.4375rem",
        border: "var(--border-thickness) solid var(--color-border-light)",
        overflow: "hidden",
        transition: "border-color 140ms ease",
    },
    "&.cm-focused": {
        outline: "none",
        borderColor: "var(--coral)",
    },
    ".cm-content": {
        caretColor: "var(--coral)",
        padding: "0.5rem 0.625rem",
    },
    ".cm-line": {
        padding: 0,
    },
    "&.cm-focused .cm-selectionBackground, .cm-selectionBackground": {
        backgroundColor: "var(--color-accent-dim-hover)",
    },
    ".cm-activeLine": {
        backgroundColor: "var(--color-bg-hover)",
    },
    ".cm-tooltip-autocomplete": {
        fontFamily: "var(--font-mono)",
        fontSize: "0.6875rem",
        backgroundColor: "var(--color-bg-elevated)",
        border: "var(--border-thickness) solid var(--color-border-light)",
    },
    ".cm-tooltip-autocomplete ul li[aria-selected]": {
        backgroundColor: "var(--color-accent-dim)",
        color: "var(--color-fg-primary)",
    },
})

export const scriptHighlightStyle = HighlightStyle.define([
    {tag: t.keyword, color: "var(--syntax-keyword)"},
    {tag: [t.string, t.special(t.string)], color: "var(--syntax-string)"},
    {tag: [t.number, t.bool, t.null], color: "var(--syntax-number)"},
    {tag: t.comment, color: "var(--syntax-comment)", fontStyle: "italic"},
    {tag: [t.propertyName, t.function(t.variableName)], color: "var(--syntax-property)"},
])

/** paramCompletionSource suggests the step's declared input/output param names (the
 * variables the generated signature/prelude actually put in scope) as CodeMirror completion
 * options. Registered via `javascriptLanguage.data.of` rather than `autocompletion({override})`
 * so it adds to, instead of replacing, the language's own keyword/snippet completions. */
export function paramCompletionSource(inputParams: ScriptParam[], outputParams: ScriptParam[]) {
    const options: Completion[] = [
        ...inputParams.map(p => ({label: p.name, type: "variable", detail: p.type.type})),
        ...outputParams.map(p => ({label: p.name, type: "variable", detail: p.type.type})),
    ]

    function source(context: CompletionContext): CompletionResult | null {
        const word = context.matchBefore(/\w*/)
        if (!word || (word.from === word.to && !context.explicit)) return null
        return {from: word.from, options}
    }

    return javascriptLanguage.data.of({autocomplete: source})
}
