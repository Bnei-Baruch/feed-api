-- Fixed MDB content for the core tests, on top of the migrations (content types, source types, the Rav person).
-- Dates are relative to now(), so "last" order and age filters keep working.
--   lesson<d><p>: 4 daily lessons (dlesson<d>, d days ago), parts 0-2 (part 0 is a short preparation), source srcpart<p+1>.
--   chaptr<g><n>: chapters 1-4 of programs program1-2, tag tagprog<g>.
--   clipunt<n>:   clips 1-4 in collection clipcol1, tag tagprog2.
--   articl0<n>:   2 articles.
--   blogpst<n>:   3 blog posts (the app feed shows recent ones), mp3 in ru and en only.
--   hidden01 (unpublished), hidden02 (secure): lesson parts that must never be returned.
-- Every other unit has mp4 and mp3 files in he, ru and en; 3 blog posts in each blog (the migrations add the blogs).
CREATE FUNCTION pg_temp.ct(text) RETURNS bigint AS 'SELECT id FROM content_types WHERE name = $1' LANGUAGE sql;
CREATE FUNCTION pg_temp.days_ago(int) RETURNS text AS 'SELECT (now() - $1 * interval ''1 day'')::date::text' LANGUAGE sql;

INSERT INTO sources (uid, parent_id, type_id, name, position) VALUES
  ('srcbook1', NULL, (SELECT id FROM source_types WHERE name = 'BOOK'), 'Book', 1);
INSERT INTO sources (uid, parent_id, type_id, name, position)
SELECT 'srcpart' || n, (SELECT id FROM sources WHERE uid = 'srcbook1'), (SELECT id FROM source_types WHERE name = 'PART'), 'Part ' || n, n
FROM generate_series(1, 3) n;

INSERT INTO tags (uid, pattern) VALUES ('tagprog1', 'tag-program-1'), ('tagprog2', 'tag-program-2');

INSERT INTO collections (uid, type_id, published, properties)
SELECT 'dlesson' || d, pg_temp.ct('DAILY_LESSON'), true, jsonb_build_object('film_date', pg_temp.days_ago(d), 'number', 1, 'original_language', 'he')
FROM generate_series(1, 4) d;
INSERT INTO collections (uid, type_id, published, properties) VALUES
  ('program1', pg_temp.ct('VIDEO_PROGRAM'), true, '{"active": true, "default_language": "he"}'),
  ('program2', pg_temp.ct('VIDEO_PROGRAM'), true, '{"active": true, "default_language": "he"}'),
  ('clipcol1', pg_temp.ct('CLIPS'), true, '{"active": true, "default_language": "he"}');

INSERT INTO content_units (uid, type_id, published, secure, properties)
SELECT 'lesson' || d || p, pg_temp.ct('LESSON_PART'), true, 0,
       jsonb_build_object('film_date', pg_temp.days_ago(d), 'part', p, 'duration', CASE p WHEN 0 THEN 600 ELSE 3600 END, 'original_language', 'he')
FROM generate_series(1, 4) d, generate_series(0, 2) p
UNION ALL
SELECT 'chaptr' || g || n, pg_temp.ct('VIDEO_PROGRAM_CHAPTER'), true, 0,
       jsonb_build_object('film_date', pg_temp.days_ago(n + 2 * g), 'duration', 1800, 'original_language', 'he')
FROM generate_series(1, 2) g, generate_series(1, 4) n
UNION ALL
SELECT 'clipunt' || n, pg_temp.ct('CLIP'), true, 0, jsonb_build_object('film_date', pg_temp.days_ago(n * 3), 'duration', 180, 'original_language', 'he')
FROM generate_series(1, 4) n
UNION ALL
SELECT 'articl0' || n, pg_temp.ct('ARTICLE'), true, 0, jsonb_build_object('film_date', pg_temp.days_ago(n), 'original_language', 'he')
FROM generate_series(1, 2) n
UNION ALL
SELECT 'blogpst' || n, pg_temp.ct('BLOG_POST'), true, 0, jsonb_build_object('film_date', pg_temp.days_ago(n), 'original_language', 'en')
FROM generate_series(1, 3) n
UNION ALL
VALUES ('hidden01', pg_temp.ct('LESSON_PART'), false, 0, jsonb_build_object('film_date', pg_temp.days_ago(1), 'part', 3, 'duration', 3600, 'original_language', 'he')),
       ('hidden02', pg_temp.ct('LESSON_PART'), true, 1, jsonb_build_object('film_date', pg_temp.days_ago(1), 'part', 4, 'duration', 3600, 'original_language', 'he'));

INSERT INTO collections_content_units (collection_id, content_unit_id, name, position)
SELECT c.id, cu.id, x.pos::text, x.pos
FROM (SELECT 'dlesson' || d AS c, 'lesson' || d || p AS cu, p AS pos FROM generate_series(1, 4) d, generate_series(0, 2) p
      UNION ALL SELECT 'program' || g, 'chaptr' || g || n, n FROM generate_series(1, 2) g, generate_series(1, 4) n
      UNION ALL SELECT 'clipcol1', 'clipunt' || n, n FROM generate_series(1, 4) n
      UNION ALL VALUES ('dlesson1', 'hidden01', 3), ('dlesson1', 'hidden02', 4)) x
JOIN collections c ON c.uid = x.c JOIN content_units cu ON cu.uid = x.cu;

INSERT INTO content_units_sources (content_unit_id, source_id)
SELECT cu.id, s.id FROM content_units cu JOIN sources s ON s.uid = 'srcpart' || ((cu.properties->>'part')::int % 3 + 1)
WHERE cu.type_id = pg_temp.ct('LESSON_PART');

INSERT INTO content_units_tags (content_unit_id, tag_id)
SELECT cu.id, t.id FROM content_units cu JOIN tags t ON t.uid = CASE WHEN cu.uid LIKE 'chaptr1%' THEN 'tagprog1' ELSE 'tagprog2' END
WHERE cu.uid LIKE 'chaptr%' OR cu.uid LIKE 'clipunt%';

INSERT INTO files (uid, name, size, type, sub_type, mime_type, content_unit_id, language, published)
SELECT left(md5(cu.uid || l || m), 8), cu.uid || '_' || l || CASE m WHEN 'video/mp4' THEN '.mp4' ELSE '.mp3' END,
       1, split_part(m, '/', 1), '', m, cu.id, l, true
FROM content_units cu, unnest(ARRAY['he', 'ru', 'en']) l, unnest(ARRAY['video/mp4', 'audio/mpeg']) m
WHERE cu.type_id <> pg_temp.ct('BLOG_POST') OR l <> 'he' AND m = 'audio/mpeg';

INSERT INTO blog_posts (blog_id, wp_id, title, content, posted_at)
SELECT b.id, 100 + n, 'Post ' || n, '', now() - n * interval '1 day' FROM blogs b, generate_series(1, 3) n;
