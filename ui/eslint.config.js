// @ts-check
import eslint from "@eslint/js";
import eslintPluginPrettierRecommended from "eslint-plugin-prettier/recommended";
import tseslint from "typescript-eslint";

export default tseslint.config(
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
  // @ts-expect-error: prettier types are not up-to-date
  eslintPluginPrettierRecommended,
  {
    ignores: ["dist", "public", "src/utils/karaberus.d.ts"],
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
