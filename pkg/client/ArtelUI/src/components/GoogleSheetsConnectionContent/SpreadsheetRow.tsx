import {Button} from "@vervstack/chures"

import cls from "@/components/GoogleSheetsConnectionContent/SpreadsheetRow.module.css"
import {Spreadsheet} from "@/app/api/artel/external_connections.pb.ts"

export default function SpreadsheetRow({sheet, onRemove}: { sheet: Spreadsheet; onRemove: () => void }) {
    return (
        <div className={cls.SpreadsheetRowContainer}>
            <span className={cls.SpreadsheetName}>{sheet.name}</span>
            <Button variant="iconDanger" onClick={onRemove} aria-label="Remove">×</Button>
        </div>
    )
}
