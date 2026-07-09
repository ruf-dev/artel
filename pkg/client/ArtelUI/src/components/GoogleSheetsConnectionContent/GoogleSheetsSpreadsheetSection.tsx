import {Button} from "@vervstack/chures"

import cls from "@/components/GoogleSheetsConnectionContent/GoogleSheetsSpreadsheetSection.module.css"
import {Spreadsheet} from "@/app/api/artel/external_connections.pb.ts"
import SpreadsheetRow from "@/components/GoogleSheetsConnectionContent/SpreadsheetRow.tsx"

export default function GoogleSheetsSpreadsheetSection({spreadsheets, spreadsheetsLoading, onAdd, onRemove}: {
    spreadsheets: Spreadsheet[]
    spreadsheetsLoading: boolean
    onAdd: () => void
    onRemove: (sheet: Spreadsheet) => void
}) {
    return (
        <div className={cls.GoogleSheetsSpreadsheetSectionContainer}>
            <div className={cls.SpreadsheetSectionHeader}>
                <span className={cls.SpreadsheetSectionLabel}>Spreadsheets</span>
                <Button variant="ghost" onClick={onAdd}>+ Add</Button>
            </div>
            {spreadsheetsLoading ? (
                <p className={cls.EmptySpreadsheets}>Loading…</p>
            ) : spreadsheets.length === 0 ? (
                <p className={cls.EmptySpreadsheets}>No spreadsheets added yet.</p>
            ) : (
                <div className={cls.SpreadsheetList}>
                    {spreadsheets.map(sheet => (
                        <SpreadsheetRow key={sheet.id} sheet={sheet} onRemove={() => onRemove(sheet)}/>
                    ))}
                </div>
            )}
        </div>
    )
}
