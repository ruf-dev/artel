import cls from "@/pages/tract-canvas/segments/HeroSegment/HeroSegment.module.css"
import {useDialog} from "@/app/hooks/Dialog"
import {useTracts} from "@/app/hooks/Tracts.ts"
import NewTractButton from "@/pages/tract-canvas/components/NewTractButton/NewTractButton.tsx"
import NewTractDialog from "@/pages/tract-canvas/dialogs/NewTractDialog/NewTractDialog.tsx"

export default function HeroSegment() {
    const {tracts, loading} = useTracts()
    const {OpenDialog} = useDialog()
    const active = tracts.filter(t => t.enabled).length
    const paused = tracts.length - active

    return (
        <div className={cls.HeroSegmentContainer}>
            <div className={cls.HeroTitles}>
                <div className={cls.Eyebrow}>Workspace · Automations</div>
                <h1 className={cls.HeroTitle}>Tracts</h1>
                <p className={cls.HeroSub}>
                    <b>{loading ? "…" : `${tracts.length} ${tracts.length === 1 ? "tract" : "tracts"}`}</b>
                    {tracts.length > 0 && <>{" · "}<span>{active} active, {paused} paused</span></>}
                </p>
            </div>
            <NewTractButton onClick={() => OpenDialog(<NewTractDialog/>)}/>
        </div>
    )
}
