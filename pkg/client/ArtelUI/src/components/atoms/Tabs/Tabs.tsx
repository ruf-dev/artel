// TODO: chures has no tab primitive yet, drop this wrapper once it does
import {Button} from "@vervstack/chures"

import cls from "@/components/atoms/Tabs/Tabs.module.css"
import {cn} from "@/app/utils/cn.ts"

interface TabsProps {
    tabs: {key: string; label: string; icon?: React.ReactNode}[]
    active: string
    onChange: (key: string) => void
}

export default function Tabs({tabs, active, onChange}: TabsProps) {
    return (
        <div className={cls.TabsContainer}>
            {tabs.map(({key, label, icon}) => (
                <Button
                    key={key}
                    variant="unstyled"
                    className={cn(cls.Tab, active === key && cls.TabActive)}
                    onClick={() => onChange(key)}
                >
                    {icon && <span className={cls.TabIcon}>{icon}</span>}
                    {label}
                </Button>
            ))}
        </div>
    )
}
