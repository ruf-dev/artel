import {describe, expect, it, vi} from 'vitest'
import {render, screen} from '@testing-library/react'

import TerminalView from '@/pages/workbench/components/TerminalView/TerminalView.tsx'

// Mock Terminal component to avoid its internal fetch-based implementation
vi.mock('@/pages/workbench/components/Terminal/Terminal.tsx', () => ({
    default: () => <span data-testid="terminal-component">Terminal</span>,
}))

// Mock TerminalTabBar as it renders cleanly on its own
vi.mock('@/pages/workbench/components/Terminal/components/TerminalTabBar/TerminalTabBar.tsx', () => ({
    default: () => <span data-testid="terminal-tab-bar">TerminalTabBar</span>,
}))

describe('TerminalView', () => {
    const mockProps = {
        vaultId: 'v1',
        onSelectTab: vi.fn(),
        onCreateTab: vi.fn(),
        onCloseTab: vi.fn(),
    }

    it('renders empty state when tabs array is empty', () => {
        render(
            <TerminalView
                {...mockProps}
                tabs={[]}
            />
        )

        expect(screen.getByText('No terminal sessions yet. Click + to start one.')).toBeInTheDocument()
        expect(screen.queryByTestId('terminal-component')).not.toBeInTheDocument()
    })

    it('renders Terminal component when tabs array is not empty', () => {
        render(
            <TerminalView
                {...mockProps}
                tabs={[
                    {id: '@1', name: 'claude', active: true},
                ]}
            />
        )

        expect(screen.queryByText('No terminal sessions yet. Click + to start one.')).not.toBeInTheDocument()
        expect(screen.getByTestId('terminal-component')).toBeInTheDocument()
    })

    it('renders TerminalTabBar regardless of tab count', () => {
        render(
            <TerminalView
                {...mockProps}
                tabs={[]}
            />
        )

        expect(screen.getByTestId('terminal-tab-bar')).toBeInTheDocument()
    })
})
