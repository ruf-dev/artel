import {useEffect} from "react"
import {useNavigate} from "react-router-dom"
import {Button} from "@vervstack/chures"

import {ExternalConnectionInfo} from "@/app/api/artel/external_connections.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"
import cls from "@/dialogs/InstantiateTemplateDialog/components/ConnectionSection/ConnectionSection.module.css"

interface Props {
    mcp: string
    connections: ExternalConnectionInfo[]
    selectedId: string
    onSelect: (id: string) => void
}

export default function ConnectionSection({mcp, connections, selectedId, onSelect}: Props) {
    const {CloseDialog} = useDialog()
    const navigate = useNavigate()

    useEffect(() => {
        if (!selectedId && connections.length === 1) {
            onSelect(connections[0].id ?? "")
        }
    }, [connections, selectedId, onSelect])

    return (
        <div className={cls.ConnectionSectionContainer}>
            <span className={cls.SegmentHeading}>{mcp}</span>
            {connections.length === 0 ? (
                <p className={cls.Empty}>
                    No {mcp} connections.{" "}
                    <Button
                        variant="ghost"
                        className={cls.LinkBtn}
                        onClick={() => {
                            CloseDialog()
                            navigate(Path.ConnectionsPage)
                        }}
                    >
                        Set one up
                    </Button>
                </p>
            ) : (
                connections.map(c => (
                    <SelectOption
                        key={c.id}
                        label={connectionLabel(c)}
                        selected={c.id === selectedId}
                        onSelect={() => onSelect(c.id ?? "")}
                    />
                ))
            )}
        </div>
    )
}
