/*
Warnings:

- You are about to drop the column `password` on the `users` table. All the data in the column will be lost.

*/
-- CreateTable: Create user_auth table first
CREATE TABLE "user_auth" (
  "id" TEXT NOT NULL,
  "userId" TEXT NOT NULL,
  "password" TEXT NOT NULL,
  "loginAttempts" INTEGER NOT NULL DEFAULT 0,
  "lockedUntil" TIMESTAMP(3),
  "lastLogin" TIMESTAMP(3),
  "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updatedAt" TIMESTAMP(3) NOT NULL,
  CONSTRAINT "user_auth_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "user_auth_userId_key" ON "user_auth" ("userId");

-- AddForeignKey
ALTER TABLE "user_auth" ADD CONSTRAINT "user_auth_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users" ("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- Migrate existing password data
INSERT INTO
  "user_auth" (
    "id",
    "userId",
    "password",
    "createdAt",
    "updatedAt"
  )
SELECT
  'cauth_' || "id" as id,
  "id" as "userId",
  "password",
  "createdAt",
  "updatedAt"
FROM
  "users"
WHERE
  "password" IS NOT NULL;

-- AlterTable: Drop password column after data migration
ALTER TABLE "users"
DROP COLUMN "password";
