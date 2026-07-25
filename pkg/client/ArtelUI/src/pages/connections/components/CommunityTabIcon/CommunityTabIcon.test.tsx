import {describe, expect, it} from 'vitest'
import {render} from '@testing-library/react'

import CommunityTabIcon from '@/pages/connections/components/CommunityTabIcon/CommunityTabIcon.tsx'

describe('CommunityTabIcon', () => {
    it('renders an svg without crashing', () => {
        const {container} = render(<CommunityTabIcon/>)
        expect(container.querySelector('svg')).toBeInTheDocument()
    })
})
