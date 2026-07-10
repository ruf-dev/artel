// TODO: chures' Input puts the caller-supplied className on its own outer wrapper
// div, never on the inner <input> it renders — text-level overrides (font-size,
// font-family, color, letter-spacing) declared against that className never reach
// the rendered text. This wrapper forces the inner input to inherit those
// properties from the wrapper div instead. Drop it if chures ever exposes an
// inputClassName/inputStyle prop for the inner element.
import { Input as ChuresInput } from "@vervstack/chures"
import type { InputProps } from "@vervstack/chures"

import { cn } from "@/app/utils/cn"
import cls from "@/components/atoms/Input/Input.module.css"

export default function Input({ className, ...rest }: InputProps) {
    return <ChuresInput className={cn(cls.InputWrapper, className)} {...rest} />
}
