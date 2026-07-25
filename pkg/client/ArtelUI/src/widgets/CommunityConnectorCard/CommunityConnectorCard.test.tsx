import {describe, expect, it, vi} from 'vitest'
import {render, screen} from '@testing-library/react'
import userEvent from '@testing-library/user-event'

import CommunityConnectorCard from '@/widgets/CommunityConnectorCard/CommunityConnectorCard.tsx'

interface RenderOverrides {
    viewerIsOwner?: boolean
    isConnected?: boolean
}

function renderCard(overrides: RenderOverrides = {}) {
    const onConnect = vi.fn()
    const onDelete = vi.fn().mockResolvedValue(undefined)
    render(
        <CommunityConnectorCard
            name="gitlab-community"
            author="alice"
            description="A community GitLab connector"
            toolCount={3}
            isConnected={overrides.isConnected ?? false}
            viewerIsOwner={overrides.viewerIsOwner ?? false}
            onConnect={onConnect}
            onDelete={onDelete}
        />
    )
    return {onConnect, onDelete}
}

describe('CommunityConnectorCard', () => {
    it('renders name, author, description and tool count', () => {
        renderCard()
        expect(screen.getByText('gitlab-community')).toBeInTheDocument()
        expect(screen.getByText('by alice')).toBeInTheDocument()
        expect(screen.getByText('A community GitLab connector')).toBeInTheDocument()
        expect(screen.getByText('3 tools')).toBeInTheDocument()
    })

    it('hides the Delete button when the viewer is not the owner', () => {
        renderCard({viewerIsOwner: false})
        expect(screen.queryByRole('button', {name: 'Delete'})).not.toBeInTheDocument()
    })

    it('shows the Delete button when the viewer is the owner', () => {
        renderCard({viewerIsOwner: true})
        expect(screen.getByRole('button', {name: 'Delete'})).toBeInTheDocument()
    })

    it('always shows Connect and calls onConnect when clicked', async () => {
        const {onConnect} = renderCard()
        const user = userEvent.setup()
        await user.click(screen.getByRole('button', {name: 'Connect'}))
        expect(onConnect).toHaveBeenCalledTimes(1)
    })
})
