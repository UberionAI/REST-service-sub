CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE subscriptions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        service_name TEXT NOT NULL,
        price INTEGER NOT NULL,
        user_id UUID NOT NULL,
        start_date DATE NOT NULL,
        end_date DATE,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "subscriptions_user_id_idx" ON "subscriptions" ("user_id");
CREATE INDEX IF NOT EXISTS "subscriptions_service_name_idx" ON "subscriptions" ("service_name");
Create INDEX IF NOT EXISTS "subscriptions_start_date_idx" ON "subscriptions" ("start_date");