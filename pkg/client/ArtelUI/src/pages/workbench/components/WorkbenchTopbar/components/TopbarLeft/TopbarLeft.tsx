import ModelSwitcher from "@/pages/workbench/components/ModelSwitcher/ModelSwitcher.tsx"
import WorkbenchStatusBadge from "@/pages/workbench/components/WorkbenchStatusBadge/WorkbenchStatusBadge.tsx"
import type {WorkbenchMode} from "@/pages/workbench/processes/useWorkbenchModeControls.ts"
import cls from "@/pages/workbench/components/WorkbenchTopbar/components/TopbarLeft/TopbarLeft.module.css"

interface ModelProps {
    models: string[]
    value: string
    isLoading?: boolean
    onChange: (model: string) => void
}

interface Props {
    effectiveMode: WorkbenchMode | "picking"
    exists: boolean
    status: string
    model: ModelProps
}

// Left cluster of WorkbenchTopbar: docker status badge and the model switcher
// (live in api mode, a fixed dimmed "Claude Code" in docker mode).
export default function TopbarLeft({effectiveMode, exists, status, model}: Props) {
    return (
        <div className={cls.TopbarLeftContainer}>
            {effectiveMode === "docker" && (
                <WorkbenchStatusBadge status={exists ? status : "not_configured"}/>
            )}
            {effectiveMode === "api" && (
                <ModelSwitcher
                    models={model.models}
                    value={model.value}
                    isLoading={model.isLoading}
                    onChange={model.onChange}
                    needsAttention={!model.isLoading && !model.value}
                />
            )}
            {effectiveMode === "docker" && (
                <ModelSwitcher disabled models={[]} value="" placeholder="Claude Code" onChange={() => {}}/>
            )}
        </div>
    )
}
