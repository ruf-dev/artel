import {createPortal} from "react-dom"
import {Button, Input} from "@vervstack/chures"

import cls from "@/components/TemplateInput/TemplateInput.module.css"
import {cn} from "@/app/utils/cn.ts"
import SourceGroups from "@/components/TemplateInput/components/SourceGroups.tsx"
import {useTemplateInputController} from "@/components/TemplateInput/processes/useTemplateInputController.ts"
import {TemplateSource} from "@/components/TemplateInput/processes/templateRefs.ts"

export type {TemplateSource}

interface Props {
    value: string
    onChange: (value: string) => void
    sources: TemplateSource[]
    placeholder?: string
}

export default function TemplateInput({value, onChange, sources, placeholder}: Props) {
    const ctl = useTemplateInputController(sources, value, onChange)

    return (
        <div className={cls.Wrap} ref={ctl.wrapRef}>
            <div className={cls.InputRow}>
                <Input
                    ref={ctl.inputRef}
                    className={cn(cls.Input, ctl.isInvalid && cls.InputInvalid)}
                    value={value}
                    setValue={ctl.handleChange}
                    onKeyDown={ctl.handleKeyDown}
                    onBlur={() => setTimeout(ctl.closeDropdown, 150)}
                    placeholder={placeholder}
                    spellCheck={false}
                    data-tooltip-id={ctl.tooltip ? "root-tooltip" : undefined}
                    data-tooltip-content={ctl.tooltip}
                />
                <Button
                    type="button"
                    variant="ghost"
                    className={cls.AffixBtn}
                    onClick={() => ctl.openDropdown(null)}
                >
                    {"{}"}
                </Button>
            </div>
            {ctl.open && ctl.dropdownRect && createPortal(
                <>
                    <div className={cls.Backdrop} onClick={ctl.closeDropdown} />
                    <div
                        className={cls.Dropdown}
                        style={{top: ctl.dropdownRect.top, left: ctl.dropdownRect.left, width: ctl.dropdownRect.width}}
                        onMouseDown={e => e.preventDefault()}
                    >
                        <SourceGroups
                            groups={ctl.filtered}
                            flatList={ctl.flatList}
                            activeIndex={ctl.activeIndex}
                            onSelect={ctl.insertRef}
                        />
                    </div>
                </>,
                document.body,
            )}
        </div>
    )
}
