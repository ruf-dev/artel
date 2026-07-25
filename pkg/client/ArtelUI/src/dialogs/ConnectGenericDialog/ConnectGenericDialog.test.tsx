import {beforeEach, describe, expect, it, vi} from 'vitest'
import {render, screen, waitFor} from '@testing-library/react'
import userEvent from '@testing-library/user-event'

import {ExternalConnectionsAPI} from '@/app/api/artel/external_connections.pb.ts'
import {useExternalConnections} from '@/app/hooks/ExternalConnections.ts'
import ConnectGenericDialog from '@/dialogs/ConnectGenericDialog/ConnectGenericDialog.tsx'

vi.mock('@/app/api/artel/external_connections.pb.ts', () => ({
    ExternalConnectionsAPI: {
        AddGenericConnection: vi.fn(),
        ListConnections: vi.fn(),
    },
}))

describe('ConnectGenericDialog', () => {
    beforeEach(() => {
        useExternalConnections.setState({connections: [], loading: false})
        vi.mocked(ExternalConnectionsAPI.AddGenericConnection).mockReset()
        vi.mocked(ExternalConnectionsAPI.ListConnections).mockReset().mockResolvedValue({connections: []})
    })

    it('allows adding another credential row', async () => {
        render(<ConnectGenericDialog name="gitlab-community" onSuccess={vi.fn()}/>)
        const user = userEvent.setup()

        expect(screen.getAllByPlaceholderText('Field name (e.g. api_key)')).toHaveLength(1)
        await user.click(screen.getByRole('button', {name: '+ Add field'}))
        expect(screen.getAllByPlaceholderText('Field name (e.g. api_key)')).toHaveLength(2)
    })

    it('submits AddGenericConnection with the provider name and entered credentials', async () => {
        vi.mocked(ExternalConnectionsAPI.AddGenericConnection).mockResolvedValue({connection: {}})
        const onSuccess = vi.fn()
        render(<ConnectGenericDialog name="gitlab-community" onSuccess={onSuccess}/>)
        const user = userEvent.setup()

        await user.type(screen.getByPlaceholderText('Field name (e.g. api_key)'), 'api_key')
        await user.type(screen.getByPlaceholderText('Value'), 'secret-token')
        await user.click(screen.getByRole('button', {name: 'Connect'}))

        await waitFor(() => expect(ExternalConnectionsAPI.AddGenericConnection).toHaveBeenCalledWith(
            {provider: 'gitlab-community', credentials: {api_key: 'secret-token'}},
            expect.anything(),
        ))
        expect(onSuccess).toHaveBeenCalledTimes(1)
    })

    it('keeps Connect disabled until at least one field has a name', () => {
        render(<ConnectGenericDialog name="gitlab-community" onSuccess={vi.fn()}/>)
        expect(screen.getByRole('button', {name: 'Connect'})).toBeDisabled()
    })
})
