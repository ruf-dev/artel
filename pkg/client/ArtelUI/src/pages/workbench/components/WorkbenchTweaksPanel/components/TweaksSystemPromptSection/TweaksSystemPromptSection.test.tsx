import {beforeEach, describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import TweaksSystemPromptSection
    // eslint-disable-next-line max-len -- path too long to wrap
    from "@/pages/workbench/components/WorkbenchTweaksPanel/components/TweaksSystemPromptSection/TweaksSystemPromptSection.tsx"

const updatePrompts = vi.fn()
const bakeError = vi.fn()
let vaults: Array<{id: string; prompt?: string; useSystemPrompt?: boolean}>

vi.mock("@/app/hooks/Vaults.ts", () => ({
    useVaults: () => ({vaults}),
    useVaultMutations: () => ({updatePrompts}),
}))

vi.mock("@/app/hooks/useErrorToast.ts", () => ({
    useBakeError: () => bakeError,
}))

describe("TweaksSystemPromptSection", () => {
    beforeEach(() => {
        updatePrompts.mockReset().mockResolvedValue(undefined)
        bakeError.mockReset()
        vaults = [{id: "v1", prompt: "seeded prompt", useSystemPrompt: false}]
    })

    it("seeds the textarea from the vault's saved prompt", () => {
        render(<TweaksSystemPromptSection vaultId="v1"/>)

        expect(screen.getByRole("textbox")).toHaveValue("seeded prompt")
    })

    it("disables Save while pristine", () => {
        render(<TweaksSystemPromptSection vaultId="v1"/>)

        expect(screen.getByRole("button", {name: "Save"})).toBeDisabled()
    })

    it("saves the edited prompt via updatePrompts", () => {
        render(<TweaksSystemPromptSection vaultId="v1"/>)

        fireEvent.change(screen.getByRole("textbox"), {target: {value: "new prompt"}})
        fireEvent.click(screen.getByRole("button", {name: "Save"}))

        expect(updatePrompts).toHaveBeenCalledWith("v1", "new prompt", false)
    })
})
