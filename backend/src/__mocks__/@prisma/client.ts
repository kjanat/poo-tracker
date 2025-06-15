export class PrismaClient {
  user = {
    findUnique: jest.fn(),
    create: jest.fn(),
    update: jest.fn()
  }
  userAuth = {
    findUnique: jest.fn(),
    update: jest.fn()
  }
}
