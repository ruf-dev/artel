import {forwardRef} from "react"
import cls from "./Input.module.css"
import {cn} from "@/app/utils/cn"

type InputProps = React.InputHTMLAttributes<HTMLInputElement>

const Input = forwardRef<HTMLInputElement, InputProps>(({className, ...props}, ref) => (
    <input ref={ref} className={cn(cls.Input, className)} {...props}/>
))

Input.displayName = "Input"

export default Input
