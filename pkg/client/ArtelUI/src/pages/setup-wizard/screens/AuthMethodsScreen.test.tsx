import {beforeEach, describe, expect, it, vi} from 'vitest'
import {fireEvent, render, screen, waitFor} from '@testing-library/react'

import {SetupWizardAPI} from '@/app/api/artel/setup_wizard.pb.ts'
import {SetupWizardContext, type SetupWizardState} from '@/pages/setup-wizard/setupWizardContext.ts'
import {RegistrationMode} from '@/app/api/artel'
import AuthMethodsScreen from '@/pages/setup-wizard/screens/AuthMethodsScreen.tsx'

vi.mock('@/app/api/artel/setup_wizard.pb.ts', async (importOriginal) => ({
    ...(await importOriginal<object>()),
    SetupWizardAPI: {
        SelectAuthMethods: vi.fn(),
    },
}))

const bakeError = vi.fn()
vi.mock('@/app/hooks/useErrorToast', () => ({
    useBakeError: () => bakeError,
}))

function buildContextValue(overrides: Partial<SetupWizardState> = {}): SetupWizardState {
    return {
        step: 'methods',
        setStep: vi.fn(),

        wizardSessionToken: 'test-token',
        setWizardSessionToken: vi.fn(),

        passwordEnabled: true,
        setPasswordEnabled: vi.fn(),
        telegramEnabled: false,
        setTelegramEnabled: vi.fn(),

        registrationMode: RegistrationMode.ADMIN_ONLY,
        setRegistrationMode: vi.fn(),

        ...overrides,
    }
}

function renderWithContext(overrides: Partial<SetupWizardState> = {}) {
    const ctx = buildContextValue(overrides)
    render(
        <SetupWizardContext.Provider value={ctx}>
            <AuthMethodsScreen/>
        </SetupWizardContext.Provider>,
    )
    return ctx
}

describe('AuthMethodsScreen', () => {
    beforeEach(() => {
        vi.mocked(SetupWizardAPI.SelectAuthMethods).mockReset()
        bakeError.mockReset()
    })

    it('calls SelectAuthMethods with the current toggle state and advances on success', async () => {
        vi.mocked(SetupWizardAPI.SelectAuthMethods).mockResolvedValue({})

        const ctx = renderWithContext({
            wizardSessionToken: 'session-abc',
            passwordEnabled: true,
            telegramEnabled: false,
        })

        fireEvent.click(screen.getByText('Continue'))

        await waitFor(() => expect(SetupWizardAPI.SelectAuthMethods).toHaveBeenCalledWith(
            {wizardSessionToken: 'session-abc', passwordEnabled: true, telegramEnabled: false},
            expect.anything(),
        ))
        await waitFor(() => expect(ctx.setStep).toHaveBeenCalledWith('mode'))
    })

    it('does not advance and surfaces an error when SelectAuthMethods rejects', async () => {
        vi.mocked(SetupWizardAPI.SelectAuthMethods).mockRejectedValue(new Error('boom'))

        const ctx = renderWithContext()

        fireEvent.click(screen.getByText('Continue'))

        await waitFor(() => expect(bakeError).toHaveBeenCalledWith('Failed to save sign-in methods', expect.any(Error)))
        expect(ctx.setStep).not.toHaveBeenCalled()
    })
})
