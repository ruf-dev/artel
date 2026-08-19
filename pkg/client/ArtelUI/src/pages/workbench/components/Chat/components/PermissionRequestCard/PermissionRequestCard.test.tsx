import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import PermissionRequestCard
    from "@/pages/workbench/components/Chat/components/PermissionRequestCard/PermissionRequestCard.tsx"

describe("PermissionRequestCard", () => {
    it("renders the three decision buttons and reports the clicked one", () => {
        const onDecide = vi.fn()
        render(
            <PermissionRequestCard
                toolName="Bash"
                inputSummary="rm -rf /tmp/x"
                options={["allow_once", "allow_always", "deny"]}
                onDecide={onDecide}
            />,
        )

        expect(screen.getByText("Permission requested: Bash")).toBeInTheDocument()
        expect(screen.getByText("rm -rf /tmp/x")).toBeInTheDocument()

        fireEvent.click(screen.getByText("Allow always"))
        expect(onDecide).toHaveBeenCalledWith("allow_always")
    })

    it("shows the resolved outcome instead of buttons once a decision is set", () => {
        render(
            <PermissionRequestCard
                toolName="Bash"
                options={["allow_once", "allow_always", "deny"]}
                decision="deny"
                onDecide={vi.fn()}
            />,
        )

        expect(screen.getByText("Denied")).toBeInTheDocument()
        expect(screen.queryByText("Allow")).not.toBeInTheDocument()
        expect(screen.queryByText("Deny")).not.toBeInTheDocument()
    })

    it("only renders buttons present in options", () => {
        render(
            <PermissionRequestCard
                toolName="Bash"
                options={["allow_once", "deny"]}
                onDecide={vi.fn()}
            />,
        )

        expect(screen.getByText("Allow")).toBeInTheDocument()
        expect(screen.getByText("Deny")).toBeInTheDocument()
        expect(screen.queryByText("Allow always")).not.toBeInTheDocument()
    })
})
