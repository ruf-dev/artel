import cls from "@/dialogs/AddTriggerDialog/AddTriggerDialog.module.css"

export default function KindSelect(
    {value, onChange}: { value: "webhook" | "manual"; onChange: (v: "webhook" | "manual") => void },
) {
    return (
        <select
            className={cls.PlainSelect}
            value={value}
            onChange={e => onChange(e.target.value as "webhook" | "manual")}
        >
            <option value="webhook">Webhook</option>
            <option value="manual">Manual</option>
        </select>
    )
}
