import TweaksSection from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksSection/TweaksSection.tsx"
import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksContextSection/TweaksContextSection.module.css"

// Placeholder "Context window" row — a static usage bar styled after the mock's
// `.usage-bar` / `.usage-label`. Non-interactive: aria-disabled with a "Coming
// soon" tooltip, matching the Stage 3/4 treatment.
export default function TweaksContextSection() {
    return (
        <TweaksSection label="Context window" disabled>
            <div
                className={cls.TweaksContextSectionContainer}
                aria-disabled="true"
                data-tooltip-id="root-tooltip"
                data-tooltip-content="Coming soon"
            >
                <span className={cls.Bar}/>
                <span className={cls.Label}>Context usage appears here once wired up</span>
            </div>
        </TweaksSection>
    )
}
