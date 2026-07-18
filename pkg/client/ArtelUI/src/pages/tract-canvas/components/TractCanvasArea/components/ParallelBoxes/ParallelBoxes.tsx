import cls from "@/pages/tract-canvas/components/TractCanvasArea/components/ParallelBoxes/ParallelBoxes.module.css"
import {cn} from "@/app/utils/cn.ts"
import {ParallelBox} from "@/pages/tract-canvas/processes/tractCanvasLayout.ts"

interface Props {
    boxes: ParallelBox[]
    selectedNodeId: string | null
    onSelectNode: (id: string) => void
}

export default function ParallelBoxes({boxes, selectedNodeId, onSelectNode}: Props) {
    return (
        <>
            {boxes.map(box => (
                <div
                    key={box.id}
                    className={cn(cls.ParallelBox, box.id === selectedNodeId && cls.ParallelBoxSelected)}
                    style={{left: box.x, top: box.y, width: box.width, height: box.height}}
                    data-tract-node
                    onClick={e => {
                        e.stopPropagation()
                        onSelectNode(box.id)
                    }}
                    role="button"
                    tabIndex={0}
                />
            ))}
        </>
    )
}
