import cls from "@/pages/tract-canvas/components/ScriptBody/components/ScriptCodeSection/ScriptCodeSection.module.css"
import {ScriptParam} from "@/processes/Tracts.ts"
import {generatedFooter, generatedHeader} from "@/pages/tract-canvas/processes/scriptSignature.ts"
import Section from "@/pages/tract-canvas/components/Section/Section.tsx"

interface Props {
    name: string
    inputParams: ScriptParam[]
    outputParams: ScriptParam[]
    code: string
    onChange: (code: string) => void
}

/** ScriptCodeSection renders the immutable generated header (JSDoc, signature, zero-valued
 * output prelude) and footer (return statement, closing brace) around the editable body —
 * only the body between them is ever persisted as step.code. */
export default function ScriptCodeSection({name, inputParams, outputParams, code, onChange}: Props) {
    return (
        <Section title="Code">
            <div className={cls.ScriptCodeSectionContainer}>
                <pre className={cls.Generated}>{generatedHeader(name, inputParams, outputParams)}</pre>
                <textarea
                    className={cls.CodeEditor}
                    value={code}
                    onChange={e => onChange(e.target.value)}
                    rows={8}
                    spellCheck={false}
                />
                <pre className={cls.Generated}>{generatedFooter(outputParams)}</pre>
            </div>
        </Section>
    )
}
