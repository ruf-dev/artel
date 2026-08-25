// TODO: chures has no like/favorite icon-button yet, drop this wrapper once it does
import {Button} from "@vervstack/chures"

import {cn} from "@/app/utils/cn.ts"
import cls from "@/components/atoms/LikeToggleButton/LikeToggleButton.module.css"

interface Props {
    liked: boolean
    onToggle: () => void
}

// Small heart toggle used inside ModelOptionRow to mark an OpenRouter model as
// liked/unliked. Lives inside a dropdown option row whose own container picks the
// model on mousedown — stopping mousedown propagation here keeps a click on the
// heart from also selecting the row's model.
export default function LikeToggleButton({liked, onToggle}: Props) {
    return (
        <div className={cls.LikeToggleButtonContainer}>
            <Button
                variant="unstyled"
                className={cn(cls.ToggleBtn, liked && cls.Liked)}
                aria-pressed={liked}
                aria-label={liked ? "Unlike model" : "Like model"}
                onMouseDown={e => e.stopPropagation()}
                onClick={onToggle}
            >
                {liked ? (
                    <svg
                        width="14"
                        height="14"
                        viewBox="0 0 14 14"
                        fill="currentColor"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            d={
                                "M7 12.35 1.98 7.6C.5 6.22.44 3.95 1.85 2.5c1.35-1.4 3.53-1.44 4.94-.1L7 2.6l.21-.2" +
                                "c1.4-1.34 3.59-1.3 4.94.1 1.4 1.45 1.35 3.72-.13 5.1L7 12.35Z"
                            }
                        />
                    </svg>
                ) : (
                    <svg
                        width="14"
                        height="14"
                        viewBox="0 0 14 14"
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="1.15"
                        strokeLinejoin="round"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            d={
                                "M7 12.35 1.98 7.6C.5 6.22.44 3.95 1.85 2.5c1.35-1.4 3.53-1.44 4.94-.1L7 2.6l.21-.2" +
                                "c1.4-1.34 3.59-1.3 4.94.1 1.4 1.45 1.35 3.72-.13 5.1L7 12.35Z"
                            }
                        />
                    </svg>
                )}
            </Button>
        </div>
    )
}
