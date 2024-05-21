CREATE TABLE "accounts" (
  "uuid" varchar PRIMARY KEY,
  "owner" varchar NOT NULL,
  "amount_e5" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "uuid" varchar PRIMARY KEY,
  "account_uuid" varchar NOT NULL,
  "amount_e5" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "uuid" varchar PRIMARY KEY,
  "from_account_uuid" varchar NOT NULL,
  "to_account_uuid" varchar NOT NULL,
  "amount_e5" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "entries" ADD FOREIGN KEY ("account_uuid") REFERENCES "accounts" ("uuid");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_uuid") REFERENCES "accounts" ("uuid");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_uuid") REFERENCES "accounts" ("uuid");

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_uuid");

CREATE INDEX ON "transfers" ("from_account_uuid");

CREATE INDEX ON "transfers" ("to_account_uuid");

CREATE INDEX ON "transfers" ("from_account_uuid", "to_account_uuid");
