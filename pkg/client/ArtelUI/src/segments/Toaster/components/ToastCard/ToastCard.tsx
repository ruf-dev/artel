import {motion} from 'framer-motion'
import {ModalClose} from '@vervstack/chures'
import type {Toast} from '@vervstack/chures'

import {cn} from '@/app/utils/cn'
import cls from '@/segments/Toaster/components/ToastCard/ToastCard.module.css'
import {shouldDismissFromSwipe} from '@/segments/Toaster/components/ToastCard/shouldDismissFromSwipe'

interface ToastCardProps {
	toast: Toast
	onDismiss: (id: number) => void
}

function handleDragEnd(
	toast: Toast,
	onDismiss: (id: number) => void,
	_: unknown,
	info: {offset: {x: number}},
): void {
	if (shouldDismissFromSwipe(info.offset.x)) {
		onDismiss(toast.id!)
	}
}

export default function ToastCard({toast, onDismiss}: ToastCardProps) {
	const toastLevel = toast.level?.toLowerCase() as 'error' | 'warn' | 'info' | undefined
	const levelClass = toastLevel ? cls[`Level${toastLevel.charAt(0).toUpperCase() + toastLevel.slice(1)}`] : undefined

	return (
		<motion.div
			key={toast.id}
			className={cls.ToastCardContainer}
			drag="x"
			dragSnapToOrigin
			onDragEnd={(e, info) => handleDragEnd(toast, onDismiss, e, info)}
			initial={{opacity: 0, x: -20}}
			animate={{opacity: 1, x: 0}}
			exit={{opacity: 0, x: 20}}
			transition={{duration: 0.2}}
		>
			<div className={cn(cls.ContentWrapper, levelClass)}>
				<div className={cls.TextContent}>
					<div className={cls.Title}>{toast.title}</div>
					{toast.description && (
						<div className={cls.Description}>{toast.description}</div>
					)}
				</div>
				<ModalClose
					className={cls.CloseButton}
					onClick={() => onDismiss(toast.id!)}
				/>
			</div>
		</motion.div>
	)
}
