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
    needsAttention?: boolean
}

// Lets the user change the OpenRouter model mid-conversation — wired to
// useSimpleChatSession's currentModel/setModel, which stamps whatever is
// currently selected onto every outgoing user_message. `disabled` renders a dimmed,
// non-interactive trigger (used by the Docker topbar, where the model is fixed to
// Claude Code and there's nothing to pick). `needsAttention` visually flags "no model
// selected yet" so the user notices where to fix it.
export default function ModelSwitcher(props: Props) {
    const {likedIds, isLiked, toggleLiked} = useLikedModels()

    function handleChange(ids: string[]) {
        if (props.disabled) return
        const next = ids[0]
        if (next) props.onChange(next)
    }

    const groupedModels = useMemo(
        () => sortLikedToTop(groupModelsByFamily(props.models), likedIds),
        [props.models, likedIds],
    )

    const latestFlags = useMemo(() => buildLatestModelFlags(props.models), [props.models])

    const searchModels = useCallback(
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (query: string) => Promise.resolve(filterModelTree(groupedModels, query) as any),
        [groupedModels],
    )

    return (
        <div
            className={cn(
                cls.ModelSwitcherContainer,
                props.disabled && cls.ModelSwitcherDisabled,
                props.needsAttention && cls.ModelSwitcherNeedsAttention,
            )}
            aria-disabled={props.disabled || undefined}
        >
            <Dropdown
                placeholder={props.placeholder ?? "Model…"}
                isLoading={props.isLoading}
                options={groupedModels}
                value={props.value ? [props.value] : []}
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
