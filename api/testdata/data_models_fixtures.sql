-- Fixed views and watching-now counts for the core tests; hidden01 is unpublished, so it must never be returned.
INSERT INTO dwh_content_units_measures (event_unit_uid, total_page_enter_count, unique_users_count, unique_users_watching_now_count) VALUES
  ('lesson11', 500, 300, 80),
  ('lesson12', 400, 250, 60),
  ('chaptr11', 300, 200, 0),
  ('clipunt1', 150, 100, 55),
  ('hidden01', 900, 900, 900);
