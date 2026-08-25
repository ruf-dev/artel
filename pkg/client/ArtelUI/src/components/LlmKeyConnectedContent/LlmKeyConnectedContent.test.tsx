/* eslint-disable func-style, max-lines-per-function */
import {describe, expect, it} from 'vitest'
import {render, screen} from '@testing-library/react'

import LlmKeyConnectedContent from '@/components/LlmKeyConnectedContent/LlmKeyConnectedContent.tsx'

describe('LlmKeyConnectedContent', () => {
    it('does not render available models section when fields.available_models is not present', () => {
        const fields = {
            key_preview: 'sk-test-key',
            default_model: 'openai/gpt-4',
        }
        const onDisconnect = () => {}

        render(<LlmKeyConnectedContent fields={fields} onDisconnect={onDisconnect} />)

        expect(screen.queryByText(/Available models/)).not.toBeInTheDocument()
    })

    it('does not render available models section when fields.available_models is empty', () => {
        const fields = {
            key_preview: 'sk-test-key',
            default_model: 'openai/gpt-4',
            available_models: '',
        }
        const onDisconnect = () => {}

        render(<LlmKeyConnectedContent fields={fields} onDisconnect={onDisconnect} />)

        expect(screen.queryByText(/Available models/)).not.toBeInTheDocument()
    })

    it('renders available models section with correct count and dropdown', () => {
        const fields = {
            key_preview: 'sk-test-key',
            default_model: 'openai/gpt-4',
            available_models: 'openai/gpt-4,openai/gpt-3.5-turbo,anthropic/claude-3-opus',
        }
        const onDisconnect = () => {}

        render(<LlmKeyConnectedContent fields={fields} onDisconnect={onDisconnect} />)

        expect(screen.getByText('Available models (3)')).toBeInTheDocument()
        expect(screen.getByText('Browse models…')).toBeInTheDocument()
    })

    it('renders default_model field when present', () => {
        const fields = {
            key_preview: 'sk-test-key',
            default_model: 'openai/gpt-4',
            available_models: 'openai/gpt-4,openai/gpt-3.5-turbo',
        }
        const onDisconnect = () => {}

        render(<LlmKeyConnectedContent fields={fields} onDisconnect={onDisconnect} />)

        expect(screen.getByText('Default model')).toBeInTheDocument()
        expect(screen.getByText('openai/gpt-4')).toBeInTheDocument()
    })

    it('does not render default_model when not present', () => {
        const fields = {
            key_preview: 'sk-test-key',
            available_models: 'openai/gpt-4,openai/gpt-3.5-turbo',
        }
        const onDisconnect = () => {}

        render(<LlmKeyConnectedContent fields={fields} onDisconnect={onDisconnect} />)

        expect(screen.queryByText('Default model')).not.toBeInTheDocument()
    })

    it('renders key preview with connected message', () => {
        const fields = {
            key_preview: 'sk-test-key-12345',
        }
        const onDisconnect = () => {}

        render(<LlmKeyConnectedContent fields={fields} onDisconnect={onDisconnect} />)

        expect(screen.getByText(/Connected with key/)).toBeInTheDocument()
        expect(screen.getByText('sk-test-key-12345')).toBeInTheDocument()
    })

    it('renders disconnect button', () => {
        const fields = {
            key_preview: 'sk-test-key',
        }
        const onDisconnect = () => {}

        render(<LlmKeyConnectedContent fields={fields} onDisconnect={onDisconnect} />)

        expect(screen.getByRole('button', {name: /Disconnect/})).toBeInTheDocument()
    })

    it('renders view usage button when onViewUsage is provided', () => {
        const fields = {
            key_preview: 'sk-test-key',
        }
        const onDisconnect = () => {}
        const onViewUsage = () => {}

        render(<LlmKeyConnectedContent fields={fields} onDisconnect={onDisconnect} onViewUsage={onViewUsage} />)

        expect(screen.getByRole('button', {name: /View usage/})).toBeInTheDocument()
    })

    it('does not render view usage button when onViewUsage is not provided', () => {
        const fields = {
            key_preview: 'sk-test-key',
        }
        const onDisconnect = () => {}

        render(<LlmKeyConnectedContent fields={fields} onDisconnect={onDisconnect} />)

        expect(screen.queryByRole('button', {name: /View usage/})).not.toBeInTheDocument()
    })

    it('uses interactive placeholder when onSelectDefaultModel is provided', () => {
        const fields = {
            key_preview: 'sk-test-key',
            available_models: 'openai/gpt-4,openai/gpt-3.5-turbo',
        }
        const onDisconnect = () => {}
        const onSelectDefaultModel = () => Promise.resolve()

        render(
            <LlmKeyConnectedContent
                fields={fields}
                onDisconnect={onDisconnect}
                onSelectDefaultModel={onSelectDefaultModel}
            />,
        )

        const dropdown = screen.getByRole('button', {expanded: false})
        expect(dropdown).toHaveTextContent('Choose default model…')
    })

    it('uses browse-only placeholder when onSelectDefaultModel is not provided', () => {
        const fields = {
            key_preview: 'sk-test-key',
            available_models: 'openai/gpt-4,openai/gpt-3.5-turbo',
        }
        const onDisconnect = () => {}

        render(<LlmKeyConnectedContent fields={fields} onDisconnect={onDisconnect} />)

        const dropdown = screen.getByRole('button', {expanded: false})
        expect(dropdown).toHaveTextContent('Browse models…')
    })
})
