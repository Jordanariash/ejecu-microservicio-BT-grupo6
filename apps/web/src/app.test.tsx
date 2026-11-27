import { renderToString } from "react-dom/server"
import { expect, test } from "vitest"
import { App } from "./app"

test("muestra el título de la base técnica", () => {
  const html = renderToString(<App />)
  expect(html).toContain("Base técnica · Grupo 6")
})
