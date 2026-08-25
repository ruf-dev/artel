import type {DropdownOption, RenderOptionState} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import LikeToggleButton from "@/components/atoms/LikeToggleButton/LikeToggleButton.tsx"
import ModelFamilyIcon from "@/components/ModelFamilyIcon/ModelFamilyIcon.tsx"
import cls from "@/components/ModelOptionRow/ModelOptionRow.module.css"
import {getDropdownOptionId, getModelFamily} from "@/processes/groupModelsByFamily.ts"

interface Props {
    opt: DropdownOption
    state: RenderOptionState
    isLatest: boolean
    liked: boolean
    onToggleLike: (id: string) => void
}

// chures' DropdownOption is `string | {id, name}`; getDropdownOptionId already
// handles the id half of that union (shared with ModelSwitcher/useOpenRouterModels'
// renderOption callbacks), this covers the name half — mirrors the spirit of
// chures' own (non-exported) getOptionLabel.
function getOptionName(opt: DropdownOption): string {
    return typeof opt === "string" ? opt : opt.name
}

// Full-row override for chures' <Dropdown renderOption>, used by ModelSwitcher and
// SimpleChatModeForm to show a provider icon, a "latest" badge, and a like toggle
// alongside the model name. Since this is a complete row override, nothing from
// chures renders on top of it — the selected-checkmark and mousedown-select
// behavior chures normally provides are replicated here (see DropdownOptionRow.tsx
// in chures for the precedent this mirrors).
export default function ModelOptionRow({opt, state, isLatest, liked, onToggleLike}: Props) {
    const id = getDropdownOptionId(opt)
    const name = getOptionName(opt)
    const family = getModelFamily(id)

    function handleMouseDown(e: React.MouseEvent) {
        e.preventDefault()
        state.onPick()
    }

    return (
        <div
            className={cn(
                cls.ModelOptionRowContainer,
                state.isSelected && cls.Selected,
                state.indented && cls.Indented,
            )}
            onMouseDown={handleMouseDown}
        >
            <ModelFamilyIcon family={family}/>
            <span className={cls.Label}>{name}</span>
            {isLatest && <span className={cls.LatestBadge}>latest</span>}
            {state.multiSelect && state.isSelected && <span className={cls.Checkmark}>✓</span>}
            <LikeToggleButton liked={liked} onToggle={() => onToggleLike(id)}/>
        </div>
    )
}
