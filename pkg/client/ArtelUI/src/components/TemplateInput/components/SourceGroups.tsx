import cls from "@/components/TemplateInput/TemplateInput.module.css"
import {cn} from "@/app/utils/cn.ts"
import {FilteredGroup, FlatListItem} from "@/components/TemplateInput/processes/templateRefs.ts"

interface Props {
    groups: FilteredGroup[]
    flatList: FlatListItem[]
    activeIndex: number
    onSelect: (ref: string) => void
}

export default function SourceGroups({groups, flatList, activeIndex, onSelect}: Props) {
    if (groups.length === 0) {
        return <div className={cls.Empty}>No matching references.</div>
    }

    function indexOf(ref: string): number {
        return flatList.findIndex(f => f.ref === ref)
    }

    return (
        <>
            {groups.map(g => (
                <div className={cls.Group} key={g.source.id}>
                    {g.rootMatches && (
                        <div
                            className={cn(cls.GroupLabel, indexOf(g.source.id) === activeIndex && cls.EntryActive)}
                            onMouseDown={e => {
                                e.preventDefault()
                                onSelect(g.source.id)
                            }}
                        >
                            {g.source.label}
                        </div>
                    )}
                    {g.entries.map(e => (
                        <div
                            key={e.key}
                            className={cn(cls.Entry, indexOf(e.ref) === activeIndex && cls.EntryActive)}
                            style={{paddingLeft: `${0.5 + e.depth * 0.75}rem`}}
                            onMouseDown={ev => {
                                ev.preventDefault()
                                onSelect(e.ref)
                            }}
                            data-tooltip-id={e.description ? "root-tooltip" : undefined}
                            data-tooltip-content={e.description}
                        >
                            <span className={cls.EntryLabel}>{e.label}</span>
                            {e.type && <span className={cls.EntryType}>{e.type}</span>}
                        </div>
                    ))}
                </div>
            ))}
        </>
    )
}
