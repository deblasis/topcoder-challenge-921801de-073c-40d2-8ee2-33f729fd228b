-- This table will be used only once for seeding in order to keep hashing responsility where it belongs
CREATE TABLE if not exists seeding_tmp(
   role VARCHAR(255),
   username VARCHAR(255) NOT NULL UNIQUE,
   password TEXT NOT NULL
);
