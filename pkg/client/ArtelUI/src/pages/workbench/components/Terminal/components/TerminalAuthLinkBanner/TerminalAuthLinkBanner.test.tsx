import {describe, expect, it} from 'vitest'
import {render, screen} from '@testing-library/react'
import userEvent from '@testing-library/user-event'

import TerminalAuthLinkBanner
    from "@/pages/workbench/components/Terminal/components/TerminalAuthLinkBanner/TerminalAuthLinkBanner"

describe('TerminalAuthLinkBanner', () => {
    it('renders nothing when url is not provided', () => {
        const {container} = render(<TerminalAuthLinkBanner/>)

        expect(container).toBeEmptyDOMElement()
    })

    it('renders the link with the correct href when url is present', () => {
        render(<TerminalAuthLinkBanner url="https://claude.ai/login/abc"/>)

        const link = screen.getByRole('link', {name: 'Authorize Claude'})
        expect(link).toHaveAttribute('href', 'https://claude.ai/login/abc')
        expect(link).toHaveAttribute('target', '_blank')
    })

    it('hides the banner after the dismiss button is clicked', async () => {
        const user = userEvent.setup()
        const {container} = render(<TerminalAuthLinkBanner url="https://claude.ai/login/abc"/>)

        await user.click(screen.getByLabelText('Dismiss'))

        expect(container).toBeEmptyDOMElement()
    })

    it('shows a new url again after the previous one was dismissed', async () => {
        const user = userEvent.setup()
        const {rerender, container} = render(<TerminalAuthLinkBanner url="https://claude.ai/login/abc"/>)

        await user.click(screen.getByLabelText('Dismiss'))
        expect(container).toBeEmptyDOMElement()

        rerender(<TerminalAuthLinkBanner url="https://claude.ai/login/xyz"/>)

        const link = screen.getByRole('link', {name: 'Authorize Claude'})
        expect(link).toHaveAttribute('href', 'https://claude.ai/login/xyz')
    })
})
