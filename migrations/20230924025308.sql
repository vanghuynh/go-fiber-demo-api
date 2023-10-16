-- Create "products" table
CREATE TABLE "public"."products" (
  "id" bigserial NOT NULL,
  "name" text NULL,
  "created_at" timestamptz NULL,
  "serial_number" text NULL,
  PRIMARY KEY ("id")
);
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "first_name" text NULL,
  "last_name" text NULL,
  PRIMARY KEY ("id")
);
-- Create "orders" table
CREATE TABLE "public"."orders" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "product_refer" bigint NULL,
  "user_refer" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_orders_product" FOREIGN KEY ("product_refer") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_orders_user" FOREIGN KEY ("user_refer") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
