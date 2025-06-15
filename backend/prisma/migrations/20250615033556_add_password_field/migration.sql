/*
  Warnings:

  - Added the required column `password` to the `users` table without a default value. This is not possible if the table is not empty.

*/

-- AlterTable: Add password column with a default value first
ALTER TABLE "users" ADD COLUMN "password" TEXT DEFAULT '$2a$12$LQv3c1yqBwlVHpPyibUye.Q7caFLlaPLSklGWOO.7r3gHzfqJHOzO';

-- Update the column to be NOT NULL and remove the default
ALTER TABLE "users" ALTER COLUMN "password" SET NOT NULL;
ALTER TABLE "users" ALTER COLUMN "password" DROP DEFAULT;
