-- Enable the UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE test(
  id serial PRIMARY KEY,
  name varchar(50) UNIQUE NOT NULL
);

CREATE TABLE users(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  phone_number varchar(15) UNIQUE NOT NULL,
  full_name varchar(100) NOT NULL,
  password_hash varchar(100) NOT NULL
);

INSERT INTO test(name)
  VALUES ('test1');

INSERT INTO test(name)
  VALUES ('test2');

INSERT INTO users(phone_number, full_name, password_hash)
  VALUES ('+1234567890', 'John Doe', 'password');

