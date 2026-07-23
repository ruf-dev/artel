// Minimal ambient typings for `bun:test`, scoped locally instead of installing
// @types/bun: that package globally augments `typeof fetch` (adding a
// required `preconnect` member) and breaks type-checking on every other file
// that assigns to `fetch` (e.g. AuthFetchInterceptor.ts). Only the bits used
// by ServerStatus.test.ts are declared here.
declare module "bun:test" {
    export function describe(name: string, fn: () => void): void
    export function it(name: string, fn: () => void | Promise<void>): void
    export function beforeEach(fn: () => void | Promise<void>): void
    export function afterEach(fn: () => void | Promise<void>): void
    export function mock<T extends (...args: never[]) => unknown>(fn: T): T
    export function expect<T>(actual: T): {toBe(expected: T): void}
}
