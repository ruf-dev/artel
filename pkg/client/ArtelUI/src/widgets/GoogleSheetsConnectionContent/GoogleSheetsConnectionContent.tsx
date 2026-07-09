import {useEffect, useRef} from "react"
import {Button, ConfirmDialog} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/widgets/GoogleSheetsConnectionContent/GoogleSheetsConnectionContent.module.css"
import {ExternalConnectionInfo, Spreadsheet} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {parseScopeList, SCOPE_INFO, trimScope} from "@/app/utils/googleScopes.ts"


type GapiWindow = Window & { gapi: { load: (lib: string, cb: () => void) => void } }

let gapiPromise: Promise<void> | null = null

function loadGapi(): Promise<void> {
    if (gapiPromise) return gapiPromise
    gapiPromise = new Promise<void>((resolve, reject) => {
        const script = document.createElement("script")
        script.src = "https://apis.google.com/js/api.js"
        script.onload = () => (window as unknown as GapiWindow).gapi.load("picker", resolve)
        script.onerror = reject
        document.head.appendChild(script)
    })
    return gapiPromise
}

function SpreadsheetRow({sheet, onRemove}: { sheet: Spreadsheet; onRemove: () => void }) {
    return (
        <div className={cls.SpreadsheetRow}>
            <span className={cls.SpreadsheetName}>{sheet.name}</span>
            <Button variant="iconDanger" onClick={onRemove} aria-label="Remove">×</Button>
        </div>
    )
}

export default function GoogleSheetsConnectionContent({connection, onDisconnect}: {
    connection: ExternalConnectionInfo
    onDisconnect: () => void
}) {
    const {OpenDialog, CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const {spreadsheets, spreadsheetsLoading, fetchSpreadsheets, addSpreadsheet, removeSpreadsheet, getPickerToken} =
        useExternalConnections()
    const pickerOpenRef = useRef(false)

    const email = connection.google?.email
    const scopes = connection.google?.scopes
    const connectedAt = connection.createdAt
        ? new Date(connection.createdAt).toLocaleDateString(undefined, {
            year: "numeric",
            month: "short",
            day: "numeric",
        })
        : "—"

    const scopeList = scopes ? parseScopeList(scopes) : []
    const scopeTooltipHtml = scopeList
        .map(s => {
            const name = trimScope(s)
            const desc = SCOPE_INFO[name]
            return desc ? `<b>${name}</b>: ${desc}` : `<b>${name}</b>`
        })
        .join("<br/>")

    useEffect(() => {
        void fetchSpreadsheets()
    }, [])

    function openPicker() {
        if (pickerOpenRef.current) return
        pickerOpenRef.current = true

        getPickerToken()
            .then(token => loadGapi().then(() => token))
            .then(token => {
                const picker = new google.picker.PickerBuilder()
                    .addView(new google.picker.DocsView(google.picker.ViewId.SPREADSHEETS))
                    .setOAuthToken(token)
                    .setDeveloperKey(import.meta.env.VITE_GOOGLE_API_KEY ?? "")
                    .setCallback((data: google.picker.ResponseObject) => {
                        const action = data[google.picker.Response.ACTION]
                        if (action === google.picker.Action.PICKED) {
                            const docs = data[google.picker.Response.DOCUMENTS]
                            const file = docs?.[0]
                            if (file) {
                                addSpreadsheet(
                                    file[google.picker.Document.ID],
                                    file[google.picker.Document.NAME] ?? file[google.picker.Document.ID],
                                ).catch(e => bakeError("Failed to add spreadsheet", e))
                            }
                        }
                        if (action === google.picker.Action.PICKED || action === google.picker.Action.CANCEL) {
                            pickerOpenRef.current = false
                        }
                    })
                    .build()
                picker.setVisible(true)
            })
            .catch(e => {
                pickerOpenRef.current = false
                bakeError("Failed to open picker", e)
            })
    }

    function handleRemove(sheet: Spreadsheet) {
        OpenDialog(
            <ConfirmDialog
                title="Remove spreadsheet"
                message={`Remove "${sheet.name}" from MCP access? The file is not deleted from Google Drive.`}
                confirmLabel="Remove"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() => removeSpreadsheet(sheet.id ?? "").catch(e => bakeError("Failed to remove spreadsheet", e))}
            />
        )
    }

    return (
        <>
            <div className={cls.InfoRows}>
                {email && (
                    <div className={cls.InfoRow}>
                        <span className={cls.InfoLabel}>Account</span>
                        <span className={cls.InfoValue}>{email}</span>
                    </div>
                )}
                {scopeList.length > 0 && (
                    <div className={cls.InfoRow}>
                        <span className={cn(cls.InfoLabel, cls.ScopesLabel)}>
                            Scopes
                            <Button
                                variant="ghost"
                                className={cls.ScopeHelp}
                                data-tooltip-id="root-tooltip"
                                data-tooltip-html={scopeTooltipHtml}
                                aria-label="What are scopes?"
                            >?</Button>
                        </span>
                        <span className={cls.InfoValue}>{scopeList.map(trimScope).join(", ")}</span>
                    </div>
                )}
                <div className={cls.InfoRow}>
                    <span className={cls.InfoLabel}>Connected</span>
                    <span className={cls.InfoValue}>{connectedAt}</span>
                </div>
            </div>

            <div className={cls.SpreadsheetSection}>
                <div className={cls.SpreadsheetSectionHeader}>
                    <span className={cls.SpreadsheetSectionLabel}>Spreadsheets</span>
                    <Button variant="ghost" onClick={openPicker}>+ Add</Button>
                </div>
                {spreadsheetsLoading ? (
                    <p className={cls.EmptySpreadsheets}>Loading…</p>
                ) : spreadsheets.length === 0 ? (
                    <p className={cls.EmptySpreadsheets}>No spreadsheets added yet.</p>
                ) : (
                    <div className={cls.SpreadsheetList}>
                        {spreadsheets.map(sheet => (
                            <SpreadsheetRow key={sheet.id} sheet={sheet} onRemove={() => handleRemove(sheet)}/>
                        ))}
                    </div>
                )}
            </div>

            <div className={cls.Actions}>
                <Button variant="danger" onClick={onDisconnect}>Disconnect</Button>
            </div>
        </>
    )
}
