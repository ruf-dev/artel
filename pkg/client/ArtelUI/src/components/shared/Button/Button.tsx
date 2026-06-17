import cls from "@/components/shared/Button/Button.module.css"

export type ButtonVariant = "primary" | "secondary" | "danger" | "ghost" | "iconDanger"

const variantClass: Record<ButtonVariant, string> = {
    primary: cls.Primary,
    secondary: cls.Secondary,
    danger: cls.Danger,
    ghost: cls.Ghost,
    iconDanger: cls.IconDanger,
}

interface Props extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: ButtonVariant
}

export default function Button({variant, className, children, ...rest}: Props) {
    const classes = [cls.Btn, variant ? variantClass[variant] : undefined, className].filter(Boolean).join(" ")
    return (
        <button type="button" {...rest} className={classes}>
            {children}
        </button>
    )
}
