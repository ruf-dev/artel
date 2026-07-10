export default {
    ignoreFiles: ['dist/**'],
    extends: [
        'stylelint-config-standard',
        'stylelint-config-standard-scss',
        'stylelint-config-css-modules',
    ],
    rules: {
        // No named colors — always use CSS variables
        'color-named': 'never',

        // No !important
        'declaration-no-important': true,

        // Allow blank lines between property groups (common grouping pattern)
        'declaration-empty-line-before': null,

        // Don't require shorthand when longhand is clearer
        'declaration-block-no-redundant-longhand-properties': null,

        // Allow single-line keyframe blocks (common shorthand)
        'declaration-block-single-line-max-declarations': null,

        // Allow camelCase keyframe names (e.g. slideIn, fadeOut)
        'keyframes-name-pattern': null,

        // Vendor prefixes still needed for backdrop-filter etc.
        'property-no-vendor-prefix': null,

        // Allow more than 4 decimal places in calculated values
        'number-max-precision': null,

        // Sizes must be expressed in rem — no px/em/vh/vw/etc. Warning only: ~580
        // pre-existing occurrences are a known baseline (see sizes.css), fix the
        // ones you touch rather than sweeping the whole codebase.
        'unit-disallowed-list': [
            ['px', 'em', 'ex', 'ch', 'vh', 'vw', 'vmin', 'vmax', 'cm', 'mm', 'in', 'pt', 'pc', 'q'],
            {severity: 'warning'},
        ],
    },
    overrides: [
        // Component CSS Modules: enforce naming and unit conventions
        {
            files: ['src/**/*.module.{css,scss}'],
            rules: {
                // No px/em for font-size — use rem
                'declaration-property-unit-disallowed-list': {
                    'font-size': ['px', 'em'],
                },

                // Never use z-index — rely on DOM order or a portal to document.body instead (see CLAUDE.md)
                'declaration-property-value-disallowed-list': {
                    'z-index': [/.*/],
                },

                // PascalCase or camelCase class names (matches CSS Modules usage)
                'selector-class-pattern': '^[A-Za-z][a-zA-Z0-9]*$',

                // `composes` references another class by its exact (case-sensitive) name —
                // not a lowercasable value keyword
                'value-keyword-case': [
                    'lower',
                    {ignoreProperties: ['composes']},
                ],
            },
        },
        // Global token files: hex color definitions and base sizing are the source of truth here
        {
            files: ['src/index.css', 'src/colors_and_type.css', 'src/sizes.css'],
            rules: {
                'color-named': null,
                'declaration-property-unit-disallowed-list': null,
                'selector-class-pattern': null,
            },
        },
        // sizes.css is the one place raw px/vh/vw values are allowed — it's the
        // source of truth that every other file's rem-based size tokens derive from
        {
            files: ['src/sizes.css'],
            rules: {
                'unit-disallowed-list': null,
            },
        },
    ],
}
