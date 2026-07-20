import {Button} from "@vervstack/chures"

import cls from "@/dialogs/BrowseTemplatesDialog/components/TemplateRow/TemplateRow.module.css"
import {TractTemplateSummary} from "@/processes/Tracts.ts"
import TractsIcon from "@/segments/Topbar/components/icons/TractsIcon.tsx"

interface Props {
    template: TractTemplateSummary
    onClick: () => void
}

export default function TemplateRow({template, onClick}: Props) {
    return (
        <Button variant="ghost" className={cls.TemplateRowContainer} onClick={onClick}>
            <div className={cls.Header}>
                <div className={cls.IconFrame}>
                    <TractsIcon className={cls.Icon}/>
                </div>
                <span className={cls.Title}>{template.name}</span>
            </div>
            {template.description && <p className={cls.Description}>{template.description}</p>}
        </Button>
    )
}
