import {beforeEach, describe, expect, it, vi} from 'vitest'
import {render, screen} from '@testing-library/react'

import TerminalTabChip
    from "@/pages/workbench/components/Terminal/components/TerminalTabBar/components/TerminalTabChip/TerminalTabChip"

describe('TerminalTabChip', () => {
    const mockTab = {id: 'tab1', name: 'Main', active: true}

    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('renders tab name', () => {
        const onSelect = vi.fn()
        const onClose = vi.fn()

        render(
            <TerminalTabChip
                tab={mockTab}
                canClose={true}
                onSelect={onSelect}
                onClose={onClose}
            />
        )

        expect(screen.getByText('Main')).toBeInTheDocument()
    })

    it('renders close button', () => {
        const onSelect = vi.fn()
        const onClose = vi.fn()

        render(
            <TerminalTabChip
                tab={mockTab}
                canClose={true}
                onSelect={onSelect}
                onClose={onClose}
            />
        )

        expect(screen.getByLabelText('Close tab')).toBeInTheDocument()
    })

    it('disables close button when canClose is false', () => {
        const onSelect = vi.fn()
        const onClose = vi.fn()

        render(
            <TerminalTabChip
                tab={mockTab}
                canClose={false}
                onSelect={onSelect}
                onClose={onClose}
            />
        )

        const closeButton = screen.getByLabelText('Close tab')
        expect(closeButton).toBeDisabled()
    })

    it('enables close button when canClose is true', () => {
        const onSelect = vi.fn()
        const onClose = vi.fn()

        render(
            <TerminalTabChip
                tab={mockTab}
                canClose={true}
                onSelect={onSelect}
                onClose={onClose}
            />
        )

        const closeButton = screen.getByLabelText('Close tab')
        expect(closeButton).not.toBeDisabled()
    })

    it('renders with active class when tab.active is true', () => {
        const onSelect = vi.fn()
        const onClose = vi.fn()

        const {container} = render(
            <TerminalTabChip
                tab={mockTab}
                canClose={true}
                onSelect={onSelect}
                onClose={onClose}
            />
        )

        const chipContainer = container.firstChild as HTMLElement
        expect(chipContainer.className).toMatch(/Active/)
    })

    it('renders without active class when tab.active is false', () => {
        const onSelect = vi.fn()
        const onClose = vi.fn()
        const inactiveTab = {id: 'tab2', name: 'Logs', active: false}

        const {container} = render(
            <TerminalTabChip
                tab={inactiveTab}
                canClose={true}
                onSelect={onSelect}
                onClose={onClose}
            />
        )

        const chipContainer = container.firstChild as HTMLElement
        expect(chipContainer.className).not.toMatch(/Active/)
    })
})
