import {forwardRef} from "react"

import {cn} from "@/app/utils/cn"

import cls from "./Input.module.css"

type InputProps = React.InputHTMLAttributes<HTMLInputElement>

const Input = forwardRef<HTMLInputElement, InputProps>(({className, ...props}, ref) => (
    <input ref={ref} className={cn(cls.Input, className)} {...props}/>
))

Input.displayName = "Input"

export default Input
