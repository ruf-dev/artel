import {Button} from "@vervstack/chures"

interface Props {
    onClick: () => void
}

export default function NewTractButton({onClick}: Props) {
    return (
        <Button variant="primary" onClick={onClick}>
            + New tract
        </Button>
    )
}
