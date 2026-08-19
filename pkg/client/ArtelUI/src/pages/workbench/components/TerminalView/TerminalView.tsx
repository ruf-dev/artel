import Terminal from "@/pages/workbench/components/Terminal/Terminal.tsx"
import TerminalTabBar from "@/pages/workbench/components/Terminal/components/TerminalTabBar/TerminalTabBar.tsx"
import cls from "@/pages/workbench/components/TerminalView/TerminalView.module.css"

interface Props {
    vaultId: string
    tabs: {id: string; name: string; active: boolean}[]
    onSelectTab: (id: string) => void
    onCreateTab: () => void
    onCloseTab: (id: string) => void
}

export default function TerminalView(props: Props) {
    return (
        <div className={cls.TerminalViewContainer}>
            <div className={cls.TerminalPanel}>
                <TerminalTabBar
                    tabs={props.tabs}
                    onSelect={props.onSelectTab}
                    onCreate={props.onCreateTab}
                    onClose={props.onCloseTab}
                />
                <Terminal vaultId={props.vaultId}/>
            </div>
            <div className={cls.TerminalSpacer}/>
        </div>
    )
}
