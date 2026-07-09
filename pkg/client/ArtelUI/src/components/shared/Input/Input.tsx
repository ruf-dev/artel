import {forwardRef} from "react"

import {cn} from "@/app/utils/cn"
import cls from "@/components/shared/Input/Input.module.css"

type InputProps = React.InputHTMLAttributes<HTMLInputElement>

// TODO: predates the components/atoms/ convention (see pkg/client/ArtelUI/CLAUDE.md
// "Known debt") — move under components/atoms/Input once callers are migrated.
const Input = forwardRef<HTMLInputElement, InputProps>(({className, ...props}, ref) => (
    <input ref={ref} className={cn(cls.Input, className)} {...props}/>
))

Input.displayName = "Input"

export default Input
