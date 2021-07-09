
insert into dwh_fact_play_units_by_minutes

(
  
select 
  event_unit_uid, 
  event_user_id,
  event_language, 
  event_user_agent_type,
  date_trunc('minute', event_end_timestamp) as event_end_minute,
  max(event_end_id) as event_end_id_max,
  count(*) as event_count,
  round(sum(event_duration_sec)) as event_duration_sec_sum,
  NOW() as dwh_update_datetime
 
  from (
  /** all play events only - play join stop **/
  select * from dblink('chronicles_conn', '
    select
    a.data->>''unit_uid'' as event_unit_uid,					   
    a.user_id as event_user_id, 
    a.data->>''language'' as event_language,
	case when a.user_agent like ''%Mobile%'' then ''mobile'' else ''stationary'' end as event_user_agent_type,							   
    b.created_at as event_end_timestamp,
	b.id as event_end_id,
    (DATE_PART(''hour'', b.created_at-a.created_at::time) * 60 +
     DATE_PART(''minute'', b.created_at-a.created_at::time)) * 60 +
     DATE_PART(''second'', b.created_at-a.created_at::time) as event_duration_sec

	from entries a join entries b on (a.client_event_id=b.client_flow_id and b.client_event_type=''player-stop'')
    where 
	a.client_event_type=''player-play''
    /** for incremental loading **/ 
	--and b.id >= ''$prev-read-id''
  ') as a(
    event_unit_uid character varying,
    event_user_id character varying(64),	  
    event_language character varying,	  
    event_user_agent_type text,
    event_end_timestamp timestamp with time zone,	  
    event_end_id character(27) COLLATE pg_catalog."POSIX", 
    event_duration_sec double precision
         )   
  ) c
	
group by 
  event_unit_uid, 
  event_user_id,
  event_language, 
  event_user_agent_type,
  date_trunc('minute', event_end_timestamp)


)

