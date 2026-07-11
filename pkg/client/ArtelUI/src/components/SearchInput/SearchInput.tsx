import { Input } from "@vervstack/chures"

import { cn } from "@/app/utils/cn"
import cls from "@/components/SearchInput/SearchInput.module.css"

interface Props {
    value: string
    setValue: (v: string) => void
    placeholder?: string
    autoFocus?: boolean
    className?: string
}

export default function SearchInput({ value, setValue, placeholder, autoFocus, className }: Props) {
    return (
        <div className={cn(cls.SearchInputContainer, className)}>
            <Input
                inputClassName={cls.Field}
                type="text"
                placeholder={placeholder}
                value={value}
                setValue={setValue}
                autoFocus={autoFocus}
            />
        </div>
    )
}
