import { PrismaClient } from '@prisma/client'
import bcrypt from 'bcryptjs'

const prisma = new PrismaClient()

async function main(): Promise<void> {
  console.log('🧹 Cleaning up existing data...')

  // Clean up data in correct order to avoid foreign key constraint issues
  await prisma.mealBowelMovementRelation.deleteMany()
  await prisma.bowelMovementDetails.deleteMany()
  await prisma.bowelMovement.deleteMany()
  await prisma.meal.deleteMany()
  await prisma.userAuth.deleteMany()
  await prisma.user.deleteMany()

  console.log('✅ Cleaned up existing data')

  // Create sample users
  console.log('👤 Creating sample users...')
  const salt = await bcrypt.genSalt(12)
  const hashedPassword = await bcrypt.hash('password123', salt)

  const user1 = await prisma.user.create({
    data: {
      email: 'john.doe@example.com',
      name: 'John Doe',
      auth: {
        create: {
          passwordHash: hashedPassword,
          salt: salt
        }
      }
    }
  })

  // Create additional users for completeness
  await prisma.user.create({
    data: {
      email: 'jane.smith@example.com',
      name: 'Jane Smith',
      auth: {
        create: {
          passwordHash: hashedPassword,
          salt: salt
        }
      }
    }
  })

  await prisma.user.create({
    data: {
      email: 'test@example.com',
      name: 'Test User',
      auth: {
        create: {
          passwordHash: hashedPassword,
          salt: salt
        }
      }
    }
  })

  console.log('✅ Created 3 users')

  // Create sample meals
  console.log('🍽️  Creating sample meals...')
  const meals = [
    {
      name: 'Oatmeal with Berries',
      description: 'Steel-cut oats with fresh blueberries and almonds',
      mealTime: new Date('2025-06-14T07:30:00Z'),
      category: 'BREAKFAST' as const,
      cuisine: 'American',
      spicyLevel: 1,
      fiberRich: true,
      dairy: false,
      gluten: true,
      notes: 'Very filling and nutritious'
    },
    {
      name: 'Spicy Thai Curry',
      description: 'Red curry with vegetables and tofu',
      mealTime: new Date('2025-06-14T12:30:00Z'),
      category: 'LUNCH' as const,
      cuisine: 'Thai',
      spicyLevel: 8,
      fiberRich: true,
      dairy: false,
      gluten: false,
      notes: 'Very spicy but delicious'
    },
    {
      name: 'Grilled Salmon',
      description: 'Atlantic salmon with roasted vegetables',
      mealTime: new Date('2025-06-14T19:00:00Z'),
      category: 'DINNER' as const,
      cuisine: 'Mediterranean',
      spicyLevel: 2,
      fiberRich: true,
      dairy: false,
      gluten: false,
      notes: 'Perfect protein and omega-3s'
    }
  ]

  for (const meal of meals) {
    await prisma.meal.create({
      data: {
        ...meal,
        userId: user1.id
      }
    })
  }

  // Create sample bowel movements
  console.log('💩 Creating sample bowel movements...')
  const bowelMovements = [
    {
      bristolType: 4,
      volume: 'MEDIUM' as const,
      color: 'BROWN' as const,
      consistency: 'SOFT' as const,
      floaters: false,
      pain: 1,
      strain: 1,
      satisfaction: 8,
      smell: 'MILD' as const,
      recordedAt: new Date('2025-06-14T09:15:00Z'),
      notes: 'Perfect morning movement after oatmeal'
    },
    {
      bristolType: 3,
      volume: 'SMALL' as const,
      color: 'DARK_BROWN' as const,
      consistency: 'SOLID' as const,
      floaters: true,
      pain: 3,
      strain: 4,
      satisfaction: 6,
      smell: 'MODERATE' as const,
      recordedAt: new Date('2025-06-13T14:30:00Z'),
      notes: 'A bit hard but manageable'
    },
    {
      bristolType: 6,
      volume: 'LARGE' as const,
      color: 'YELLOW' as const,
      consistency: 'LOOSE' as const,
      floaters: false,
      pain: 7,
      strain: 2,
      satisfaction: 2,
      smell: 'STRONG' as const,
      recordedAt: new Date('2025-06-12T16:20:00Z'),
      notes: 'Emergency bathroom visit! Something did not agree with me.'
    }
  ]

  for (const bowelMovement of bowelMovements) {
    const { notes, ...bowelMovementData } = bowelMovement
    const createdBowelMovement = await prisma.bowelMovement.create({
      data: {
        ...bowelMovementData,
        userId: user1.id
      }
    })

    // Add details separately if notes exist
    if (notes) {
      await prisma.bowelMovementDetails.create({
        data: {
          bowelMovementId: createdBowelMovement.id,
          notes: notes
        }
      })
    }
  }

  // Get counts for summary
  const userCount = await prisma.user.count()
  const mealCount = await prisma.meal.count()
  const bowelMovementCount = await prisma.bowelMovement.count()

  console.log('✅ Seed data created successfully!')
  console.log(`📊 Summary:`)
  console.log(`   - Users: ${userCount}`)
  console.log(`   - Meals: ${mealCount}`)
  console.log(`   - Bowel Movements: ${bowelMovementCount}`)
  console.log('')
  console.log('🔑 Test login credentials:')
  console.log('   Email: john.doe@example.com')
  console.log('   Password: password123')
}

main()
  .catch((e) => {
    console.error('❌ Error running seed:', e)
    process.exit(1)
  })
  .finally(async () => {
    await prisma.$disconnect()
  })
