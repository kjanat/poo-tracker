-- CreateTable
CREATE TABLE "meal_entry_relations" (
    "id" TEXT NOT NULL,
    "mealId" TEXT NOT NULL,
    "entryId" TEXT NOT NULL,

    CONSTRAINT "meal_entry_relations_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "meal_entry_relations_mealId_entryId_key" ON "meal_entry_relations"("mealId", "entryId");

-- AddForeignKey
ALTER TABLE "meal_entry_relations" ADD CONSTRAINT "meal_entry_relations_mealId_fkey" FOREIGN KEY ("mealId") REFERENCES "meals"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "meal_entry_relations" ADD CONSTRAINT "meal_entry_relations_entryId_fkey" FOREIGN KEY ("entryId") REFERENCES "entries"("id") ON DELETE CASCADE ON UPDATE CASCADE;
