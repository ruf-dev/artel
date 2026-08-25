import SimpleChatToolbarActions from
    "@/pages/workbench/components/SimpleChatTopBar/components/SimpleChatToolbarActions/SimpleChatToolbarActions.tsx"
import WorkbenchTopBarShell from "@/pages/workbench/components/WorkbenchTopBarShell/WorkbenchTopBarShell.tsx"

interface Props {
    vaultName: string
    models: string[]
    currentModel: string
    modelsLoading?: boolean
    onChangeModel: (model: string) => void
    onNewChat: () => void
    onToggleHistory: () => void
}

// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function SimpleChatTopBar(props: Props) {
    return (
        <WorkbenchTopBarShell
            vaultName={props.vaultName}
            actions={(
                <SimpleChatToolbarActions
                    models={props.models}
                    currentModel={props.currentModel}
                    modelsLoading={props.modelsLoading}
                    onChangeModel={props.onChangeModel}
                    onNewChat={props.onNewChat}
                    onToggleHistory={props.onToggleHistory}
                />
            )}
        />
    )
}
