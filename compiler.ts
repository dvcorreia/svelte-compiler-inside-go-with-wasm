import { compile as compileSvelte } from "svelte/compiler"

export type Input = {
  code: string
  path: string
  target: "ssr" | "dom"
  dev: boolean
  css: boolean
}

export type Output =
  | {
      js: string
      css: string
    }
  | { 
      Error: { // Capitalized for Go
        Path: string
        Name: string
        Message: string
        Stack?: string
      }
    }

// Compile svelte code
export function compile(input: string): string {
  const opts: Input = JSON.parse(input)

  const { code, path, target, dev, css } = opts
  const svelte = compileSvelte(code, {
    filename: path,
    generate: target,
    hydratable: true,
    dev: dev,
    css: css,
  })
  return JSON.stringify({
    css: svelte.css.code,
    js: svelte.js.code,
  } as Output)
}
