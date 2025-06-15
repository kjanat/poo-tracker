import { PrismaClient } from '@prisma/client'

const prisma = new PrismaClient()

async function main() {
  console.log('🌱 Starting database seed...')

  // Clean existing data
  console.log('🧹 Cleaning existing data...')
  await prisma.entry.deleteMany()
  await prisma.meal.deleteMany()
  await prisma.user.deleteMany()

  // Create sample users
  console.log('👤 Creating sample users...')
  const user1 = await prisma.user.create({
    data: {
      email: 'john.doe@example.com',
      name: 'John Doe'
    }
  })

  const user2 = await prisma.user.create({
    data: {
      email: 'jane.smith@example.com',
      name: 'Jane Smith'
    }
  })

  const user3 = await prisma.user.create({
    data: {
      email: 'test@example.com',
      name: 'Test User'
    }
  })

  console.log(`✅ Created ${3} users`)

  // Create sample meals for user1
  console.log('🍽️  Creating sample meals...')
  const meals = [
    {
      name: 'Oatmeal with Berries',
      description: 'Steel-cut oats with fresh blueberries and almonds',
      mealTime: new Date('2025-06-14T07:30:00Z'),
      category: 'Breakfast',
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
      category: 'Lunch',
      cuisine: 'Thai',
      spicyLevel: 8,
      fiberRich: true,
      dairy: false,
      gluten: false,
      notes: 'Extremely spicy but delicious'
    },
    {
      name: 'Grilled Salmon',
      description: 'Atlantic salmon with roasted vegetables',
      mealTime: new Date('2025-06-14T19:00:00Z'),
      category: 'Dinner',
      cuisine: 'Mediterranean',
      spicyLevel: 2,
      fiberRich: true,
      dairy: false,
      gluten: false,
      notes: 'Perfect protein and omega-3s'
    },
    {
      name: 'Fast Food Burger',
      description: 'Double cheeseburger with fries',
      mealTime: new Date('2025-06-13T13:00:00Z'),
      category: 'Lunch',
      cuisine: 'Fast Food',
      spicyLevel: 3,
      fiberRich: false,
      dairy: true,
      gluten: true,
      notes: 'Guilty pleasure meal'
    },
    {
      name: 'Vegetable Smoothie',
      description: 'Kale, spinach, banana, and protein powder',
      mealTime: new Date('2025-06-13T08:00:00Z'),
      category: 'Breakfast',
      cuisine: 'Health Food',
      spicyLevel: 1,
      fiberRich: true,
      dairy: true,
      gluten: false,
      notes: 'Green and mean!'
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

  // Create sample meals for user2
  const user2Meals = [
    {
      name: 'Avocado Toast',
      description: 'Sourdough with smashed avocado and everything seasoning',
      mealTime: new Date('2025-06-15T08:30:00Z'),
      category: 'Breakfast',
      cuisine: 'California',
      spicyLevel: 2,
      fiberRich: true,
      dairy: false,
      gluten: true,
      notes: 'Instagram worthy'
    },
    {
      name: 'Pepperoni Pizza',
      description: 'Large pepperoni pizza with extra cheese',
      mealTime: new Date('2025-06-15T20:00:00Z'),
      category: 'Dinner',
      cuisine: 'Italian',
      spicyLevel: 4,
      fiberRich: false,
      dairy: true,
      gluten: true,
      notes: 'Friday night treat'
    }
  ]

  for (const meal of user2Meals) {
    await prisma.meal.create({
      data: {
        ...meal,
        userId: user2.id
      }
    })
  }

  console.log(`✅ Created ${meals.length + user2Meals.length} meals`)

  // Create sample entries for user1
  console.log('💩 Creating sample poop entries...')
  const entries = [
    {
      bristolType: 4,
      volume: 'Medium',
      color: 'Brown',
      consistency: 'Soft',
      floaters: false,
      pain: 1,
      strain: 2,
      satisfaction: 8,
      smell: 'Mild',
      notes: 'Perfect consistency after the oatmeal breakfast',
      createdAt: new Date('2025-06-14T09:15:00Z')
    },
    {
      bristolType: 6,
      volume: 'Large',
      color: 'Dark Brown',
      consistency: 'Loose',
      floaters: true,
      pain: 3,
      strain: 1,
      satisfaction: 4,
      smell: 'Strong',
      notes: 'That spicy curry definitely had an impact! 🌶️',
      createdAt: new Date('2025-06-14T14:30:00Z')
    },
    {
      bristolType: 3,
      volume: 'Small',
      color: 'Light Brown',
      consistency: 'Solid',
      floaters: false,
      pain: 2,
      strain: 4,
      satisfaction: 6,
      smell: 'Moderate',
      notes: 'Needed more fiber today',
      createdAt: new Date('2025-06-13T10:00:00Z')
    },
    {
      bristolType: 5,
      volume: 'Large',
      color: 'Brown',
      consistency: 'Soft',
      floaters: false,
      pain: 1,
      strain: 1,
      satisfaction: 9,
      smell: 'Mild',
      notes: 'After the salmon dinner - excellent!',
      createdAt: new Date('2025-06-14T21:45:00Z')
    },
    {
      bristolType: 2,
      volume: 'Small',
      color: 'Dark Brown',
      consistency: 'Solid',
      floaters: false,
      pain: 4,
      strain: 8,
      satisfaction: 3,
      smell: 'Strong',
      notes: 'Too much fast food lately 😣',
      createdAt: new Date('2025-06-13T15:30:00Z')
    },
    {
      bristolType: 7,
      volume: 'Medium',
      color: 'Yellow',
      consistency: 'Watery',
      floaters: true,
      pain: 6,
      strain: 1,
      satisfaction: 1,
      smell: 'Toxic',
      notes: 'Emergency bathroom visit! Something did not agree with me.',
      createdAt: new Date('2025-06-12T16:20:00Z')
    }
  ]

  for (const entry of entries) {
    await prisma.entry.create({
      data: {
        ...entry,
        userId: user1.id
      }
    })
  }

  // Create sample entries for user2
  const user2Entries = [
    {
      bristolType: 4,
      volume: 'Medium',
      color: 'Brown',
      consistency: 'Soft',
      floaters: false,
      pain: 1,
      strain: 2,
      satisfaction: 8,
      smell: 'None',
      notes: 'Morning routine after avocado toast',
      createdAt: new Date('2025-06-15T10:30:00Z')
    },
    {
      bristolType: 5,
      volume: 'Large',
      color: 'Brown',
      consistency: 'Soft',
      floaters: false,
      pain: 2,
      strain: 3,
      satisfaction: 7,
      smell: 'Moderate',
      notes: 'Pizza aftermath - not terrible!',
      createdAt: new Date('2025-06-15T22:15:00Z')
    }
  ]

  for (const entry of user2Entries) {
    await prisma.entry.create({
      data: {
        ...entry,
        userId: user2.id
      }
    })
  }

  // Create a few entries for test user
  const testUserEntries = [
    {
      bristolType: 4,
      volume: 'Medium',
      color: 'Brown',
      consistency: 'Soft',
      floaters: false,
      pain: 1,
      strain: 1,
      satisfaction: 9,
      smell: 'Mild',
      notes: 'Perfect test entry',
      createdAt: new Date('2025-06-15T12:00:00Z')
    }
  ]

  for (const entry of testUserEntries) {
    await prisma.entry.create({
      data: {
        ...entry,
        userId: user3.id
      }
    })
  }

  console.log(`✅ Created ${entries.length + user2Entries.length + testUserEntries.length} entries`)

  // Summary
  const userCount = await prisma.user.count()
  const mealCount = await prisma.meal.count()
  const entryCount = await prisma.entry.count()

  console.log('\n🎉 Database seeding completed!')
  console.log(`📊 Summary:`)
  console.log(`   👤 Users: ${userCount}`)
  console.log(`   🍽️  Meals: ${mealCount}`)
  console.log(`   💩 Entries: ${entryCount}`)
  console.log('\n🚽 Your poo-tracker database is ready to roll!')
}

main()
  .catch((e) => {
    console.error('❌ Error seeding database:', e)
    process.exit(1)
  })
  .finally(async () => {
    await prisma.$disconnect()
  })
