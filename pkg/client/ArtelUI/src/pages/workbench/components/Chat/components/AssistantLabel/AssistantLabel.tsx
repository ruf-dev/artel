import cls from "@/pages/workbench/components/Chat/components/AssistantLabel/AssistantLabel.module.css"

interface Props {
    label: string
}

// The muted "who is speaking" row above an assistant message — a coral dot plus
// the model name (api mode) or the literal "Claude Code" (docker mode).
export default function AssistantLabel({label}: Props) {
    return (
        <div className={cls.AssistantLabelContainer}>
            <span className={cls.Dot}/>
            <span className={cls.Label}>{label}</span>
        </div>
    )
}
