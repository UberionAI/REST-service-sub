CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE sub (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        service_name TEXT NOT NULL,
        price INTEGER NOT NULL,
        user_id UUID NOT NULL,
        start_date DATE NOT NULL,
        end_date DATE NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS "sub_user_id_idx" ON "sub" ("user_id");
CREATE INDEX IF NOT EXISTS "sub_service_name_idx" ON "sub" ("service_name");
Create INDEX IF NOT EXISTS "sub_start_date_idx" ON "sub" ("start_date");