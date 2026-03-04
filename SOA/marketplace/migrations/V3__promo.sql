CREATE TABLE promo_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code TEXT NOT NULL UNIQUE,
    discount_percent INT NOT NULL CHECK (discount_percent > 0 AND discount_percent <= 100),
    expires_at TIMESTAMPTZ,
    usage_limit INT CHECK (usage_limit IS NULL OR usage_limit > 0),
    used_count INT NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_promo_code_code ON promo_codes(code);
CREATE INDEX idx_promo_code_active ON promo_codes(active);

ALTER TABLE orders
    ADD CONSTRAINT fk_orders_promo
    FOREIGN KEY (promo_id)
    REFERENCES promo_codes(id)
    ON DELETE SET NULL;
