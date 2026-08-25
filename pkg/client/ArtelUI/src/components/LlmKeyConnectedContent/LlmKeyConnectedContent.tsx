import {useCallback, useMemo} from "react"
import {Button, Dropdown} from "@vervstack/chures"

import cls from "@/components/LlmKeyConnectedContent/LlmKeyConnectedContent.module.css"
import {useLikedModels} from "@/app/hooks/useLikedModels.ts"
import ModelOptionRow from "@/components/ModelOptionRow/ModelOptionRow.tsx"
import {
    buildLatestModelFlags,
    filterModelTree,
    getDropdownOptionId,
    groupModelsByFamily,
    sortLikedToTop,
} from "@/processes/groupModelsByFamily.ts"

interface LlmKeyConnectedContentProps {
    fields: Record<string, string>
    onDisconnect: () => void
    onViewUsage?: () => void
    // Present only for OpenRouter's connected card — turns the "Available models"
    // dropdown below from a browse-only list into a real single-select that
    // persists the pick as this connection's default model. Every other
    // Manage*Dialog caller (Anthropic, OpenAI, S3, Postgres, CouchDB) omits this,
    // keeping their cards read-only exactly as before.
    onSelectDefaultModel?: (model: string) => Promise<void>
}

export default function LlmKeyConnectedContent(
    {fields, onDisconnect, onViewUsage, onSelectDefaultModel}: LlmKeyConnectedContentProps,
) {
    const availableModels = fields.available_models ? fields.available_models.split(",").filter(Boolean) : []
    const {likedIds, isLiked, toggleLiked} = useLikedModels()

    const groupedModels = useMemo(
        () => sortLikedToTop(groupModelsByFamily(availableModels), likedIds),
        [availableModels, likedIds],
    )
    const latestFlags = useMemo(() => buildLatestModelFlags(availableModels), [availableModels])

    const searchModels = useCallback(
        (query: string) => Promise.resolve(filterModelTree(groupedModels, query)),
        [groupedModels],
    )

    function handleChange(ids: string[]) {
        const next = ids[0]
        if (next && onSelectDefaultModel) {
            onSelectDefaultModel(next).catch(() => {})
        }
    }

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
                    <span className={cls.FieldLabel}>{`Available models (${availableModels.length})`}</span>
                    <div className={cls.DropdownWrapper}>
                        <Dropdown
                            options={groupedModels}
                            value={onSelectDefaultModel && fields.default_model ? [fields.default_model] : []}
                            onChange={onSelectDefaultModel ? handleChange : () => {}}
                            onSearch={searchModels}
                            renderOption={(opt, state) => (
                                <ModelOptionRow
                                    opt={opt}
                                    state={state}
                                    isLatest={latestFlags.get(getDropdownOptionId(opt)) ?? false}
                                    liked={isLiked(getDropdownOptionId(opt))}
                                    onToggleLike={toggleLiked}
                                />
                            )}
                            placeholder={onSelectDefaultModel ? "Choose default model…" : "Browse models…"}
                            portal
                        />
                    </div>
                </label>
            )}
            <div className={cls.ModalActions}>
                {onViewUsage && <Button variant="secondary" onClick={onViewUsage}>View usage</Button>}
                <Button variant="danger" onClick={onDisconnect}>Disconnect</Button>
            </div>
        </div>
    )
}
