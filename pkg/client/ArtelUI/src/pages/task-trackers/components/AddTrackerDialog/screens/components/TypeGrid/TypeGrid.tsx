import {Button} from "@vervstack/chures"

import TypeCardText from "@/pages/task-trackers/components/AddTrackerDialog/screens/components/TypeGrid/components/TypeCardText/TypeCardText.tsx"
import cls from "@/pages/task-trackers/components/AddTrackerDialog/screens/components/TypeGrid/TypeGrid.module.css"

export default function TypeGrid({onChooseTrello}: {onChooseTrello: () => void}) {
    return (
        <div className={cls.TypeGridContainer}>
            <Button variant="ghost" className={cls.TypeCard} onClick={onChooseTrello}>
                <div className={cls.TypeCardIcon}>T</div>
                <TypeCardText label="Trello" desc="Connect via API key and token"/>
            </Button>
        </div>
    )
}
