import cls from "@/pages/segments/Dialog.module.css"

import {useDialog} from "@/app/hooks/Dialog"

export default function Dialog() {
    const {children, CloseDialog} = useDialog()

    if (!children) return null

    return (
        <div
            className={cls.DialogContainer}
            onMouseDown={(e) => {
                if (e.target === e.currentTarget) CloseDialog()
            }}
        >
            <div
                onMouseDown={(e) => {
                    e.stopPropagation()
                }}
            >
                {
                    children.map((v, idx) => {
                        return (<div key={idx}>
                            {v}
                        </div>)
                    })
                }
            </div>
        </div>
    )
}
