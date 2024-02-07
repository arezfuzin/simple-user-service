CREATE TABLE test(
  id serial PRIMARY KEY,
  name varchar(50) UNIQUE NOT NULL
);

INSERT INTO test(name)
  VALUES ('test1');

INSERT INTO test(name)
  VALUES ('test2');

