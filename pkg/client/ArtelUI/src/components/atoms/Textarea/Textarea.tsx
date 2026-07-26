// TODO: chures has no multiline variant yet, drop this wrapper once it does
import {cn} from "@/app/utils/cn.ts"
import cls from "@/components/atoms/Textarea/Textarea.module.css"

interface Props {
    value: string
    setValue: (v: string) => void
    disabled?: boolean
    placeholder?: string
    className?: string
    rows?: number
}

export default function Textarea(props: Props) {
    return (
        <div className={cn(cls.TextareaContainer, props.className)}>
            <textarea
                value={props.value}
                onChange={e => props.setValue(e.target.value)}
                disabled={props.disabled}
                placeholder={props.placeholder}
                rows={props.rows ?? 4}
            />
        </div>
    )
}
