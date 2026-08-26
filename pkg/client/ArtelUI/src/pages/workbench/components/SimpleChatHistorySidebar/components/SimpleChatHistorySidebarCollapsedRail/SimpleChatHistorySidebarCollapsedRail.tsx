import {Button, ChevronDownIcon} from "@vervstack/chures"

import PlusIcon from "@/icons/common/PlusIcon.tsx"
import cls from
    // eslint-disable-next-line max-len -- import path too long to wrap under 120 chars
    "@/pages/workbench/components/SimpleChatHistorySidebar/components/SimpleChatHistorySidebarCollapsedRail/SimpleChatHistorySidebarCollapsedRail.module.css"

interface Props {
    onExpand: () => void
    onNewChat: () => void
}

export default function SimpleChatHistorySidebarCollapsedRail(props: Props) {
    return (
        <div className={cls.SimpleChatHistorySidebarCollapsedRailContainer}>
            <Button
                variant="secondary"
                onClick={props.onExpand}
                aria-label="Expand sidebar"
                title="Expand sidebar"
            >
                <ChevronDownIcon className={cls.ExpandIcon}/>
            </Button>
            <Button
                variant="secondary"
                onClick={props.onNewChat}
                aria-label="New chat"
                title="New chat"
            >
                <PlusIcon/>
            </Button>
        </div>
    )
}
