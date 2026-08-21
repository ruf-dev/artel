import {useToaster} from '@vervstack/chures'
import {AnimatePresence} from 'framer-motion'

import cls from '@/segments/Toaster/Toaster.module.css'
import ToastCard from '@/segments/Toaster/components/ToastCard/ToastCard'

export default function Toaster() {
	const {toasts, dismiss} = useToaster()

	return (
		<div className={cls.ToasterContainer}>
			<AnimatePresence mode="popLayout">
				{toasts.map((toast) => (
					<ToastCard
						key={toast.id}
						toast={toast}
						onDismiss={dismiss}
					/>
				))}
			</AnimatePresence>
		</div>
	)
}
