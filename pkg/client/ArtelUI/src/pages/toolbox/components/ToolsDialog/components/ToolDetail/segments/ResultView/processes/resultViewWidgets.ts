import type {ComponentType} from "react"

import TrelloTableWidget
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/components/TrelloTableWidget/TrelloTableWidget.tsx"
import {ResultViewWidgetProps}
    // eslint-disable-next-line max-len -- deep nested import path can't be shortened without changing the import
    from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/segments/ResultView/processes/resultViewWidgetTypes.ts"

interface ResultViewWidgetEntry {
    // Label for the toggle button that switches into this widget (see ViewModeToggle) — e.g.
    // "Table", "Chart". Kept per-entry rather than hardcoded since not every widget is a table.
    label: string
    Component: ComponentType<ResultViewWidgetProps>
}

// Keyed by "<connector>:<tool>" for a widget scoped to one tool, or "<connector>:*" for a
// widget that applies to every tool under that MoM connector (e.g. any trello list/board
// response renders as a generic table). Add new entries here to register a custom widget for
// a connector/tool pair; anything unregistered falls back to raw JsonView in ResultView.
const RESULT_VIEW_WIDGETS: Record<string, ResultViewWidgetEntry> = {
    "trello:*": {label: "Table", Component: TrelloTableWidget},
}

export function getResultViewWidget(connector: string, tool: string): ResultViewWidgetEntry | null {
    return RESULT_VIEW_WIDGETS[`${connector}:${tool}`] ?? RESULT_VIEW_WIDGETS[`${connector}:*`] ?? null
}
