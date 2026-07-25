import {fileURLToPath, URL} from 'node:url'
import {defineConfig} from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
        },
    },
    test: {
        environment: 'jsdom',
        setupFiles: ['./src/test/setup.ts'],
        css: true,
        // A couple of pre-existing test files in this repo are written for Bun's
        // built-in test runner (`bun:test`, run via `bun test`), not Vitest — exclude
        // them here. This Vitest setup is scoped to the community-connectors feature's
        // new components; it isn't a retrofit of test coverage across the rest of the app.
        exclude: [
            '**/node_modules/**',
            '**/dist/**',
            'src/app/hooks/ServerStatus.test.ts',
            'src/dialogs/InstantiateTemplateDialog/processes/requiredConnections.test.ts',
        ],
    },
})
