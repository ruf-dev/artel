import ComposerCtxChip from "@/pages/workbench/components/Chat/components/ComposerCtxChip/ComposerCtxChip.tsx"
import cls from "@/pages/workbench/components/Chat/components/ComposerCtxRow/ComposerCtxRow.module.css"

interface Props {
    attachedPaths: string[]
    onRemoveAttachment: (path: string) => void
}

// Wraps the composer's attachment chips in a collapsible track so the composer
// card smoothly grows/shrinks as files are attached/removed instead of jumping —
// see ComposerCtxRow.module.css for the grid-template-rows technique.
export default function ComposerCtxRow({attachedPaths, onRemoveAttachment}: Props) {
    return (
        <div className={cls.ComposerCtxRowContainer}>
            <div className={cls.Inner}>
                {attachedPaths.map(p => (
                    <ComposerCtxChip key={p} path={p} onRemove={onRemoveAttachment}/>
                ))}
            </div>
        </div>
    )
}
