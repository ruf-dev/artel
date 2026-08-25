import {useState} from "react"
import {Button} from "@vervstack/chures"

import cls from "@/pages/workbench/components/PickWorkbenchModeScreen/PickWorkbenchModeScreen.module.css"
import OptionCard from "@/components/OptionCard/OptionCard.tsx"
import SimpleChatModeForm from
    "@/pages/workbench/components/PickWorkbenchModeScreen/components/SimpleChatModeForm/SimpleChatModeForm.tsx"

type Mode = "simple-chat" | "docker"

interface Props {
    vaultId: string
    onStartDocker: () => void
    startingDocker: boolean
    onSimpleChatCreated: (chatId: string) => void
}

// New "which workbench mode" picker, shown from WorkbenchPage's empty state
// (structurally similar to PickAuthModeScreen — local rows + footer — but not a
// modification of it, since that screen picks how Docker Claude Code
// authenticates, a step below this one). Simple Chat is the default selection:
// it has no Docker dependency and is the lighter-weight option.
export default function PickWorkbenchModeScreen({vaultId, onStartDocker, startingDocker, onSimpleChatCreated}: Props) {
    const [mode, setMode] = useState<Mode>("simple-chat")

    return (
        <div className={cls.PickWorkbenchModeScreenContainer}>
            <p className={cls.ModalSub}>Choose how you want to work in this vault.</p>
            <div className={cls.OptionList}>
                <OptionCard
                    selected={mode === "simple-chat"}
                    label="Simple Chat"
                    desc="Chat with an in-process agent that can call MCP tools. No Docker required."
                    onSelect={() => setMode("simple-chat")}
                />
                <OptionCard
                    selected={mode === "docker"}
                    label="Docker Claude Code"
                    desc="Full Claude Code workbench running in its own container, with a terminal."
                    onSelect={() => setMode("docker")}
                />
            </div>
            {mode === "simple-chat" && (
                <SimpleChatModeForm vaultId={vaultId} onCreated={onSimpleChatCreated}/>
            )}
            {mode === "docker" && (
                <div className={cls.ModalFooter}>
                    <Button variant="primary" onClick={onStartDocker} disabled={startingDocker}>
                        {startingDocker ? "Setting up…" : "Set up Workbench"}
                    </Button>
                </div>
            )}
        </div>
    )
}
