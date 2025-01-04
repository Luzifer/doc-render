import globals from 'globals'
import js from '@eslint/js'
import typescriptEslint from '@typescript-eslint/eslint-plugin'
import vue from 'eslint-plugin-vue'

export default [
  js.configs.recommended,
  ...vue.configs['flat/recommended'],
  {
    languageOptions: {
      ecmaVersion: 'latest',
      globals: {
        ...globals.browser,
        ...globals.node,
        process: true,
      },

      parserOptions: {
        parser: '@typescript-eslint/parser',
        requireConfigFile: true,
      },
    },

    linterOptions: {
      reportUnusedDisableDirectives: true,
    },

    plugins: {
      '@typescript-eslint': typescriptEslint,
      vue,
    },

    rules: {
      'array-bracket-newline': ['error', { multiline: true }],
      'array-bracket-spacing': ['error'],
      'arrow-body-style': ['error', 'as-needed'],
      'arrow-parens': ['error', 'as-needed'],
      'arrow-spacing': ['error', { after: true, before: true }],
      'block-spacing': ['error'],
      'brace-style': ['error', '1tbs'],
      'camelcase': ['warn'],
      'comma-dangle': ['error', 'always-multiline'],
      'comma-spacing': ['error'],
      'comma-style': ['error', 'last'],
      'curly': ['error'],
      'default-case-last': ['error'],
      'default-param-last': ['error'],
      'dot-location': ['error', 'property'],
      'dot-notation': ['error'],
      'eol-last': ['error', 'always'],
      'eqeqeq': ['error', 'always', { null: 'ignore' }],
      'func-call-spacing': ['error', 'never'],
      'function-paren-newline': ['error', 'multiline'],
      'generator-star-spacing': ['off'],
      'implicit-arrow-linebreak': ['error'],
      'indent': ['error', 2],
      'key-spacing': ['error', { afterColon: true, beforeColon: false, mode: 'strict' }],
      'keyword-spacing': ['error'],
      'linebreak-style': ['error', 'unix'],
      'lines-between-class-members': ['error'],
      'multiline-comment-style': ['off'],
      'newline-per-chained-call': ['error'],
      'no-alert': ['error'],
      'no-console': ['off'],
      'no-debugger': 'off',
      'no-duplicate-imports': ['error'],
      'no-else-return': ['error'],
      'no-empty-function': ['error'],
      'no-extra-parens': ['error'],
      'no-implicit-coercion': ['error'],
      'no-lonely-if': ['error'],
      'no-multi-spaces': ['error'],
      'no-multiple-empty-lines': ['warn', { max: 2, maxBOF: 0, maxEOF: 0 }],
      'no-promise-executor-return': ['error'],
      'no-return-assign': ['error'],
      'no-script-url': ['error'],
      'no-template-curly-in-string': ['error'],
      'no-trailing-spaces': ['error'],
      'no-unneeded-ternary': ['error'],
      'no-unreachable-loop': ['error'],
      'no-unsafe-optional-chaining': ['error'],
      'no-useless-return': ['error'],
      'no-var': ['error'],
      'no-warning-comments': ['error'],
      'no-whitespace-before-property': ['error'],
      'object-curly-newline': ['error', { consistent: true }],
      'object-curly-spacing': ['error', 'always'],
      'object-shorthand': ['error'],
      'padded-blocks': ['error', 'never'],
      'prefer-arrow-callback': ['error'],
      'prefer-const': ['error'],
      'prefer-object-spread': ['error'],
      'prefer-rest-params': ['error'],
      'prefer-template': ['error'],
      'quote-props': ['error', 'consistent-as-needed', { keywords: false }],
      'quotes': ['error', 'single', { allowTemplateLiterals: true }],
      'require-atomic-updates': ['error'],
      'require-await': ['error'],
      'semi': ['error', 'never'],
      'sort-imports': ['error', { ignoreCase: true, ignoreDeclarationSort: false, ignoreMemberSort: false }],
      'sort-keys': ['error', 'asc', { caseSensitive: true, natural: false }],
      'space-before-blocks': ['error', 'always'],
      'space-before-function-paren': ['error', 'never'],
      'space-in-parens': ['error', 'never'],
      'space-infix-ops': ['error'],
      'space-unary-ops': ['error', { nonwords: false, words: true }],
      'spaced-comment': ['warn', 'always'],
      'switch-colon-spacing': ['error'],
      'template-curly-spacing': ['error', 'never'],
      'unicode-bom': ['error', 'never'],
      'vue/comment-directive': 'off',
      'vue/new-line-between-multi-line-property': ['error'],
      'vue/no-empty-component-block': ['error'],
      'vue/no-reserved-component-names': ['error'],
      'vue/no-template-target-blank': ['error'],
      'vue/no-unused-properties': ['error'],
      'vue/no-unused-refs': ['error'],
      'vue/no-useless-mustaches': ['error'],
      'vue/order-in-components': ['off'],
      'vue/require-name-property': ['error'],
      'vue/v-for-delimiter-style': ['error'],
      'vue/v-on-function-call': ['error'],
      'wrap-iife': ['error'],
      'yoda': ['error'],
    },
  },
]
