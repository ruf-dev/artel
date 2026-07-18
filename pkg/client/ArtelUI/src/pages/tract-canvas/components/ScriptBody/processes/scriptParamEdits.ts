// Pure step-transition helpers for ScriptBody's Inputs/Outputs editors — split out of the
// component purely to keep it under the project's max-lines-per-function limit; each
// function takes the current step and returns the next one, ready to hand to replaceStep's
// updater callback.

import {ScriptParam, TractStep} from "@/processes/Tracts.ts"

function uniqueParamName(base: string, taken: Set<string>): string {
    if (!taken.has(base)) return base
    let i = 2
    while (taken.has(`${base}${i}`)) i++
    return `${base}${i}`
}

export function addInputParam(s: TractStep): TractStep {
    const params = s.inputParams ?? []
    const name = uniqueParamName("param", new Set(params.map(p => p.name)))
    const next: ScriptParam = {name, type: {type: "string"}}
    return {...s, inputParams: [...params, next], params: {...s.params, [name]: ""}}
}

/** updateInputParam keeps step.params (the template-expr bindings, keyed by param name) in
 * sync when a param is renamed — moving its binding value from the old key to the new one
 * rather than leaving an orphaned entry under the old name. */
export function updateInputParam(s: TractStep, i: number, patch: Partial<ScriptParam>): TractStep {
    const params = [...(s.inputParams ?? [])]
    const prevName = params[i].name
    params[i] = {...params[i], ...patch}

    if (patch.name === undefined || patch.name === prevName) return {...s, inputParams: params}

    const bindings = {...s.params}
    bindings[patch.name] = bindings[prevName] ?? ""
    delete bindings[prevName]
    return {...s, inputParams: params, params: bindings}
}

export function removeInputParam(s: TractStep, i: number): TractStep {
    const params = [...(s.inputParams ?? [])]
    const [removed] = params.splice(i, 1)
    const bindings = {...s.params}
    if (removed) delete bindings[removed.name]
    return {...s, inputParams: params, params: bindings}
}

export function addOutputParam(s: TractStep): TractStep {
    const params = s.outputParams ?? []
    const name = uniqueParamName("result", new Set(params.map(p => p.name)))
    const next: ScriptParam = {name, type: {type: "string"}}
    return {...s, outputParams: [...params, next]}
}

export function updateOutputParam(s: TractStep, i: number, patch: Partial<ScriptParam>): TractStep {
    const params = [...(s.outputParams ?? [])]
    params[i] = {...params[i], ...patch}
    return {...s, outputParams: params}
}

export function removeOutputParam(s: TractStep, i: number): TractStep {
    const params = [...(s.outputParams ?? [])]
    params.splice(i, 1)
    return {...s, outputParams: params}
}
