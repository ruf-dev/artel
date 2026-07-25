import {beforeEach, describe, expect, it, vi} from 'vitest'
import {render, screen, waitFor} from '@testing-library/react'

import {McpKeysAPI} from '@/app/api/artel/mcp_keys.pb.ts'
import {useMcpKeys} from '@/app/hooks/McpKeys.ts'
import CommunitySection from '@/pages/connections/components/CommunitySection/CommunitySection.tsx'

vi.mock('@/app/api/artel/mcp_keys.pb.ts', () => ({
    McpKeysAPI: {
        ListCommunityConnectors: vi.fn(),
        DeleteCommunityConnector: vi.fn(),
    },
}))

describe('CommunitySection', () => {
    beforeEach(() => {
        useMcpKeys.setState({communityConnectors: [], communityConnectorsLoading: false})
        vi.mocked(McpKeysAPI.ListCommunityConnectors).mockReset()
        vi.mocked(McpKeysAPI.DeleteCommunityConnector).mockReset()
    })

    it('renders a list of community connectors from a mocked fetch', async () => {
        vi.mocked(McpKeysAPI.ListCommunityConnectors).mockResolvedValue({
            connectors: [
                {
                    name: 'gitlab-community',
                    author: 'alice',
                    description: 'A community GitLab connector',
                    tools: [{name: 'list_issues'}, {name: 'create_issue'}],
                    viewerIsOwner: false,
                    viewerIsConnected: true,
                },
            ],
        })

        render(<CommunitySection/>)

        await waitFor(() => expect(screen.getByText('gitlab-community')).toBeInTheDocument())
        expect(screen.getByText('2 tools')).toBeInTheDocument()
    })

    it('shows an empty state when there are no connectors', async () => {
        vi.mocked(McpKeysAPI.ListCommunityConnectors).mockResolvedValue({connectors: []})

        render(<CommunitySection/>)

        await waitFor(() => expect(screen.getByText(/No community connectors yet/)).toBeInTheDocument())
    })
})
