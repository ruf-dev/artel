import {beforeEach, describe, expect, it, vi} from 'vitest'

import {SkillsAPI} from '@/app/api/artel/skills.pb.ts'
import {useSkills} from '@/app/hooks/Skills.ts'

vi.mock('@/app/api/artel/skills.pb.ts', () => ({
    SkillsAPI: {
        ListSkills: vi.fn(),
        GetSkill: vi.fn(),
        CreateSkill: vi.fn(),
        UpdateSkill: vi.fn(),
        DeleteSkill: vi.fn(),
        SetSkillHotPlug: vi.fn(),
    },
}))

describe('useSkills', () => {
    beforeEach(() => {
        useSkills.setState({skills: [], loading: false})
        vi.mocked(SkillsAPI.ListSkills).mockReset()
        vi.mocked(SkillsAPI.CreateSkill).mockReset()
        vi.mocked(SkillsAPI.DeleteSkill).mockReset()
        vi.mocked(SkillsAPI.SetSkillHotPlug).mockReset()
    })

    it('fetch populates skills from the API response', async () => {
        vi.mocked(SkillsAPI.ListSkills).mockResolvedValue({
            skills: [{slug: 'foo', name: 'Foo', isHotPlug: false, isSystem: false}],
        })

        await useSkills.getState().fetch('vault-1')

        expect(SkillsAPI.ListSkills).toHaveBeenCalledWith({vaultId: 'vault-1'}, expect.anything())
        expect(useSkills.getState().skills).toHaveLength(1)
        expect(useSkills.getState().skills[0].slug).toBe('foo')
        expect(useSkills.getState().loading).toBe(false)
    })

    it('fetch toggles loading while in flight', async () => {
        let resolveList!: (v: { skills: never[] }) => void
        const responsePromise = new Promise<{ skills: never[] }>(resolve => { resolveList = resolve })
        vi.mocked(SkillsAPI.ListSkills).mockReturnValue(responsePromise)

        const pending = useSkills.getState().fetch('vault-1')
        expect(useSkills.getState().loading).toBe(true)

        resolveList({skills: []})
        await pending

        expect(useSkills.getState().loading).toBe(false)
    })

    it('create calls the API then refetches the list', async () => {
        vi.mocked(SkillsAPI.CreateSkill).mockResolvedValue({skill: {slug: 'new-skill', name: 'New skill'}})
        vi.mocked(SkillsAPI.ListSkills).mockResolvedValue({
            skills: [{slug: 'new-skill', name: 'New skill', isHotPlug: false, isSystem: false}],
        })

        const result = await useSkills.getState().create('vault-1', 'New skill', 'desc', 'none', 'body', false)

        expect(SkillsAPI.CreateSkill).toHaveBeenCalledWith(
            {
                vaultId: 'vault-1', name: 'New skill', description: 'desc', storageMode: 'none',
                body: 'body', hotPlug: false,
            },
            expect.anything(),
        )
        expect(SkillsAPI.ListSkills).toHaveBeenCalledWith({vaultId: 'vault-1'}, expect.anything())
        expect(result.slug).toBe('new-skill')
        expect(useSkills.getState().skills).toHaveLength(1)
    })

    it('setHotPlug calls the API then refetches the list', async () => {
        vi.mocked(SkillsAPI.SetSkillHotPlug).mockResolvedValue({
            skill: {slug: 'foo', name: 'Foo', isHotPlug: true},
        })
        vi.mocked(SkillsAPI.ListSkills).mockResolvedValue({
            skills: [{slug: 'foo', name: 'Foo', isHotPlug: true, isSystem: false}],
        })

        await useSkills.getState().setHotPlug('vault-1', 'foo', true)

        expect(SkillsAPI.SetSkillHotPlug).toHaveBeenCalledWith(
            {vaultId: 'vault-1', slug: 'foo', hotPlug: true}, expect.anything(),
        )
        expect(useSkills.getState().skills[0].isHotPlug).toBe(true)
    })

    it('remove calls the API then refetches the list', async () => {
        vi.mocked(SkillsAPI.DeleteSkill).mockResolvedValue({})
        vi.mocked(SkillsAPI.ListSkills).mockResolvedValue({skills: []})

        await useSkills.getState().remove('vault-1', 'foo')

        expect(SkillsAPI.DeleteSkill).toHaveBeenCalledWith({vaultId: 'vault-1', slug: 'foo'}, expect.anything())
        expect(useSkills.getState().skills).toHaveLength(0)
    })
})
