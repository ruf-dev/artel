import {beforeEach, describe, expect, it, vi} from 'vitest'
import {render, screen} from '@testing-library/react'

import TerminalTabBar
    from '@/pages/workbench/components/Terminal/components/TerminalTabBar/TerminalTabBar.tsx'

function renderTabBar(tabs = [
    {id: 'tab1', name: 'Main', active: true},
    {id: 'tab2', name: 'Logs', active: false},
]) {
    const onSelect = vi.fn()
    const onCreate = vi.fn()
    const onClose = vi.fn()
    render(
        <TerminalTabBar
            tabs={tabs}
            onSelect={onSelect}
            onCreate={onCreate}
            onClose={onClose}
        />
    )
    return {onSelect, onCreate, onClose}
}

describe('TerminalTabBar', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('renders all tabs', () => {
        renderTabBar()
        expect(screen.getByText('Main')).toBeInTheDocument()
        expect(screen.getByText('Logs')).toBeInTheDocument()
    })

    it('renders new tab button', () => {
        renderTabBar()
        expect(screen.getByLabelText('New terminal tab')).toBeInTheDocument()
    })

    it('renders correct number of close buttons', () => {
        renderTabBar()
        const closeButtons = screen.getAllByLabelText('Close tab')
        expect(closeButtons).toHaveLength(2)
    })

    it('disables close buttons when there is only one tab', () => {
        renderTabBar([{id: 'tab1', name: 'Main', active: true}])
        const closeButton = screen.getByLabelText('Close tab')
        expect(closeButton).toBeDisabled()
    })

    it('allows closing when there are multiple tabs', () => {
        renderTabBar()
        const closeButtons = screen.getAllByLabelText('Close tab')
        closeButtons.forEach(btn => {
            expect(btn).not.toBeDisabled()
        })
    })
})
