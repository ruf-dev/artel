// Pure text generation for a script step's read-only header/footer around its editable
// code body (ScriptBody's Code section). This must produce EXACTLY what the backend's
// script.javascript engine wraps the stored body with (internal/service/v1/tract/script/
// javascript.go's buildSource) — signature, JSDoc, output prelude, and return statement are
// never persisted, only ever derived from inputParams/outputParams, here and server-side.

import {ScriptParam} from "@/processes/Tracts.ts"

const JSDOC_TYPE: Record<string, string> = {
    string: "string",
    integer: "number",
    number: "number",
    boolean: "boolean",
    array: "Array",
    object: "object",
}

function jsDocType(param: ScriptParam): string {
    return JSDOC_TYPE[param.type.type] ?? "any"
}

function zeroValueLiteral(type: string): string {
    switch (type) {
        case "string":
            return "\"\""
        case "integer":
        case "number":
            return "0"
        case "boolean":
            return "false"
        case "array":
            return "[]"
        case "object":
            return "{}"
        default:
            return "null"
    }
}

/** generatedHeader is the read-only block shown above the editable code textarea: JSDoc,
 * the function signature, and the zero-valued declaration for every output param. */
export function generatedHeader(
    name: string, inputParams: ScriptParam[], outputParams: ScriptParam[],
): string {
    const lines: string[] = ["/**"]

    for (const p of inputParams) {
        const desc = p.type.description ? ` - ${p.type.description}` : ""
        lines.push(` * @param {${jsDocType(p)}} ${p.name}${desc}`)
    }

    if (outputParams.length > 0) {
        const shape = outputParams.map(p => `${p.name}: ${jsDocType(p)}`).join(", ")
        lines.push(` * @returns {{${shape}}}`)
    }

    lines.push(" */")
    lines.push(`function ${name}(${inputParams.map(p => p.name).join(", ")}) {`)

    if (outputParams.length > 0) {
        const decls = outputParams.map(p => `${p.name} = ${zeroValueLiteral(p.type.type)}`).join(", ")
        lines.push(`  let ${decls};`)
    }

    return lines.join("\n")
}

/** generatedFooter is the read-only block shown below the editable code textarea: the
 * return statement listing every declared output param, and the closing brace. */
export function generatedFooter(outputParams: ScriptParam[]): string {
    const shape = outputParams.map(p => `${p.name}: ${p.name}`).join(", ")
    return `  return {${shape}};\n}`
}
