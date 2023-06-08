CREATE TABLE "public"."user_credentials" (
  "id" int8 NOT NULL,
  "username" varchar(255) NOT NULL,
  "password" varchar(255) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "unq_username_user_credentials" UNIQUE ("username")
)
;

ALTER TABLE "public"."user_credentials" ADD CONSTRAINT "fk_id_user_credentials_user_profiles" FOREIGN KEY ("id") REFERENCES "public"."user_profiles" ("id") ON DELETE CASCADE ON UPDATE CASCADE;