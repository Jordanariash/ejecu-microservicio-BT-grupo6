import react from "@vitejs/plugin-react"
import { defineConfig } from "vitest/config"

export default defineConfig({
  plugins: [react()],
  test: {
    // lcov es el formato que lee SonarQube
    coverage: { provider: "v8", reporter: ["text", "lcov"] },
  },
})
