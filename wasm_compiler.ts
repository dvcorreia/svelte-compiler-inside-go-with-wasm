import { compile } from "./compiler.ts"
import { STDIO, readFileSync, writeFileSync} from "javy/fs"

const inputBytes = readFileSync(STDIO.Stdin);
const inputString = new TextDecoder().decode(inputBytes);

const output = compile(inputString);

const outputBytes = new TextEncoder().encode(output)
writeFileSync(STDIO.Stdout, outputBytes)
