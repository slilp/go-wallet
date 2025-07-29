CREATE TABLE "transactions" (
    "id" VARCHAR(20) PRIMARY KEY,
    "from" UUID,
    "to" UUID,
    "amount" DECIMAL(20, 2) NOT NULL,
    "type" VARCHAR(20) NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX "idx_transactions_from_created_at" ON "transactions"("from", "created_at" DESC);
CREATE INDEX "idx_transactions_to_created_at" ON "transactions"("to", "created_at" DESC);