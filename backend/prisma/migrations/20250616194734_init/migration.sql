-- CreateEnum
CREATE TYPE "audit_action_enum" AS ENUM ('CREATE', 'UPDATE', 'DELETE');

-- CreateEnum
CREATE TYPE "volume_enum" AS ENUM ('SMALL', 'MEDIUM', 'LARGE', 'MASSIVE');

-- CreateEnum
CREATE TYPE "color_enum" AS ENUM ('BROWN', 'DARK_BROWN', 'LIGHT_BROWN', 'YELLOW', 'GREEN', 'RED', 'BLACK');

-- CreateEnum
CREATE TYPE "consistency_enum" AS ENUM ('SOLID', 'SOFT', 'LOOSE', 'WATERY');

-- CreateEnum
CREATE TYPE "smell_level_enum" AS ENUM ('NONE', 'MILD', 'MODERATE', 'STRONG', 'TOXIC');

-- CreateEnum
CREATE TYPE "meal_category_enum" AS ENUM ('BREAKFAST', 'LUNCH', 'DINNER', 'SNACK');

-- CreateEnum
CREATE TYPE "symptom_type_enum" AS ENUM ('BLOATING', 'CRAMPS', 'NAUSEA', 'HEARTBURN', 'CONSTIPATION', 'DIARRHEA', 'GAS', 'FATIGUE', 'OTHER');

-- CreateTable
CREATE TABLE "users" (
    "id" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "name" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "user_auth" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "passwordHash" TEXT NOT NULL,
    "salt" TEXT NOT NULL,
    "loginAttempts" INTEGER NOT NULL DEFAULT 0,
    "lockedUntil" TIMESTAMP(3),
    "lastLogin" TIMESTAMP(3),
    "resetToken" TEXT,
    "resetTokenExpiry" TIMESTAMP(3),
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "user_auth_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "user_two_factor" (
    "id" TEXT NOT NULL,
    "userAuthId" TEXT NOT NULL,
    "secret" TEXT NOT NULL,
    "enabled" BOOLEAN NOT NULL DEFAULT false,
    "backupCodes" TEXT[],
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "user_two_factor_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "user_settings" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "timezone" TEXT NOT NULL DEFAULT 'UTC',
    "reminderEnabled" BOOLEAN NOT NULL DEFAULT true,
    "reminderTime" TEXT NOT NULL DEFAULT '09:00',
    "dataRetentionDays" INTEGER NOT NULL DEFAULT 365,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "user_settings_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "bowel_movements" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "recordedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "bristolType" SMALLINT NOT NULL,
    "volume" "volume_enum",
    "color" "color_enum",
    "consistency" "consistency_enum",
    "floaters" BOOLEAN NOT NULL DEFAULT false,
    "pain" SMALLINT NOT NULL DEFAULT 1,
    "strain" SMALLINT NOT NULL DEFAULT 1,
    "satisfaction" SMALLINT NOT NULL DEFAULT 5,
    "photoUrl" TEXT,
    "smell" "smell_level_enum",

    CONSTRAINT "bowel_movements_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "bowel_movement_details" (
    "id" TEXT NOT NULL,
    "bowelMovementId" TEXT NOT NULL,
    "notes" TEXT,
    "aiAnalysis" JSONB,

    CONSTRAINT "bowel_movement_details_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "meals" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "mealTime" TIMESTAMP(3) NOT NULL,
    "category" "meal_category_enum",
    "cuisine" TEXT,
    "spicyLevel" SMALLINT,
    "fiberRich" BOOLEAN NOT NULL DEFAULT false,
    "dairy" BOOLEAN NOT NULL DEFAULT false,
    "gluten" BOOLEAN NOT NULL DEFAULT false,
    "photoUrl" TEXT,
    "notes" TEXT,

    CONSTRAINT "meals_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "meal_bowel_movement_relations" (
    "id" TEXT NOT NULL,
    "mealId" TEXT NOT NULL,
    "bowelMovementId" TEXT NOT NULL,

    CONSTRAINT "meal_bowel_movement_relations_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "symptoms" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "bowelMovementId" TEXT,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "recordedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "type" "symptom_type_enum" NOT NULL,
    "severity" SMALLINT NOT NULL,
    "notes" TEXT,

    CONSTRAINT "symptoms_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "medications" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,
    "name" TEXT NOT NULL,
    "dosage" TEXT,
    "frequency" TEXT,
    "startDate" TIMESTAMP(3) NOT NULL,
    "endDate" TIMESTAMP(3),
    "notes" TEXT,

    CONSTRAINT "medications_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "meal_symptom_relations" (
    "id" TEXT NOT NULL,
    "mealId" TEXT NOT NULL,
    "symptomId" TEXT NOT NULL,

    CONSTRAINT "meal_symptom_relations_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "audit_logs" (
    "id" TEXT NOT NULL,
    "userId" TEXT NOT NULL,
    "tableName" TEXT NOT NULL,
    "recordId" TEXT NOT NULL,
    "action" "audit_action_enum" NOT NULL,
    "oldValues" JSONB,
    "newValues" JSONB,
    "timestamp" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "ipAddress" TEXT,
    "userAgent" TEXT,

    CONSTRAINT "audit_logs_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "users_email_key" ON "users"("email");

-- CreateIndex
CREATE INDEX "users_createdAt_idx" ON "users"("createdAt");

-- CreateIndex
CREATE UNIQUE INDEX "user_auth_userId_key" ON "user_auth"("userId");

-- CreateIndex
CREATE INDEX "user_auth_resetToken_idx" ON "user_auth"("resetToken");

-- CreateIndex
CREATE UNIQUE INDEX "user_two_factor_userAuthId_key" ON "user_two_factor"("userAuthId");

-- CreateIndex
CREATE UNIQUE INDEX "user_settings_userId_key" ON "user_settings"("userId");

-- CreateIndex
CREATE INDEX "bowel_movements_userId_createdAt_idx" ON "bowel_movements"("userId", "createdAt");

-- CreateIndex
CREATE INDEX "bowel_movements_userId_recordedAt_idx" ON "bowel_movements"("userId", "recordedAt");

-- CreateIndex
CREATE INDEX "bowel_movements_createdAt_idx" ON "bowel_movements"("createdAt");

-- CreateIndex
CREATE INDEX "bowel_movements_bristolType_idx" ON "bowel_movements"("bristolType");

-- CreateIndex
CREATE INDEX "bowel_movements_pain_idx" ON "bowel_movements"("pain");

-- CreateIndex
CREATE INDEX "bowel_movements_satisfaction_idx" ON "bowel_movements"("satisfaction");

-- CreateIndex
CREATE UNIQUE INDEX "bowel_movement_details_bowelMovementId_key" ON "bowel_movement_details"("bowelMovementId");

-- CreateIndex
CREATE INDEX "meals_userId_mealTime_idx" ON "meals"("userId", "mealTime");

-- CreateIndex
CREATE INDEX "meals_mealTime_idx" ON "meals"("mealTime");

-- CreateIndex
CREATE UNIQUE INDEX "meal_bowel_movement_relations_mealId_bowelMovementId_key" ON "meal_bowel_movement_relations"("mealId", "bowelMovementId");

-- CreateIndex
CREATE INDEX "symptoms_userId_createdAt_idx" ON "symptoms"("userId", "createdAt");

-- CreateIndex
CREATE INDEX "symptoms_userId_recordedAt_idx" ON "symptoms"("userId", "recordedAt");

-- CreateIndex
CREATE INDEX "symptoms_type_idx" ON "symptoms"("type");

-- CreateIndex
CREATE INDEX "symptoms_severity_idx" ON "symptoms"("severity");

-- CreateIndex
CREATE INDEX "medications_userId_startDate_idx" ON "medications"("userId", "startDate");

-- CreateIndex
CREATE INDEX "medications_userId_endDate_idx" ON "medications"("userId", "endDate");

-- CreateIndex
CREATE UNIQUE INDEX "meal_symptom_relations_mealId_symptomId_key" ON "meal_symptom_relations"("mealId", "symptomId");

-- CreateIndex
CREATE INDEX "audit_logs_userId_timestamp_idx" ON "audit_logs"("userId", "timestamp");

-- CreateIndex
CREATE INDEX "audit_logs_tableName_recordId_idx" ON "audit_logs"("tableName", "recordId");

-- CreateIndex
CREATE INDEX "audit_logs_timestamp_idx" ON "audit_logs"("timestamp");

-- AddForeignKey
ALTER TABLE "user_auth" ADD CONSTRAINT "user_auth_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "user_two_factor" ADD CONSTRAINT "user_two_factor_userAuthId_fkey" FOREIGN KEY ("userAuthId") REFERENCES "user_auth"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "user_settings" ADD CONSTRAINT "user_settings_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "bowel_movements" ADD CONSTRAINT "bowel_movements_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "bowel_movement_details" ADD CONSTRAINT "bowel_movement_details_bowelMovementId_fkey" FOREIGN KEY ("bowelMovementId") REFERENCES "bowel_movements"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "meals" ADD CONSTRAINT "meals_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "meal_bowel_movement_relations" ADD CONSTRAINT "meal_bowel_movement_relations_mealId_fkey" FOREIGN KEY ("mealId") REFERENCES "meals"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "meal_bowel_movement_relations" ADD CONSTRAINT "meal_bowel_movement_relations_bowelMovementId_fkey" FOREIGN KEY ("bowelMovementId") REFERENCES "bowel_movements"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "symptoms" ADD CONSTRAINT "symptoms_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "symptoms" ADD CONSTRAINT "symptoms_bowelMovementId_fkey" FOREIGN KEY ("bowelMovementId") REFERENCES "bowel_movements"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "medications" ADD CONSTRAINT "medications_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "meal_symptom_relations" ADD CONSTRAINT "meal_symptom_relations_mealId_fkey" FOREIGN KEY ("mealId") REFERENCES "meals"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "meal_symptom_relations" ADD CONSTRAINT "meal_symptom_relations_symptomId_fkey" FOREIGN KEY ("symptomId") REFERENCES "symptoms"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "audit_logs" ADD CONSTRAINT "audit_logs_userId_fkey" FOREIGN KEY ("userId") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
