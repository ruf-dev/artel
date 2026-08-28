import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import SidebarBrand from "@/pages/workbench/components/WorkbenchSidebar/components/SidebarBrand/SidebarBrand.tsx"

const navigate = vi.fn()

vi.mock("react-router-dom", () => ({
    useNavigate: () => navigate,
}))

vi.mock("@/app/routing/Router.tsx", () => ({
    Path: {HomePage: "/"},
}))

describe("SidebarBrand", () => {
    it("navigates home when the wordmark is clicked", () => {
        render(<SidebarBrand/>)

        fireEvent.click(screen.getByRole("button", {name: "Go to home"}))

        expect(navigate).toHaveBeenCalledWith("/")
    })

    it("does not render the close button without showClose", () => {
        render(<SidebarBrand/>)

        expect(screen.queryByRole("button", {name: "Close chat"})).not.toBeInTheDocument()
    })

    it("renders the close button and navigates home on click when showClose is set", () => {
        render(<SidebarBrand showClose/>)

        fireEvent.click(screen.getByRole("button", {name: "Close chat"}))

        expect(navigate).toHaveBeenCalledWith("/")
    })
})
