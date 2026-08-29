import {Button} from "@vervstack/chures"

import {getItemName} from "@/components/FileTree/fileTree.ts"
import FileIcon from "@/components/FileTree/icons/FileIcon.tsx"
import cls from "@/pages/workbench/components/Chat/components/AttachmentChip/AttachmentChip.module.css"

interface Props {
    path: string
    onClick: (path: string) => void
}

// Read-only counterpart to ComposerCtxChip: shows an attached vault file as a pill chip
// with file icon and basename. Used in UserMessageBubble to display attachments after a
// message is sent (vs ComposerCtxChip which shows pre-send/editable attachments in the
// composer). Clicking the chip opens the attached note for viewing.
export default function AttachmentChip({path, onClick}: Props) {
    return (
        <Button
            variant="unstyled"
            className={cls.AttachmentChipContainer}
            onClick={() => onClick(path)}
            title={path}
            aria-label={`View attached file: ${path}`}
        >
            <FileIcon/>
            <span className={cls.Label}>{getItemName({path})}</span>
        </Button>
    )
}
