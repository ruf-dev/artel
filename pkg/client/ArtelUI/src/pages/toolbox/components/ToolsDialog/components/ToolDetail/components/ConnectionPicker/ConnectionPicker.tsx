import {useNavigate} from "react-router-dom"
import {Button} from "@vervstack/chures"

import {ExternalConnectionInfo} from "@/app/api/artel/external_connections.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"
import cls from "@/pages/toolbox/components/ToolsDialog/components/ToolDetail/components/ConnectionPicker/ConnectionPicker.module.css"

export default function ConnectionPicker({connections, selectedId, onSelect}: {
    connections: ExternalConnectionInfo[]
    selectedId: string
    onSelect: (id: string) => void
}) {
    const navigate = useNavigate()

    if (connections.length === 0) {
        return (
            <div className={cls.ConnectionPickerContainer}>
                <span className={cls.SectionTitle}>Connection</span>
                <Button
                    variant="ghost"
                    onClick={() => navigate(Path.ConnectionsPage)}>
                    No connections — add one
                </Button>
            </div>
        )
    }

    return (
        <div className={cls.ConnectionPickerContainer}>
            <span className={cls.SectionTitle}>Connection</span>
            {connections.map(c => (
                <SelectOption
                    key={c.id}
                    label={connectionLabel(c)}
                    selected={c.id === selectedId}
                    onSelect={() => onSelect(c.id ?? "")}
                />
            ))}
        </div>
    )
}
