import {Button} from "@vervstack/chures"

import cls from "@/components/LlmKeyConnectedContent/LlmKeyConnectedContent.module.css"

interface LlmKeyConnectedContentProps {
    fields: Record<string, string>
    onDisconnect: () => void
}

export default function LlmKeyConnectedContent({fields, onDisconnect}: LlmKeyConnectedContentProps) {
    const availableModels = fields.available_models ? fields.available_models.split(",").filter(Boolean) : []

    return (
        <div className={cls.LlmKeyConnectedContentContainer}>
            <p className={cls.ModalSub}>
                Connected with key <b>{fields.key_preview}</b>.
            </p>
            {fields.default_model && (
                <label className={cls.Field}>
                    <span className={cls.FieldLabel}>Default model</span>
                    <div className={cls.FieldValue}>{fields.default_model}</div>
                </label>
            )}
            {availableModels.length > 0 && (
                <label className={cls.Field}>
                    <span className={cls.FieldLabel}>Available models</span>
                    <div className={cls.FieldValue}>{availableModels.join(", ")}</div>
                </label>
            )}
            <div className={cls.ModalActions}>
                <Button variant="danger" onClick={onDisconnect}>Disconnect</Button>
            </div>
        </div>
    )
}
