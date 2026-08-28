import {beforeEach, describe, expect, it, vi} from "vitest"
import {render} from "@testing-library/react"
import {act} from "react"

import Router, {Path, REDIRECT_AFTER_LOGIN_KEY, isPublicPath} from "@/app/routing/Router.tsx"

const navigate = vi.fn()
const isAuthenticated = vi.fn()
let mockPathname = "/"

vi.mock("react-router-dom", async (importOriginal) => ({
    ...(await importOriginal<object>()),
    useNavigate: () => navigate,
    useLocation: () => ({pathname: mockPathname, search: "", hash: "", state: null, key: "test"}),
    useRoutes: () => <span data-testid="route-element"/>,
}))

vi.mock("@/hooks/user/User.ts", () => ({
    default: () => ({
        auth: {isAuthenticated},
        setUserInfo: vi.fn(),
    }),
}))

vi.mock("@/app/hooks/ServerStatus.ts", () => ({
    useServerStatus: () => true,
}))

vi.mock("@/app/hooks/AppConfig.ts", () => ({
    useAppConfig: () => true,
}))

vi.mock("@/processes/Auth.ts", () => ({
    authService: {FetchUserInfo: vi.fn().mockResolvedValue({})},
}))

vi.mock("@/segments/UnsecureBanner/UnsecureBanner.tsx", () => ({default: () => null}))
vi.mock("@/pages/segments/Dialog.tsx", () => ({default: () => null}))
vi.mock("@/segments/Toaster/Toaster.tsx", () => ({default: () => null}))
vi.mock("react-tooltip", () => ({Tooltip: () => null}))

describe("isPublicPath", () => {
    it("treats auth / setup / docs paths as public", () => {
        expect(isPublicPath(Path.InitPage)).toBe(true)
        expect(isPublicPath(Path.SetupWizard)).toBe(true)
        expect(isPublicPath(Path.ClosedAlpha)).toBe(true)
        expect(isPublicPath(Path.McpAuth)).toBe(true)
        expect(isPublicPath(Path.DocsPageDefault)).toBe(true)
        expect(isPublicPath("/docs/getting-started")).toBe(true)
    })

    it("treats app paths as non-public", () => {
        expect(isPublicPath(Path.HomePage)).toBe(false)
        expect(isPublicPath("/notes")).toBe(false)
        expect(isPublicPath("/workbench/vault-1")).toBe(false)
    })
})

describe("Router auth guard", () => {
    beforeEach(() => {
        navigate.mockClear()
        isAuthenticated.mockReset()
        localStorage.clear()
        mockPathname = "/"
    })

    it("redirects an unauthenticated user off a protected path and stores it for post-login", () => {
        isAuthenticated.mockReturnValue(false)
        mockPathname = "/notes"

        render(<Router/>)

        expect(navigate).toHaveBeenCalledWith(Path.InitPage)
        expect(localStorage.getItem(REDIRECT_AFTER_LOGIN_KEY)).toBe("/notes")
    })

    it("does not redirect an unauthenticated user on a public docs path", () => {
        isAuthenticated.mockReturnValue(false)
        mockPathname = "/docs/some-slug"

        render(<Router/>)

        expect(navigate).not.toHaveBeenCalled()
        expect(localStorage.getItem(REDIRECT_AFTER_LOGIN_KEY)).toBeNull()
    })

    it("does not redirect an unauthenticated user on /init", () => {
        isAuthenticated.mockReturnValue(false)
        mockPathname = Path.InitPage

        render(<Router/>)

        expect(navigate).not.toHaveBeenCalled()
    })

    it("does not redirect an authenticated user on a protected path", () => {
        isAuthenticated.mockReturnValue(true)
        mockPathname = "/notes"

        render(<Router/>)

        expect(navigate).not.toHaveBeenCalled()
        expect(localStorage.getItem(REDIRECT_AFTER_LOGIN_KEY)).toBeNull()
    })

    it("redirects to /init when auth is lost without a route change (logout in place)", () => {
        isAuthenticated.mockReturnValue(true)
        mockPathname = "/notes"

        const {rerender} = render(<Router/>)
        expect(navigate).not.toHaveBeenCalled()

        // logout() swaps in a fresh AuthMiddleware — the useUser mock returns a
        // new `auth` object each render, so a rerender with isAuthenticated now
        // false stands in for that, pathname unchanged.
        isAuthenticated.mockReturnValue(false)
        act(() => {
            rerender(<Router/>)
        })

        expect(navigate).toHaveBeenCalledWith(Path.InitPage)
        expect(localStorage.getItem(REDIRECT_AFTER_LOGIN_KEY)).toBe("/notes")
    })
})
