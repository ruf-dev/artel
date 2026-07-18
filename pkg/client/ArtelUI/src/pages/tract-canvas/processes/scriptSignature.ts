// Pure text generation for a script step's read-only header/footer around its editable
// code body (ScriptBody's Code section). This must produce EXACTLY what the backend's
// script.javascript engine wraps the stored body with (internal/service/v1/tract/script/
// javascript.go's buildSource) — signature, JSDoc, output prelude, and return statement are
// never persisted, only ever derived from inputParams/outputParams, here and server-side.

import {ScriptParam, SchemaProperty} from "@/processes/Tracts.ts"

const JSDOC_PRIMITIVE: Record<string, string> = {
    string: "string",
    integer: "number",
    number: "number",
    boolean: "boolean",
    array: "Array",
    object: "object",
}

function jsDocPrimitive(type: string): string {
    return JSDOC_PRIMITIVE[type] ?? "any"
}

function toPascalCase(name: string): string {
    return name
        .replace(/[^a-zA-Z0-9]+(.)/g, (_, c: string) => c.toUpperCase())
        .replace(/^./, c => c.toUpperCase())
}

/** objectTypedefLines renders a `@typedef {Object} Name` block documenting a flat object
 * schema's properties (one level deep — a property's own type is always primitive here,
 * never a nested array/object, matching what the param-shape editor lets users declare). */
function objectTypedefLines(typedefName: string, schema: SchemaProperty): string[] {
    const lines = ["/**", ` * @typedef {Object} ${typedefName}`]

    for (const [propName, propSchema] of Object.entries(schema.properties ?? {})) {
        const desc = propSchema.description ? ` - ${propSchema.description}` : ""
        lines.push(` * @property {${jsDocPrimitive(propSchema.type)}} ${propName}${desc}`)
    }

    lines.push(" */")
    return lines
}

/** resolvedJsDocType maps a param to the type expression used in its `@param`/`@returns`
 * entry, plus an optional `@typedef` block (printed above the main JSDoc) when the param's
 * type — or an array param's item type — is an object with a declared property shape.
 * Falls back to the bare "Array"/"object" JSDoc type when no shape has been declared, so
 * params created before the shape editor existed still render sensibly. */
function resolvedJsDocType(param: ScriptParam): { expr: string, typedef?: string[] } {
    const type = param.type

    if (type.type === "array" && type.items) {
        const itemType = type.items

        if (itemType.type === "object" && itemType.properties) {
            const typedefName = `${toPascalCase(param.name)}Item`
            return {expr: `${typedefName}[]`, typedef: objectTypedefLines(typedefName, itemType)}
        }

        return {expr: `${jsDocPrimitive(itemType.type)}[]`}
    }

    if (type.type === "object" && type.properties) {
        const typedefName = `${toPascalCase(param.name)}Shape`
        return {expr: typedefName, typedef: objectTypedefLines(typedefName, type)}
    }

    return {expr: jsDocPrimitive(type.type)}
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
    const typedefBlocks: string[][] = []
    const paramLines: string[] = []

    for (const p of inputParams) {
        const {expr, typedef} = resolvedJsDocType(p)
        if (typedef) typedefBlocks.push(typedef)

        const desc = p.type.description ? ` - ${p.type.description}` : ""
        paramLines.push(` * @param {${expr}} ${p.name}${desc}`)
    }

    const outputExprs = outputParams.map(p => {
        const {expr, typedef} = resolvedJsDocType(p)
        if (typedef) typedefBlocks.push(typedef)
        return `${p.name}: ${expr}`
    })

    const lines: string[] = typedefBlocks.flatMap(block => [...block, ""])
    lines.push("/**", ...paramLines)

    if (outputExprs.length > 0) {
        lines.push(` * @returns {{${outputExprs.join(", ")}}}`)
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
