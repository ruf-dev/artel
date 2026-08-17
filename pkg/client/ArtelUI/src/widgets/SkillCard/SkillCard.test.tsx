import {beforeEach, describe, expect, it, vi} from 'vitest'
import {render, screen} from '@testing-library/react'
import userEvent from '@testing-library/user-event'

import {SkillInfo} from '@/app/api/artel/skills.pb.ts'
import {useSkills} from '@/app/hooks/Skills.ts'
import SkillCard from '@/widgets/SkillCard/SkillCard.tsx'

vi.mock('@/app/hooks/Skills.ts', () => ({
    useSkills: vi.fn(),
}))

vi.mock('@/dialogs/EditSkillDialog/EditSkillDialog.tsx', () => ({
    default: () => null,
}))

function renderCard(skill: Partial<SkillInfo>) {
    const setHotPlug = vi.fn().mockResolvedValue(undefined)
    const remove = vi.fn().mockResolvedValue(undefined)
    vi.mocked(useSkills).mockReturnValue({setHotPlug, remove} as unknown as ReturnType<typeof useSkills>)

    render(<SkillCard vaultId="vault-1" skill={{slug: 'foo', name: 'Foo', description: 'desc', ...skill}}/>)
    return {setHotPlug, remove}
}

describe('SkillCard', () => {
    beforeEach(() => {
        vi.mocked(useSkills).mockReset()
    })

    it('renders name, description and storage-mode chip', () => {
        renderCard({storageMode: 'freeform_notes'})
        expect(screen.getByText('Foo')).toBeInTheDocument()
        expect(screen.getByText('desc')).toBeInTheDocument()
        expect(screen.getByText('freeform notes')).toBeInTheDocument()
    })

    it('calls setHotPlug with the flipped value when the toggle is clicked', async () => {
        const {setHotPlug} = renderCard({isHotPlug: false})
        const user = userEvent.setup()

        const toggle = screen.getByRole('switch')
        await user.click(toggle)

        expect(setHotPlug).toHaveBeenCalledWith('vault-1', 'foo', true)
    })

    it('flips the other way when the skill is already hot-plugged', async () => {
        const {setHotPlug} = renderCard({isHotPlug: true})
        const user = userEvent.setup()

        await user.click(screen.getByRole('switch'))

        expect(setHotPlug).toHaveBeenCalledWith('vault-1', 'foo', false)
    })

    it('renders a Built-in badge and disables the toggle for system skills, hiding edit/delete', () => {
        renderCard({isSystem: true, isHotPlug: true})

        expect(screen.getByText('Built-in')).toBeInTheDocument()
        expect(screen.getByRole('switch')).toBeDisabled()
        expect(screen.queryByRole('button', {name: 'Edit skill'})).not.toBeInTheDocument()
        expect(screen.queryByRole('button', {name: 'Delete skill'})).not.toBeInTheDocument()
    })
})
