DROP TABLE IF EXISTS dwh_content_units_measures;
CREATE TABLE IF NOT EXISTS dwh_content_units_measures
(
    event_unit_uid               text COLLATE "default" NOT NULL,
    all_events_count             bigint,
    unique_users_last10min_count bigint,
    unique_users_count           bigint,
    dwh_update_datetime          timestamp with time zone,

    CONSTRAINT dwh_content_units_measures_pkey PRIMARY KEY (event_unit_uid)
);


DROP TABLE IF EXISTS dwh_fact_play_units_by_minutes;
CREATE TABLE IF NOT EXISTS dwh_fact_play_units_by_minutes
(
  event_unit_uid         text COLLATE "default", 
  event_user_id          character varying COLLATE "default",
  event_language         character varying COLLATE "default", 
  event_user_agent_type  character varying COLLATE "default",
  event_end_minute       timestamp with time zone,
  event_end_id_max       character(27) COLLATE "POSIX" NOT NULL,
  event_count            bigint,
  event_duration_sec_sum double precision,
  dwh_update_datetime    timestamp with time zone,

  CONSTRAINT dwh_fact_play_units_by_minutes_pkey PRIMARY KEY (event_end_id_max)
)
WITH (
    OIDS = FALSE
);

DROP TABLE IF EXISTS dwh_dim_content_units;
CREATE TABLE IF NOT EXISTS dwh_dim_content_units
(
    content_unit_id         integer,
    content_unit_uid        character(8) COLLATE "default" NOT NULL,
    content_unit_created_at timestamp with time zone,
    content_unit_duration   bigint,

    content_unit_type_id    bigint,
    content_unit_type_name  character varying COLLATE "default",
    content_unit_name       character varying COLLATE "default",
    content_unit_language   character varying COLLATE "default" NOT NULL,

    CONSTRAINT dwh_dim_content_units_pkey PRIMARY KEY (content_unit_uid, content_unit_language)
);

