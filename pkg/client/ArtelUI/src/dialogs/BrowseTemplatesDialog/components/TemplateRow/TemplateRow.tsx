import {Button} from "@vervstack/chures"

import cls from "@/dialogs/BrowseTemplatesDialog/BrowseTemplatesDialog.module.css"
import {TractTemplateSummary} from "@/processes/Tracts.ts"
import TractsIcon from "@/segments/Topbar/components/icons/TractsIcon.tsx"

interface Props {
    template: TractTemplateSummary
    onClick: () => void
}

export default function TemplateRow({template, onClick}: Props) {
    return (
        <Button variant="ghost" className={cls.TemplateRowContainer} onClick={onClick}>
            <TractsIcon className={cls.RowIcon}/>
            <div className={cls.RowInfo}>
                <span className={cls.RowName}>{template.name}</span>
                {template.description && <span className={cls.RowDesc}>{template.description}</span>}
            </div>
        </Button>
    )
}
