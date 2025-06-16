// Quick script to test Zod validation error format
const { z } = require('zod')

const registerSchema = z.object({
  email: z.string().email('Invalid email format'),
  name: z.string().min(1, 'Name is required').optional(),
  password: z.string().min(6, 'Password must be at least 6 characters')
})

try {
  registerSchema.parse({
    email: 'invalid-email',
    password: '123'
  })
} catch (error) {
  console.log('Validation error structure:')
  console.log(JSON.stringify(error.errors, null, 2))
}

try {
  registerSchema.parse({})
} catch (error) {
  console.log('\nMissing fields error structure:')
  console.log(JSON.stringify(error.errors, null, 2))
}
