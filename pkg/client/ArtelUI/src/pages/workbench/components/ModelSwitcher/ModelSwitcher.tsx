import {useCallback, useMemo} from "react"
import {Dropdown} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import {useLikedModels} from "@/app/hooks/useLikedModels.ts"
import ModelOptionRow from "@/components/ModelOptionRow/ModelOptionRow.tsx"
import cls from "@/pages/workbench/components/ModelSwitcher/ModelSwitcher.module.css"
import {
    buildLatestModelFlags,
    filterModelTree,
    getDropdownOptionId,
    groupModelsByFamily,
    sortLikedToTop,
} from "@/processes/groupModelsByFamily.ts"

interface Props {
    models: string[]
    value: string
    isLoading?: boolean
    onChange: (model: string) => void
    disabled?: boolean
    placeholder?: string
}

// Lets the user change the OpenRouter model mid-conversation — wired to
// useSimpleChatSession's currentModel/setModel, which stamps whatever is
// currently selected onto every outgoing user_message. `disabled` renders a dimmed,
// non-interactive trigger (used by the Docker topbar, where the model is fixed to
// Claude Code and there's nothing to pick).
export default function ModelSwitcher({models, value, isLoading, onChange, disabled, placeholder}: Props) {
    const {likedIds, isLiked, toggleLiked} = useLikedModels()

    function handleChange(ids: string[]) {
        if (disabled) return
        const next = ids[0]
        if (next) onChange(next)
    }

    const groupedModels = useMemo(
        () => sortLikedToTop(groupModelsByFamily(models), likedIds),
        [models, likedIds],
    )

    const latestFlags = useMemo(() => buildLatestModelFlags(models), [models])

    const searchModels = useCallback(
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (query: string) => Promise.resolve(filterModelTree(groupedModels, query) as any),
        [groupedModels],
    )

    return (
        <div
            className={cn(cls.ModelSwitcherContainer, disabled && cls.ModelSwitcherDisabled)}
            aria-disabled={disabled || undefined}
        >
            <Dropdown
                placeholder={placeholder ?? "Model…"}
                isLoading={isLoading}
                options={groupedModels}
                value={value ? [value] : []}
                onChange={handleChange}
                onSearch={searchModels}
                portal
                renderOption={(opt, state) => (
                    <ModelOptionRow
                        opt={opt}
                        state={state}
                        isLatest={latestFlags.get(getDropdownOptionId(opt)) ?? false}
                        liked={isLiked(getDropdownOptionId(opt))}
                        onToggleLike={toggleLiked}
                    />
                )}
            />
        </div>
    )
}
