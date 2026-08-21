import {describe, expect, it, vi} from 'vitest'
import {render, screen} from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import type {Toast} from '@vervstack/chures'

import ToastCard from '@/segments/Toaster/components/ToastCard/ToastCard'
import {shouldDismissFromSwipe} from '@/segments/Toaster/components/ToastCard/shouldDismissFromSwipe'

describe('ToastCard', () => {
	const mockToast: Toast = {
		id: 1,
		title: 'Test Toast',
		description: 'This is a test description',
		level: 'Info',
	}

	it('renders title and description', () => {
		const onDismiss = vi.fn()
		render(<ToastCard toast={mockToast} onDismiss={onDismiss} />)

		expect(screen.getByText('Test Toast')).toBeInTheDocument()
		expect(screen.getByText('This is a test description')).toBeInTheDocument()
	})

	it('calls onDismiss with toast id when close button is clicked', async () => {
		const onDismiss = vi.fn()
		const user = userEvent.setup()
		render(<ToastCard toast={mockToast} onDismiss={onDismiss} />)

		const closeButton = screen.getByRole('button')
		await user.click(closeButton)

		expect(onDismiss).toHaveBeenCalledWith(1)
	})

	it('renders with different toast levels', () => {
		const onDismiss = vi.fn()
		const errorToast: Toast = {
			id: 2,
			title: 'Error Toast',
			description: 'Something went wrong',
			level: 'Error',
		}
		const {unmount} = render(<ToastCard toast={errorToast} onDismiss={onDismiss} />)

		expect(screen.getByText('Error Toast')).toBeInTheDocument()
		expect(screen.getByText('Something went wrong')).toBeInTheDocument()

		unmount()

		const warnToast: Toast = {
			id: 3,
			title: 'Warning Toast',
			description: 'Be careful',
			level: 'Warn',
		}
		render(<ToastCard toast={warnToast} onDismiss={onDismiss} />)

		expect(screen.getByText('Warning Toast')).toBeInTheDocument()
		expect(screen.getByText('Be careful')).toBeInTheDocument()
	})
})

describe('shouldDismissFromSwipe', () => {
	it('returns true when offset exceeds default threshold of 80', () => {
		expect(shouldDismissFromSwipe(90)).toBe(true)
		expect(shouldDismissFromSwipe(-90)).toBe(true)
	})

	it('returns false when offset is below default threshold', () => {
		expect(shouldDismissFromSwipe(30)).toBe(false)
		expect(shouldDismissFromSwipe(-30)).toBe(false)
	})

	it('returns true when offset equals threshold', () => {
		expect(shouldDismissFromSwipe(80)).toBe(true)
		expect(shouldDismissFromSwipe(-80)).toBe(true)
	})

	it('respects custom threshold', () => {
		expect(shouldDismissFromSwipe(50, 40)).toBe(true)
		expect(shouldDismissFromSwipe(30, 40)).toBe(false)
	})

	it('handles zero offset', () => {
		expect(shouldDismissFromSwipe(0)).toBe(false)
	})
})
