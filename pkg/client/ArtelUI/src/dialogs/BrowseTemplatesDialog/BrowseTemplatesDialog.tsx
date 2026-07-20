import {useEffect, useState} from "react"

import {useTractTemplates} from "@/app/hooks/TractTemplates.ts"
import {TractTemplateSummary} from "@/processes/Tracts.ts"
import ListScreen from "@/dialogs/BrowseTemplatesDialog/screens/ListScreen.tsx"
import DetailScreen from "@/dialogs/BrowseTemplatesDialog/screens/DetailScreen.tsx"

export default function BrowseTemplatesDialog() {
    const {templates, loading, fetch} = useTractTemplates()
    const [selected, setSelected] = useState<TractTemplateSummary | null>(null)

    useEffect(() => {
        void fetch()
    }, [fetch])

    if (selected) {
        return <DetailScreen template={selected} onBack={() => setSelected(null)}/>
    }

    return <ListScreen templates={templates} loading={loading} onSelect={setSelected}/>
}
