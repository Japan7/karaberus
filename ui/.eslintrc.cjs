/* eslint-env node */
/** @type {import('eslint').Linter.Config} */
module.exports = {
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:prettier/recommended",
  ],
  parser: "@typescript-eslint/parser",
  plugins: ["@typescript-eslint"],
  root: true,
  ignorePatterns: ["dist"],
  rules: {
    "@typescript-eslint/no-unused-vars": [
      "warn",
      { varsIgnorePattern: "^_", args: "none" },
    ],
  },
};
