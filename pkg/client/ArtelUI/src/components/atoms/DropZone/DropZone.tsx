// TODO: chures has no drag-and-drop file dropzone yet, drop this wrapper once it does
import { useState } from "react"

import { cn } from "@/app/utils/cn.ts"
import cls from "@/components/atoms/DropZone/DropZone.module.css"

interface Props {
    accept?: string
    selectedFileName?: string | null
    onFileSelected: (file: File) => void
}

export default function DropZone({ accept, selectedFileName, onFileSelected }: Props) {
    const [dragOver, setDragOver] = useState(false)

    function handleDrop(e: React.DragEvent<HTMLLabelElement>) {
        e.preventDefault()
        setDragOver(false)
        const file = e.dataTransfer.files[0]
        if (file) onFileSelected(file)
    }

    function handleInputChange(e: React.ChangeEvent<HTMLInputElement>) {
        const file = e.target.files?.[0]
        if (file) onFileSelected(file)
    }

    return (
        <label
            className={cn(cls.DropZoneContainer, dragOver && cls.DropZoneContainerActive)}
            onDragOver={e => { e.preventDefault(); setDragOver(true) }}
            onDragLeave={() => setDragOver(false)}
            onDrop={handleDrop}
        >
            <input type="file" accept={accept} className={cls.HiddenInput} onChange={handleInputChange} />
            <span className={cls.DropZoneText}>
                {selectedFileName ?? "Drop a .zip file here, or click to browse"}
            </span>
        </label>
    )
}
