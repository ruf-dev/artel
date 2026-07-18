import {useMemo} from "react"
import CodeMirror from "@uiw/react-codemirror"
import {javascript} from "@codemirror/lang-javascript"
import {syntaxHighlighting} from "@codemirror/language"

import cls from "@/pages/tract-canvas/components/ScriptBody/components/ScriptCodeSection/ScriptCodeSection.module.css"
import {ScriptParam} from "@/processes/Tracts.ts"
import {generatedFooter, generatedHeader} from "@/pages/tract-canvas/processes/scriptSignature.ts"
import Section from "@/pages/tract-canvas/components/Section/Section.tsx"
import {
    paramCompletionSource,
    scriptEditorTheme,
    scriptHighlightStyle,
} from "@/pages/tract-canvas/components/ScriptBody/components/ScriptCodeSection/processes/scriptEditorExtensions.ts"

interface Props {
    name: string
    inputParams: ScriptParam[]
    outputParams: ScriptParam[]
    code: string
    onChange: (code: string) => void
}

/** ScriptCodeSection renders the immutable generated header (JSDoc, signature, zero-valued
 * output prelude) and footer (return statement, closing brace) around the editable body —
 * only the body between them is ever persisted as step.code. The body itself is a CodeMirror
 * editor (JS syntax highlighting + autocomplete hints for the declared param names), styled
 * via scriptEditorTheme to match the surrounding generated <pre> blocks. */
export default function ScriptCodeSection({name, inputParams, outputParams, code, onChange}: Props) {
    const extensions = useMemo(() => [
        javascript(),
        syntaxHighlighting(scriptHighlightStyle),
        scriptEditorTheme,
        paramCompletionSource(inputParams, outputParams),
    ], [inputParams, outputParams])

    return (
        <Section title="Code">
            <div className={cls.ScriptCodeSectionContainer}>
                <pre className={cls.Generated}>{generatedHeader(name, inputParams, outputParams)}</pre>
                <CodeMirror
                    className={cls.CodeEditor}
                    value={code}
                    onChange={onChange}
                    extensions={extensions}
                    theme="none"
                    minHeight="9rem"
                    basicSetup={{lineNumbers: false, foldGutter: false}}
                />
                <pre className={cls.Generated}>{generatedFooter(outputParams)}</pre>
            </div>
        </Section>
    )
}
