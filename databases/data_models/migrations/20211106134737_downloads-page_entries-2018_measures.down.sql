DROP TABLE IF EXISTS dwh_content_units_measures_2018_2021;
DROP TABLE IF EXISTS dwh_fact_download_units_by_minutes;
DROP TABLE IF EXISTS dwh_fact_page_enter_units_by_minutes;

ALTER TABLE dwh_content_units_measures
	DROP COLUMN IF EXISTS page_enter_events_count bigint, 
	DROP COLUMN IF EXISTS download_events_count bigint,
	DROP COLUMN IF EXISTS page_enter_events_count_2018 bigint;

