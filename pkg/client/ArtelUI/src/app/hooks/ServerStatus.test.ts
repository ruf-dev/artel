import {afterEach, beforeEach, describe, expect, it, mock} from "bun:test"

import {pingServer} from "@/app/hooks/ServerStatus"
import {useAppConfig} from "@/app/hooks/AppConfig"

const originalFetch = globalThis.fetch

afterEach(() => {
    globalThis.fetch = originalFetch
})

beforeEach(() => {
    useAppConfig.setState({noAuthEnabled: false})
})

describe("pingServer", () => {
    it("reports down when fetch rejects with a TypeError (network unreachable)", async () => {
        globalThis.fetch = mock(() => Promise.reject(new TypeError("Failed to fetch")))

        const result = await pingServer()

        expect(result).toBe(false)
    })

    it("reports down when fetch rejects with an AbortError DOMException (timeout)", async () => {
        globalThis.fetch = mock(() => Promise.reject(new DOMException("The operation was aborted", "AbortError")))

        const result = await pingServer()

        expect(result).toBe(false)
    })

    // Regression test for the black-screen bug: in local dev, when the Go
    // backend is down and there's no Vite proxy for /api, fetch('/api/auth/config')
    // resolves against Vite's own dev server, which serves its SPA index.html
    // fallback — 200 OK, but the body isn't valid JSON. fetchRequest() then
    // throws the raw Response (since r.json() failed to parse), and that must
    // be classified as "down", not "up".
    it("reports down when the response is 200 OK but the body isn't valid JSON (Vite SPA fallback)", async () => {
        const htmlResponse = new Response("<!doctype html><html>...</html>", {
            status: 200,
            headers: {"Content-Type": "text/html"},
        })
        globalThis.fetch = mock(() => Promise.resolve(htmlResponse))

        const result = await pingServer()

        expect(result).toBe(false)
    })

    it("reports up when the response is a non-OK status with a valid JSON error body", async () => {
        const errorResponse = new Response(JSON.stringify({message: "some error"}), {
            status: 500,
            headers: {"Content-Type": "application/json"},
        })
        globalThis.fetch = mock(() => Promise.resolve(errorResponse))

        const result = await pingServer()

        expect(result).toBe(true)
    })

    it("reports up and updates noAuthEnabled when the response is 200 OK with a valid JSON body", async () => {
        const okResponse = new Response(JSON.stringify({noAuthEnabled: true}), {
            status: 200,
            headers: {"Content-Type": "application/json"},
        })
        globalThis.fetch = mock(() => Promise.resolve(okResponse))

        const result = await pingServer()

        expect(result).toBe(true)
        expect(useAppConfig.getState().noAuthEnabled).toBe(true)
    })
})
