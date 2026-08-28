import {describe, expect, it, vi} from "vitest"
import {fireEvent, render, screen} from "@testing-library/react"

import StartStopButton
    from "@/pages/workbench/components/WorkbenchTopbar/components/StartStopButton/StartStopButton.tsx"

const base = {isRunning: false, onStart: vi.fn(), onStop: vi.fn(), stopping: false, starting: false}

describe("StartStopButton", () => {
    it("calls onStart from the Play state", () => {
        const onStart = vi.fn()
        render(<StartStopButton {...base} onStart={onStart}/>)

        fireEvent.click(screen.getByRole("button", {name: "Start"}))

        expect(onStart).toHaveBeenCalledTimes(1)
    })

    it("calls onStop from the Square state", () => {
        const onStop = vi.fn()
        render(<StartStopButton {...base} isRunning onStop={onStop}/>)

        fireEvent.click(screen.getByRole("button", {name: "Stop"}))

        expect(onStop).toHaveBeenCalledTimes(1)
    })

    it("is disabled and labelled \"Stopping\" while stopping", () => {
        render(<StartStopButton {...base} isRunning stopping/>)

        expect(screen.getByRole("button", {name: "Stopping"})).toBeDisabled()
    })

    it("is disabled while starting", () => {
        render(<StartStopButton {...base} starting/>)

        expect(screen.getByRole("button", {name: "Start"})).toBeDisabled()
    })
})
