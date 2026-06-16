import cls from "@/components/SelectOption/SelectOption.module.css"

export default function SelectOption({label, selected, onSelect}: {
    label: string
    selected: boolean
    onSelect: () => void
}) {
    return (
        <button
            className={selected ? `${cls.OptionRow} ${cls.OptionRowSelected}` : cls.OptionRow}
            onClick={onSelect}
            type="button"
        >
            <span className={cls.OptionRadio}>{selected ? "●" : "○"}</span>
            <span>{label}</span>
        </button>
    )
}
