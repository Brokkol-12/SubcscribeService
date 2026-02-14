CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    
    service_name VARCHAR(255) NOT NULL,

    price INTEGER NOT NULL CHECK (price > 0),

    start_date DATE NULL,
    end_date DATE NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_service_name ON subscriptions(service_name);
CREATE INDEX idx_subscriptions_period ON subscriptions(start_date, end_date);
CREATE INDEX idx_subscriptions_active
    ON subscriptions(deleted_at)
    WHERE deleted_at IS NULL;
