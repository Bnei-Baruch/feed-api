DROP TABLE IF EXISTS dwh_content_units_measures_2018_2021;
CREATE TABLE IF NOT EXISTS dwh_content_units_measures_2018_2021
(
	event_unit_uid      text COLLATE "default" NOT NULL,
	page_enter_count    bigint,
	dwh_update_datetime timestamp with time zone,

	CONSTRAINT dwh_content_units_measures_2018_2021_pkey PRIMARY KEY (event_unit_uid)
);

DROP TABLE IF EXISTS dwh_fact_download_units_by_minutes;
CREATE TABLE IF NOT EXISTS dwh_fact_download_units_by_minutes
(
  event_unit_uid        text COLLATE "default",
  event_user_id         character varying COLLATE "default",
  event_language        character varying COLLATE "default",
  event_user_agent_type character varying COLLATE "default",
  event_minute          timestamp with time zone,
  event_id_max          character(27) COLLATE "POSIX" NOT NULL,
  event_count           bigint,
  dwh_update_datetime   timestamp with time zone
);

DROP TABLE IF EXISTS dwh_fact_page_enter_units_by_minutes;
CREATE TABLE IF NOT EXISTS dwh_fact_page_enter_units_by_minutes
(
  event_unit_uid        text COLLATE "default",
  event_user_id         character varying COLLATE "default",
  event_language        character varying COLLATE "default",
  event_user_agent_type character varying COLLATE "default",
  event_minute          timestamp with time zone,
	event_id_max          character(27) COLLATE "POSIX" NOT NULL,
  event_count           bigint,
  dwh_update_datetime   timestamp with time zone
);

ALTER TABLE dwh_content_units_measures
	ADD COLUMN page_enter_events_count bigint, 
	ADD COLUMN download_events_count bigint,
	ADD COLUMN page_enter_events_count_2018 bigint;

