import {useState} from "react"
import {Button} from "@vervstack/chures"
import {MorphIcon} from "morphicons/react"
import {MessageSquare, SquarePen} from "lucide"

import cls from "@/pages/workbench/components/WorkbenchTopbar/components/NewChatButton/NewChatButton.module.css"

interface Props {
    onClick: () => void
}

// Starts a new chat thread. A chat-bubble glyph morphs to a compose/edit glyph
// on hover/focus as a lively hover affordance, then morphs back on hover-out —
// mirrors SidebarToggleButton/SettingsFab's chures Button + MorphIcon shape.
export default function NewChatButton({onClick}: Props) {
    const [hovered, setHovered] = useState(false)

    return (
        <div className={cls.NewChatButtonContainer}>
            <Button
                variant="unstyled"
                className={cls.Btn}
                aria-label="New chat"
                onClick={onClick}
                onMouseEnter={() => setHovered(true)}
                onMouseLeave={() => setHovered(false)}
                onFocus={() => setHovered(true)}
                onBlur={() => setHovered(false)}
            >
                <MorphIcon
                    icon={hovered ? SquarePen : MessageSquare}
                    size={20}
                    strokeWidth={1.6}
                    className={cls.Icon}
                />
            </Button>
        </div>
    )
}
