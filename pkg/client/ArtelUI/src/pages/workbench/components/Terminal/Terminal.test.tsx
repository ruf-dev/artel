import {beforeEach, describe, expect, it, vi} from 'vitest'
import {render, screen, waitFor} from '@testing-library/react'

import Terminal from '@/pages/workbench/components/Terminal/Terminal.tsx'

const bakeError = vi.fn()
vi.mock('@/app/hooks/useErrorToast', () => ({
    useBakeError: () => bakeError,
}))

describe('Terminal', () => {
    beforeEach(() => {
        bakeError.mockReset()
        vi.stubGlobal('fetch', vi.fn())
    })

    it('renders the iframe when fetch succeeds with ok: true', async () => {
        vi.mocked(fetch).mockResolvedValue({
            ok: true,
        } as Response)

        render(<Terminal vaultId="v1"/>)

        await waitFor(() => {
            expect(screen.getByTitle('Workbench terminal')).toBeInTheDocument()
        })
    })

    it('shows error state and calls bakeError when fetch fails with ok: false', async () => {
        vi.mocked(fetch).mockResolvedValue({
            ok: false,
            text: () => Promise.resolve('boom'),
        } as Response)

        render(<Terminal vaultId="v1"/>)

        await waitFor(() => {
            expect(bakeError).toHaveBeenCalledWith(
                'Failed to connect to workbench terminal',
                expect.any(Error),
            )
        })

        expect(screen.queryByTitle('Workbench terminal')).not.toBeInTheDocument()
        expect(screen.getByText('Retry')).toBeInTheDocument()
    })

    it('shows error state when fetch rejects', async () => {
        vi.mocked(fetch).mockRejectedValue(new Error('network error'))

        render(<Terminal vaultId="v1"/>)

        await waitFor(() => {
            expect(bakeError).toHaveBeenCalledWith(
                'Failed to connect to workbench terminal',
                expect.any(Error),
            )
        })

        expect(screen.queryByTitle('Workbench terminal')).not.toBeInTheDocument()
        expect(screen.getByText('Retry')).toBeInTheDocument()
    })
})
