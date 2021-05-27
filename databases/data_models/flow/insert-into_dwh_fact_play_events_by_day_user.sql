insert into dwh_fact_play_events_by_day_user

(
  
select 
  a.event_user_id,
  a.event_user_agent_type,
  a.event_language, 
  cast(event_start_date as date) as event_start_date, event_unit_uid, 

  content_unit_name,
  content_unit_type_name,
  content_unit_created_at,
  content_unit_language,
  round(sum(content_unit_duration)) as content_unit_duration,

  round(sum(event_duration_sec)) as event_duration_sec,
  round(sum(event_current_duration_sec)) as event_current_duration_sec,

  round(case when sum(content_unit_duration) is null or sum(content_unit_duration)=0 or sum(event_duration_sec)>sum(content_unit_duration) then 0
      else sum(event_current_duration_sec)*100/sum(content_unit_duration) end )  as event_of_unit_duration_percent,

  count(*) as event_count,
  NOW() as dwh_update_datetime

  from (
  /** all play events only - play join stop **/
  select * from dblink('chronicles_conn', '
    select 
      a.user_id as event_user_id, 
      case when a.user_agent like ''%Mobile%'' then ''mobile'' else ''stationary'' end as event_user_agent_type,
      a.data->>''language'' as event_language,
      a.data->>''unit_uid'' as event_unit_uid,
      a.created_at as event_start_date,
      b.created_at as event_end_date,
      (DATE_PART(''hour'', b.created_at-a.created_at::time) * 60 +
       DATE_PART(''minute'', b.created_at-a.created_at::time)) * 60 +
       DATE_PART(''second'', b.created_at-a.created_at::time) as event_duration_sec,
      a.data->>''current_time'' as event_current_start_sec,
      b.data->>''current_time'' as event_current_end_sec,
      case when cast(b.data->>''current_time'' as float)>cast(a.data->>''current_time'' as float)
        then cast(b.data->>''current_time'' as float)-cast(a.data->>''current_time'' as float)
        else 0 end as event_current_duration_sec
    from entries a 
    left join entries b on (a.client_event_id=b.client_flow_id and b.client_event_type=''player-stop'')
    
    where 
    a.client_event_type=''player-play''
    /** for incremental loading **/ 
    /** and cast(a.created_at as date) between ''2021-05-01'' and ''2021-05-20'' **/
		and a.id >= ''$prev-read-id''
  ') as a(
    event_user_id character varying(64),
    event_user_agent_type text,
    event_language character varying,
    event_unit_uid character varying,
    event_start_date timestamp with time zone,
    event_end_date timestamp with time zone,
    event_duration_sec double precision,
    event_current_start_sec float,
    event_current_end_sec float,
    event_current_duration_sec double precision
   )   
  ) a
	
	
/** 06/07/2021 updated to join instead of left join because of the units that are being deleted from mdb **/
	
join

/** all content units from mdb **/

  dwh_dim_content_units c on (a.event_unit_uid=c.content_unit_uid
              /** remove the filter when there will be a language field to join to **/
              and c.content_unit_language='he'
              )

  group by cast(event_start_date as date), event_unit_uid, content_unit_name, content_unit_type_name, content_unit_created_at, content_unit_language
  , event_language, event_user_id, a.event_user_agent_type

)
