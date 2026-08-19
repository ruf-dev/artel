import cls from "@/pages/workbench/components/WorkbenchAuthScreen/WorkbenchAuthScreen.module.css"
import AuthSignInPanel
    from "@/pages/workbench/components/WorkbenchAuthScreen/components/AuthSignInPanel/AuthSignInPanel.tsx"
import AuthCodeForm from "@/pages/workbench/components/WorkbenchAuthScreen/components/AuthCodeForm/AuthCodeForm.tsx"
import ChatStatusBanner from "@/pages/workbench/components/Chat/components/ChatStatusBanner/ChatStatusBanner.tsx"
import {ChatItem} from "@/pages/workbench/processes/chatReducer.ts"
import {ChatConnectionStatus} from "@/pages/workbench/processes/useChatSession.ts"

interface Props {
    items: ChatItem[]
    status: ChatConnectionStatus
    onSubmitCode: (code: string) => void
}

export default function WorkbenchAuthScreen({items, status, onSubmitCode}: Props) {
    // Latest, not first: a reconnect after the underlying container/bridge process restarts
    // (same vaultId, so items isn't reset) can leave a stale auth_link from a dead
    // `claude setup-token` run in the list — its state/code_challenge no longer matches any
    // pending flow server-side, so opening it fails as a CSRF/state error.
    const linkUrl = [...items].reverse().find(i => i.kind === "auth_link")?.url
    const codeNeeded = items.some(i => i.kind === "auth_code_needed" && !i.resolved)
    const lastError = [...items].reverse().find(i => i.kind === "error")?.text

    return (
        <div className={cls.WorkbenchAuthScreenContainer}>
            <p className={cls.Heading}>Sign in with your Claude subscription to continue.</p>
            <AuthSignInPanel url={linkUrl}/>
            <AuthCodeForm onSubmit={onSubmitCode} disabled={!codeNeeded || status !== "open"}/>
            <ChatStatusBanner status={status}/>
            {lastError && <p className={cls.ErrorText}>{lastError}</p>}
        </div>
    )
}
