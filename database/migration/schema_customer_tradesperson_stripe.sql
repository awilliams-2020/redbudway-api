-- Stripe Customer on a Connect account (direct charges). One row per (customer, tradesperson).
-- Apply to redbudway / redbudway_dev when rolling out direct charges.

CREATE TABLE IF NOT EXISTS customer_tradesperson_stripe (
  customerId CHAR(36) NOT NULL,
  tradespersonId CHAR(36) NOT NULL,
  stripeCustomerId VARCHAR(64) NOT NULL,
  createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (customerId, tradespersonId),
  KEY idx_tp_stripe_cust (tradespersonId, stripeCustomerId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
