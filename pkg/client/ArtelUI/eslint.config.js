import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import reactPlugin from 'eslint-plugin-react'
import importX from 'eslint-plugin-import-x'
import {createTypeScriptImportResolver} from 'eslint-import-resolver-typescript'
import tseslint from 'typescript-eslint'
import {globalIgnores} from 'eslint/config'

export default tseslint.config([
    globalIgnores(['dist', '**/*.pb.ts']),
    {
        files: ['**/*.{ts,tsx}'],
        extends: [
            js.configs.recommended,
            tseslint.configs.recommended,
            reactHooks.configs['recommended-latest'],
            reactRefresh.configs.vite,
        ],
        plugins: {
            react: reactPlugin,
            'import-x': importX,
        },
        languageOptions: {
            ecmaVersion: 2020,
            globals: globals.browser,
        },
        settings: {
            'import-x/resolver-next': [
                createTypeScriptImportResolver(),
            ],
        },
        rules: {
            'react-hooks/exhaustive-deps': 'off',
            'import-x/no-cycle': 'error',
            'no-empty': ['error', {allowEmptyCatch: true}],
            'react/jsx-max-depth': ['error', {max: 3}],
            'no-restricted-syntax': [
                'warn',
                {
                    selector: 'JSXOpeningElement[name.name="div"]:not(:has(JSXAttribute[name.name="className"]))',
                    message: '<div> must have a className (use a CSS module class).',
                },
                {
                    selector: 'Property[key.name="zIndex"]',
                    message: 'Never use z-index — rely on DOM order or a portal to document.body instead (see pkg/client/ArtelUI/CLAUDE.md).',
                },
                {
                    selector: 'JSXAttribute[name.name="className"] > JSXExpressionContainer > TemplateLiteral',
                    message: 'Do not build className with a template literal — use cn() from @/app/utils/cn instead.',
                },
            ],
        },
    },
])
