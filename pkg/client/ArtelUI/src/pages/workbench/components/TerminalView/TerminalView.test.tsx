import {describe, expect, it, vi} from 'vitest'
import {render, screen} from '@testing-library/react'

import TerminalView from '@/pages/workbench/components/TerminalView/TerminalView.tsx'

// Mock Terminal component to avoid its internal fetch-based implementation
vi.mock('@/pages/workbench/components/Terminal/Terminal.tsx', () => ({
    default: () => <span data-testid="terminal-component">Terminal</span>,
}))

// Mock TerminalTabBar and TerminalAuthLinkBanner as they render cleanly
vi.mock('@/pages/workbench/components/Terminal/components/TerminalTabBar/TerminalTabBar.tsx', () => ({
    default: () => <span data-testid="terminal-tab-bar">TerminalTabBar</span>,
}))

vi.mock('@/pages/workbench/components/Terminal/components/TerminalAuthLinkBanner/TerminalAuthLinkBanner.tsx', () => ({
    default: () => <span data-testid="terminal-auth-link-banner">TerminalAuthLinkBanner</span>,
}))

describe('TerminalView', () => {
    const mockProps = {
        vaultId: 'v1',
        onSelectTab: vi.fn(),
        onCreateTab: vi.fn(),
        onCloseTab: vi.fn(),
        pendingTerminalAuthLink: undefined,
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

    it('renders TerminalTabBar and TerminalAuthLinkBanner regardless of tab count', () => {
        render(
            <TerminalView
                {...mockProps}
                tabs={[]}
            />
        )

        expect(screen.getByTestId('terminal-tab-bar')).toBeInTheDocument()
        expect(screen.getByTestId('terminal-auth-link-banner')).toBeInTheDocument()
    })
})
