import {useNavigate} from "react-router-dom"
import {Button} from "@vervstack/chures"

import {MomCandidate} from "@/app/api/artel/mcp_keys.pb.ts"
import {Path} from "@/app/routing/Router.tsx"
import {useDialog} from "@/app/hooks/Dialog"
import SelectOption from "@/components/SelectOption/SelectOption.tsx"
import {connectionLabel} from "@/components/ConnectorChip/connectionLabel.ts"
import cls from "@/dialogs/ManageKeyDialog/components/ConnectionOptionList/ConnectionOptionList.module.css"

interface ConnectionOptionListProps {
    available: NonNullable<MomCandidate["connections"]>
    selectedId: string
    onSelect: (id: string) => void
}

export default function ConnectionOptionList({available, selectedId, onSelect}: ConnectionOptionListProps) {
    const {CloseDialog} = useDialog()
    const navigate = useNavigate()
    return (
        <div className={cls.OptionList}>
            {available.map(c => (
                <SelectOption
                    key={c.id}
                    label={connectionLabel(c)}
                    selected={selectedId === c.id}
                    onSelect={() => onSelect(c.id ?? "")}
                />
            ))}
            {available.length === 0 && (
                <p className={cls.Empty}>
                    No connections yet.{" "}
                    <Button variant="ghost" className={cls.LinkBtn}
                            onClick={() => {
                                CloseDialog();
                                navigate(Path.ConnectionsPage)
                            }}>Set one up
                    </Button>
                </p>
            )}
        </div>
    )
}
