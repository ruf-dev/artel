import Button from "@/components/shared/Button/Button.tsx"
import {PlayIcon} from "@/pages/tract-canvas/components/TractIcons/TractIcons.tsx"

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
