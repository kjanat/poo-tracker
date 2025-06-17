import fs from 'fs'
import path from 'path'
import dotenv from 'dotenv'

// Path to the repository root `.env.example`
const envExamplePath = path.resolve(__dirname, '../../.env.example')

if (!fs.existsSync(envExamplePath)) {
  console.error('Unable to locate .env.example file at', envExamplePath)
  process.exit(1)
}

const envExampleContent = fs.readFileSync(envExamplePath, 'utf-8')
const exampleVariables = dotenv.parse(envExampleContent)

const missingVariables: string[] = []

for (const key of Object.keys(exampleVariables)) {
  if (!process.env[key] || process.env[key] === '') {
    missingVariables.push(key)
  }
}

if (missingVariables.length > 0) {
  console.error('Missing required environment variables:\n')
  for (const key of missingVariables) {
    console.error(`- ${key}`)
  }
  process.exit(1)
}
