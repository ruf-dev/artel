import ThemeToggle from "@/components/ThemeToggle/ThemeToggle.tsx"
import TweaksSection from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksSection/TweaksSection.tsx"
import cls
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksThemeSection/TweaksThemeSection.module.css"

// Theme row of the Tweaks panel. Reuses the shared ThemeToggle (the same control
// promoted out of the global Topbar) rather than the mock's segmented Dark/Light
// picker, so a theme flip here matches the toggle on every other route.
export default function TweaksThemeSection() {
    return (
        <div className={cls.TweaksThemeSectionContainer}>
            <TweaksSection label="Theme">
                <ThemeToggle/>
            </TweaksSection>
        </div>
    )
}
