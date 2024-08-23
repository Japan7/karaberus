// @ts-check
import eslint from "@eslint/js";
import eslintPluginPrettierRecommended from "eslint-plugin-prettier/recommended";
import tseslint from "typescript-eslint";

export default tseslint.config(
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
  eslintPluginPrettierRecommended,
  {
    ignores: ["dist", "src/utils/karaberus.d.ts", "src-tauri"],
  },
  {
    rules: {
      "@typescript-eslint/no-unused-vars": [
        "warn",
        { varsIgnorePattern: "^_", args: "none" },
      ],
    },
  },
);
