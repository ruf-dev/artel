import {Button} from "@vervstack/chures"

import {PlayIcon} from "@/pages/tract-canvas/icons/PlayIcon/PlayIcon.tsx"

interface Props {
    running: boolean
    onClick: () => void
}

export default function RunButton({running, onClick}: Props) {
    return (
        <Button variant="primary" onClick={onClick} disabled={running}>
            <PlayIcon/> {running ? "Running…" : "Run"}
        </Button>
    )
}
