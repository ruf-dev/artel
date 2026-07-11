import {useEffect} from "react"
import {Button, ConfirmDialog} from "@vervstack/chures"

import cls from "@/widgets/GoogleSheetsConnectionContent/GoogleSheetsConnectionContent.module.css"
import {ExternalConnectionInfo, Spreadsheet} from "@/app/api/artel/external_connections.pb.ts"
import {useDialog} from "@/app/hooks/Dialog"
import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import {parseScopeList} from "@/app/utils/googleScopes.ts"
import {useSpreadsheetPicker} from "@/widgets/GoogleSheetsConnectionContent/processes/useSpreadsheetPicker.ts"
import GoogleSheetsInfoRows from "@/components/GoogleSheetsConnectionContent/GoogleSheetsInfoRows.tsx"
import GoogleSheetsSpreadsheetSection
    from "@/components/GoogleSheetsConnectionContent/GoogleSheetsSpreadsheetSection.tsx"

export default function GoogleSheetsConnectionContent({connection, onDisconnect}: {
    connection: ExternalConnectionInfo
    onDisconnect: () => void
}) {
    const {OpenDialog, CloseDialog} = useDialog()
    const bakeError = useBakeError()
    const {spreadsheets, spreadsheetsLoading, fetchSpreadsheets, removeSpreadsheet} = useExternalConnections()
    const {openPicker} = useSpreadsheetPicker()

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

    useEffect(() => {
        void fetchSpreadsheets()
    }, [])

    function handleRemove(sheet: Spreadsheet) {
        OpenDialog(
            <ConfirmDialog
                title="Remove spreadsheet"
                message={`Remove "${sheet.name}" from MCP access? The file is not deleted from Google Drive.`}
                confirmLabel="Remove"
                cancelLabel="Cancel"
                danger
                onClose={CloseDialog}
                onConfirm={() =>
                    removeSpreadsheet(sheet.id ?? "").catch(e => bakeError("Failed to remove spreadsheet", e))
                }
            />
        )
    }

    return (
        <>
            <GoogleSheetsInfoRows email={email} connectedAt={connectedAt} scopeList={scopeList}/>

            <GoogleSheetsSpreadsheetSection
                spreadsheets={spreadsheets}
                spreadsheetsLoading={spreadsheetsLoading}
                onAdd={openPicker}
                onRemove={handleRemove}
            />

            <div className={cls.Actions}>
                <Button variant="danger" onClick={onDisconnect}>Disconnect</Button>
            </div>
        </>
    )
}
