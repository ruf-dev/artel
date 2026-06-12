import cls from "./LineNumbers.module.css"

interface LineNumbersProps {
    lineCount: number
    scrollTop: number
}

export default function LineNumbers({ lineCount, scrollTop }: LineNumbersProps) {
    const numbers = Array.from({ length: lineCount }, (_, i) => i + 1)

    return (
        <div className={cls.LineNumbersContainer}>
            <div className={cls.Inner} style={{ transform: `translateY(-${scrollTop}px)` }}>
                {numbers.map(n => (
                    <div key={n} className={cls.LineNumber}>{n}</div>
                ))}
            </div>
        </div>
    )
}
