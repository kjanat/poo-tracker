import { PrismaClient } from '@prisma/client'

const prisma = new PrismaClient()

async function addConstraints() {
  const constraints = [
    `ALTER TABLE users ADD CONSTRAINT IF NOT EXISTS check_email_format
     CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$')`,

    `ALTER TABLE user_settings ADD CONSTRAINT IF NOT EXISTS check_reminder_time_format
     CHECK (reminder_time ~* '^([01]?[0-9]|2[0-3]):[0-5][0-9]$')`,

    `ALTER TABLE user_settings ADD CONSTRAINT IF NOT EXISTS check_data_retention_positive
     CHECK (data_retention_days > 0)`,

    `ALTER TABLE bowel_movements ADD CONSTRAINT IF NOT EXISTS check_bristol_type_range
     CHECK (bristol_type >= 1 AND bristol_type <= 7)`,

    `ALTER TABLE bowel_movements ADD CONSTRAINT IF NOT EXISTS check_pain_range
     CHECK (pain >= 1 AND pain <= 10)`,

    `ALTER TABLE bowel_movements ADD CONSTRAINT IF NOT EXISTS check_strain_range
     CHECK (strain >= 1 AND strain <= 10)`,

    `ALTER TABLE bowel_movements ADD CONSTRAINT IF NOT EXISTS check_satisfaction_range
     CHECK (satisfaction >= 1 AND satisfaction <= 10)`,

    `ALTER TABLE meals ADD CONSTRAINT IF NOT EXISTS check_spicy_level_range
     CHECK (spicy_level IS NULL OR (spicy_level >= 1 AND spicy_level <= 10))`,

    `ALTER TABLE symptoms ADD CONSTRAINT IF NOT EXISTS check_severity_range
     CHECK (severity >= 1 AND severity <= 10)`
  ]

  for (const constraint of constraints) {
    try {
      await prisma.$executeRawUnsafe(constraint)
      const constraintParts = constraint.split('ADD CONSTRAINT')
      if (constraintParts.length > 1) {
        const checkParts = constraintParts[1]?.split('CHECK')
        if (checkParts && checkParts.length > 0) {
          console.log('✅ Added constraint:', checkParts[0]?.trim())
        }
      }
    } catch (error) {
      console.error('❌ Failed to add constraint:', error)
    }
  }
}

addConstraints()
  .then(() => console.log('🎉 All constraints added successfully!'))
  .catch(console.error)
  .finally(() => prisma.$disconnect())
