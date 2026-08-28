import TweaksSection from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksSection/TweaksSection.tsx"
import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksMaxTokensSection/TweaksMaxTokensSection.module.css"

// Placeholder "Max tokens" row — a static, non-interactive track styled after the
// mock's `.slider-row`. No real <input type=range>; the row is marked
// aria-disabled with a "Coming soon" tooltip, matching the Stage 3/4 treatment.
export default function TweaksMaxTokensSection() {
    return (
        <TweaksSection label="Max tokens" disabled>
            <div
                className={cls.TweaksMaxTokensSectionContainer}
                aria-disabled="true"
                data-tooltip-id="root-tooltip"
                data-tooltip-content="Coming soon"
            >
                <span className={cls.Track}/>
                <span className={cls.Value}>4096</span>
            </div>
        </TweaksSection>
    )
}
