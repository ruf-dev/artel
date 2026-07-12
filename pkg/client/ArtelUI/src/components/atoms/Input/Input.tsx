// The project's single themed drop-in for chures' Input — the one place to edit
// the app's input look (border/bg/color/focus) instead of duplicating it as a
// local override per call site. chures' `className` prop only ever lands on its
// outer wrapper div, never on the inner <input>, so this styles the input via a
// `.InputWrapper input` descendant selector rather than `inputClassName`.
import { Input as ChuresInput } from "@vervstack/chures"
import type { InputProps } from "@vervstack/chures"

import { cn } from "@/app/utils/cn"
import cls from "@/components/atoms/Input/Input.module.css"

export default function Input({ className, ...rest }: InputProps) {
    return <ChuresInput className={cn(cls.InputWrapper, className)} {...rest} />
}
