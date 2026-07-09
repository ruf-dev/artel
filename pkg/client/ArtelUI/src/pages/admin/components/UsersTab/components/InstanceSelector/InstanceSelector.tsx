import {GetCouchInstanceResponse} from "@/app/api/artel/couch_instances.pb.ts"
import cls from "@/pages/admin/components/UsersTab/components/InstanceSelector/InstanceSelector.module.css"

interface InstanceSelectorProps {
    instances: GetCouchInstanceResponse[]
    value: string
    onChange: (id: string) => void
}

export default function InstanceSelector({instances, value, onChange}: InstanceSelectorProps) {
    return (
        <div className={cls.InstanceSelectorContainer}>
            <span className={cls.FieldLabel}>CouchDB instance</span>
            <select
                className={cls.Select}
                value={value}
                onChange={e => onChange(e.target.value)}
            >
                <option value="">— select an instance —</option>
                {instances.map(inst => (
                    <option key={inst.id} value={inst.id ?? ""}>{inst.url}</option>
                ))}
            </select>
        </div>
    )
}
