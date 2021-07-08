DROP TABLE IF EXISTS dwh_content_units_measures;
CREATE TABLE IF NOT EXISTS dwh_content_units_measures
(
    event_unit_uid               text COLLATE "default" NOT NULL,
    content_unit_type_name       character varying COLLATE "default",
    all_events_count             bigint,
    unique_users_curday_count    bigint,
    dwh_update_datetime          timestamp with time zone,

    CONSTRAINT dwh_content_units_measures_pkey PRIMARY KEY (event_unit_uid)
);

CREATE SEQUENCE IF NOT EXISTS dwh_fact_play_events_by_day_user_id_seq;

DROP TABLE IF EXISTS dwh_fact_play_events_by_day_user;
CREATE TABLE IF NOT EXISTS dwh_fact_play_events_by_day_user
(
    event_stop_id_max 			   character(27) COLLATE pg_catalog."POSIX" NOT NULL,
    event_user_id                  character varying COLLATE "default",
    event_user_agent_type          character varying COLLATE "default",
    event_language                 character varying COLLATE "default",
    event_end_date                 timestamp with time zone,
    event_unit_uid                 text COLLATE "default",
    content_unit_name              character varying COLLATE "default",
    content_unit_type_name         character varying COLLATE "default",
    content_unit_created_at        timestamp with time zone,
    content_unit_language          character varying COLLATE "default",
    content_unit_duration          numeric,
    event_duration_sec             double precision,
    event_current_duration_sec     double precision,
    event_of_unit_duration_percent double precision,
    event_count                    bigint,
    dwh_update_datetime            timestamp with time zone,
    id                             integer NOT NULL DEFAULT nextval('dwh_fact_play_events_by_day_user_id_seq'::regclass),

    CONSTRAINT dwh_fact_play_events_by_day_user_pkey PRIMARY KEY (id)
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
