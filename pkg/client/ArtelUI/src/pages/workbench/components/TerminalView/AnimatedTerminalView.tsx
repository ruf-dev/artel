import {motion} from "framer-motion"

import cls from "@/pages/workbench/components/TerminalView/AnimatedTerminalView.module.css"
import TerminalView from "@/pages/workbench/components/TerminalView/TerminalView.tsx"

interface Props {
    vaultId: string
    tabs: {id: string; name: string; active: boolean}[]
    onSelectTab: (id: string) => void
    onCreateTab: () => void
    onCloseTab: (id: string) => void
    pendingTerminalAuthLink?: string
}

// Wraps TerminalView with a fade transition so it isn't a jarring instant pop-in
// against the Chat panel's own enter/exit animation when toggling views — see
// WorkbenchPage.tsx's AnimatePresence around the chat/terminal swap.
export default function AnimatedTerminalView(props: Props) {
    return (
        <motion.div
            className={cls.AnimatedTerminalViewContainer}
            initial={{opacity: 0}}
            animate={{opacity: 1}}
            exit={{opacity: 0}}
            transition={{duration: 0.18}}
        >
            <TerminalView {...props}/>
        </motion.div>
    )
}
