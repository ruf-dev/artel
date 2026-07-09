import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import reactPlugin from 'eslint-plugin-react'
import importX from 'eslint-plugin-import-x'
import {createTypeScriptImportResolver} from 'eslint-import-resolver-typescript'
import tseslint from 'typescript-eslint'
import {globalIgnores} from 'eslint/config'

const localPlugin = {
    rules: {
        'no-relative-imports': {
            meta: {type: 'suggestion'},
            create(context) {
                return {
                    ImportDeclaration(node) {
                        const src = node.source.value
                        if (src.startsWith('./') || src.startsWith('../')) {
                            context.report({
                                node,
                                message: "Use '@/' path alias instead of relative imports.",
                            })
                        }
                    },
                }
            },
        },
    },
}

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
            local: localPlugin,
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
            'local/no-relative-imports': 'error',
            'func-style': ['error', 'declaration', {allowArrowFunctions: false}],
            'react/no-multi-comp': ['error', {ignoreStateless: false}],
            'no-console': ['error', {allow: ['warn', 'error']}],
            'react/forbid-component-props': ['error', {forbid: ['style']}],
            'max-lines': ['error', {max: 300, skipBlankLines: true, skipComments: true}],
            'max-lines-per-function': ['error', {max: 100, skipBlankLines: true, skipComments: true, IIFEs: true}],
            'no-restricted-syntax': [
                'error',
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
                {
                    selector: 'JSXOpeningElement[name.name="button"]',
                    message: 'Never use a raw <button> — use Button from @vervstack/chures, or a components/atoms/ wrapper if chures doesn\'t cover the need.',
                },
                {
                    selector: 'JSXOpeningElement[name.name="input"]',
                    message: 'Never use a raw <input> — use Input from @vervstack/chures, or a components/atoms/ wrapper if chures doesn\'t cover the need.',
                },
                {
                    selector: 'CallExpression[callee.object.name="window"][callee.property.name=/^(alert|confirm)$/]',
                    message: 'Never use window.alert/window.confirm — use useBakeError() for errors or chures ConfirmDialog for confirmations (see pkg/client/ArtelUI/CLAUDE.md).',
                },
                {
                    selector: 'ObjectPattern[properties.length>6]',
                    message: 'Destructuring more than 6 properties — keep the whole object as one variable (e.g. `const context = useX(...)`) instead of exploding it into separate bindings/props.',
                },
            ],
            'import-x/order': ['error', {
                groups: [
                    ['builtin', 'external'],
                    ['internal'],
                    ['parent', 'sibling', 'index'],
                    ['unknown'],
                ],
                pathGroups: [
                    {
                        pattern: 'react|react-dom|react-router-dom',
                        group: 'builtin',
                        position: 'before',
                    },
                    {
                        pattern: '@/**',
                        group: 'internal',
                        position: 'before',
                    },
                    {
                        pattern: '*.{css,scss}',
                        patternOptions: {matchBase: true},
                        group: 'index',
                        position: 'after',
                    },
                ],
                'newlines-between': 'always',
            }],
            'import-x/no-restricted-paths': ['error', {
                basePath: './src',
                zones: [
                    {
                        target: './components/atoms',
                        from: './components',
                        except: ['./atoms'],
                        message: 'components/atoms/** may only import chures and other atoms — not project-wide components.',
                    },
                    {
                        target: './components',
                        from: ['./widgets', './segments', './pages'],
                        message: 'src/components/** is below widgets/segments/pages — it may not import from them (dialogs/ is exempt: any tier may import a dialog to call OpenDialog()).',
                    },
                    {
                        target: ['./widgets', './segments'],
                        from: ['./pages'],
                        message: 'src/widgets/** and src/segments/** are below pages — they may not import from them (dialogs/ is exempt: any tier may import a dialog to call OpenDialog()).',
                    },
                    {
                        target: './dialogs',
                        from: './pages',
                        message: 'Dialogs must not import from pages — pages open dialogs via OpenDialog(), never the other way around.',
                    },
                ],
            }],
        },
    },
    {
        files: ['**/components/atoms/**/*.{ts,tsx}'],
        rules: {
            'no-restricted-syntax': [
                'error',
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
                {
                    selector: 'ObjectPattern[properties.length>6]',
                    message: 'Destructuring more than 6 properties — keep the whole object as one variable (e.g. `const context = useX(...)`) instead of exploding it into separate bindings/props.',
                },
            ],
        },
    },
])
