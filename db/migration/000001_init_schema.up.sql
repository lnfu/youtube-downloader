CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "file_type" AS ENUM ('audio', 'video');

CREATE TYPE "file_status" AS ENUM ('in_progress', 'completed', 'error');

CREATE TABLE "media" (
    "id" int PRIMARY KEY NOT NULL,
    "title" varchar(100),
    "duration" int
);

CREATE TABLE "file" (
    "id" uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
    "media_id" int NOT NULL,
    "type" file_type NOT NULL,
    "status" file_status NOT NULL,
    "access_key" varchar(16) NOT NULL,
    "access_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE
    "file"
ADD
    FOREIGN KEY ("media_id") REFERENCES "media" ("id");