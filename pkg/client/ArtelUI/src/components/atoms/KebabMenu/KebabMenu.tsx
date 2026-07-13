// TODO: chures has no action/context-menu primitive yet, drop this wrapper once it does
import { useEffect, useState } from "react"
import { createPortal } from "react-dom"
import { Button } from "@vervstack/chures"

import { cn } from "@/app/utils/cn.ts"
import cls from "@/components/atoms/KebabMenu/KebabMenu.module.css"

export interface KebabMenuItem {
    label: string
    onClick: () => void
    danger?: boolean
}

interface Props {
    items: KebabMenuItem[]
    title?: string
}

export default function KebabMenu({ items, title }: Props) {
    const [rect, setRect] = useState<{ top: number; left: number } | null>(null)

    useEffect(() => {
        if (!rect) return

        function close() {
            setRect(null)
        }

        document.addEventListener("mousedown", close)
        return () => document.removeEventListener("mousedown", close)
    }, [rect])

    function handleToggle(e: React.MouseEvent<HTMLButtonElement>) {
        e.stopPropagation()
        if (rect) {
            setRect(null)
            return
        }
        const box = e.currentTarget.getBoundingClientRect()
        setRect({ top: box.bottom + 4, left: box.right })
    }

    return (
        <div className={cls.KebabMenuContainer}>
            <Button
                variant="ghost"
                className={cls.TriggerBtn}
                onClick={handleToggle}
                title={title ?? "More actions"}
            >
                <svg viewBox="0 0 12 12" width={11} height={11} fill="currentColor">
                    <circle cx="6" cy="2" r="1.15" />
                    <circle cx="6" cy="6" r="1.15" />
                    <circle cx="6" cy="10" r="1.15" />
                </svg>
            </Button>
            {rect && createPortal(
                <div
                    className={cls.Menu}
                    style={{ top: rect.top, left: rect.left }}
                    onMouseDown={e => e.stopPropagation()}
                >
                    {items.map(item => (
                        <Button
                            key={item.label}
                            variant="ghost"
                            className={cn(cls.MenuItem, item.danger && cls.MenuItemDanger)}
                            onClick={() => {
                                setRect(null)
                                item.onClick()
                            }}
                        >
                            {item.label}
                        </Button>
                    ))}
                </div>,
                document.body,
            )}
        </div>
    )
}
