CREATE TABLE "transactions" (
    "id" VARCHAR(20) PRIMARY KEY,
    "from" UUID NOT NULL,
    "to" UUID NOT NULL,
    "amount" DECIMAL(20, 2) NOT NULL,
    "description" TEXT,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY ("from") REFERENCES "wallets"("id") ON DELETE CASCADE,
    FOREIGN KEY ("to") REFERENCES "wallets"("id") ON DELETE CASCADE
);

CREATE INDEX "idx_transactions_from" ON "transactions"("from");
CREATE INDEX "idx_transactions_to" ON "transactions"("to");
CREATE INDEX "idx_transactions_updated_at" ON "transactions"("updated_at" DESC);
