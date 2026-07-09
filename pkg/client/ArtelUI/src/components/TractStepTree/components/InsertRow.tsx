import {Button} from "@vervstack/chures"

import cls from "@/components/TractStepTree/TractStepTree.module.css"

interface Props {
    onClick: () => void
}

export default function InsertRow({onClick}: Props) {
    return (
        <div className={cls.InsertRow}>
            <span className={cls.InsertLine}/>
            <Button variant="ghost" className={cls.InsertBtn} onClick={onClick} aria-label="Add step">+</Button>
            <span className={cls.InsertLine}/>
        </div>
    )
}
