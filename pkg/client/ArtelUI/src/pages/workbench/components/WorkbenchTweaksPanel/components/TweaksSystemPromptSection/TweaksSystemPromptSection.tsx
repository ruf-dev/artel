import {useEffect, useState} from "react"
import {Button, Toggle} from "@vervstack/chures"

import Textarea from "@/components/atoms/Textarea/Textarea.tsx"
import {useVaults, useVaultMutations} from "@/app/hooks/Vaults.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import TweaksSection from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksSection/TweaksSection.tsx"
import cls
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksSystemPromptSection/TweaksSystemPromptSection.module.css"

interface Props {
    vaultId?: string
}

// api-mode only: edit this vault's per-vault system prompt and whether the
// workspace-wide system prompt is layered on top. Seeds from the vault record and
// re-syncs if it changes underneath; Save persists via UpdateVaultPrompts.
export default function TweaksSystemPromptSection({vaultId}: Props) {
    const {vaults} = useVaults()
    const {updatePrompts} = useVaultMutations()
    const bakeError = useBakeError()

    const vault = vaults.find(v => v.id === vaultId)
    const savedPrompt = vault?.prompt ?? ""
    const savedUseSystemPrompt = vault?.useSystemPrompt ?? false

    const [prompt, setPrompt] = useState(savedPrompt)
    const [useSystemPrompt, setUseSystemPrompt] = useState(savedUseSystemPrompt)

    useEffect(() => {
        setPrompt(vault?.prompt ?? "")
    }, [vault?.prompt])

    useEffect(() => {
        setUseSystemPrompt(vault?.useSystemPrompt ?? false)
    }, [vault?.useSystemPrompt])

    const dirty = prompt !== savedPrompt || useSystemPrompt !== savedUseSystemPrompt

    function handleSave() {
        if (!vaultId) return
        updatePrompts(vaultId, prompt, useSystemPrompt)
            .catch(e => bakeError("Failed to save prompt", e))
    }

    return (
        <TweaksSection label="System prompt">
            <div className={cls.TweaksSystemPromptSectionContainer}>
                <Textarea
                    className={cls.PromptInput}
                    value={prompt}
                    setValue={setPrompt}
                    rows={4}
                    placeholder="You are a senior engineer working inside this vault…"
                />
                <Toggle
                    checked={useSystemPrompt}
                    onChange={setUseSystemPrompt}
                    label="Also apply workspace system prompt"
                />
                <Button
                    className={cls.SaveButton}
                    variant="primary"
                    onClick={handleSave}
                    disabled={!vaultId || !dirty}
                >
                    Save
                </Button>
            </div>
        </TweaksSection>
    )
}
