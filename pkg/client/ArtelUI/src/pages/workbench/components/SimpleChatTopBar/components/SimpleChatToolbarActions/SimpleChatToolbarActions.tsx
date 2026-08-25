import {Button} from "@vervstack/chures"
import {motion} from "framer-motion"

import HistoryIcon from "@/icons/common/HistoryIcon.tsx"
import PlusIcon from "@/icons/common/PlusIcon.tsx"
import ModelSwitcher from "@/pages/workbench/components/SimpleChat/components/ModelSwitcher/ModelSwitcher.tsx"
// eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
import cls from "@/pages/workbench/components/SimpleChatTopBar/components/SimpleChatToolbarActions/SimpleChatToolbarActions.module.css"

interface Props {
    models: string[]
    currentModel: string
    modelsLoading?: boolean
    onChangeModel: (model: string) => void
    onNewChat: () => void
    onToggleHistory: () => void
}

// >6 props — kept as one object instead of exploding into separate destructured
// bindings, per the ObjectPattern[properties.length>6] lint rule.
export default function SimpleChatToolbarActions(props: Props) {
    return (
        <div className={cls.SimpleChatToolbarActionsContainer}>
            <motion.div layout className={cls.ModelSwitcherWrapper}>
                <ModelSwitcher
                    models={props.models}
                    value={props.currentModel}
                    isLoading={props.modelsLoading}
                    onChange={props.onChangeModel}
                />
            </motion.div>
            <motion.div layout className={cls.NewChatButtonWrapper}>
                <Button
                    variant="secondary"
                    className={cls.NewChatButton}
                    onClick={props.onNewChat}
                    aria-label="New chat"
                    title="New chat"
                >
                    <PlusIcon/>
                </Button>
            </motion.div>
            <motion.div layout className={cls.HistoryButtonWrapper}>
                <Button
                    variant="secondary"
                    className={cls.HistoryButton}
                    onClick={props.onToggleHistory}
                    aria-label="View chat history"
                    title="View chat history"
                >
                    <HistoryIcon/>
                </Button>
            </motion.div>
        </div>
    )
}
