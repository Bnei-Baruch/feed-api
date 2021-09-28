CREATE INDEX ind_dwh_fact_play_units_by_minutes_event_user_id ON dwh_fact_play_units_by_minutes USING btree (event_unit_uid, event_user_id);
ALTER TABLE dwh_content_units_measures ADD COLUMN IF NOT EXISTS unique_users_watching_now_count bigint;
UPDATE dwh_content_units_measures SET unique_users_watching_now_count = unique_users_last10min_count;
ALTER TABLE dwh_content_units_measures DROP COLUMN IF EXISTS unique_users_last10min_count;

