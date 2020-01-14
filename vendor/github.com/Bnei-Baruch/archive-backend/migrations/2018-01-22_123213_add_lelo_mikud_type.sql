-- MDB generated migration file
-- rambler up

WITH data(name) AS (VALUES
  ('LELO_MIKUD')
)
INSERT INTO content_types (name)
  SELECT d.name
  FROM data AS d
  WHERE NOT EXISTS(SELECT ct.name
                   FROM content_types AS ct
                   WHERE ct.name = d.name);

-- rambler down

DELETE FROM content_types
WHERE name IN ('LELO_MIKUD');
