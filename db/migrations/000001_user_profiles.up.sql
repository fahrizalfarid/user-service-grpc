CREATE TABLE "public"."user_profiles" (
  "id" serial8,
  "firstname" varchar(255) NOT NULL,
  "lastname" varchar(255),
  "email" varchar(255) NOT NULL,
  "created_at" int8,
  "phone" varchar(20) NOT NULL,
  "address" varchar(255) NOT NULL,
  "deleted_at" int8,
  PRIMARY KEY ("id"),
  CONSTRAINT "unq_email" UNIQUE ("email")
);