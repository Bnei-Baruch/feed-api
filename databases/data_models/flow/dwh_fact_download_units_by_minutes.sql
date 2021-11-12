insert into dwh_fact_download_units_by_minutes

(
  
select 
  unit_uid as event_unit_uid, 
  event_user_id,
  event_language, 
  event_user_agent_type,
  date_trunc('minute', event_timestamp) as event_minute,
  max(event_id) as event_id_max,
  count(*) as event_count,
  NOW() as dwh_update_datetime
 
  from (
  /** all unit-page-enter events only **/

	  select * from dblink('chronicles_conn', '
    select
    data->>''uid'' as event_file_uid,					   
    user_id as event_user_id, 
    data->>''content_language'' as event_language,
	case when user_agent like ''%Mobile%'' then ''mobile'' else ''stationary'' end as event_user_agent_type,							   
    created_at as event_timestamp,
    id as event_id
	from entries
    where client_event_type=''download''
    /** for incremental loading **/ 
	and id > ''$download-minutes-prev-read-id''
  ') as a(
    event_file_uid character varying,
    event_user_id character varying(64),	  
    event_language character varying,	  
    event_user_agent_type text,
    event_timestamp timestamp with time zone,
    event_id character(27) COLLATE "POSIX"
          )   
	  			
  ) c
	
		join 	(
					select * from dblink('mdb_conn',' select distinct f.uid as file_uid, u.uid as unit_uid  
						 from files f join content_units u on (f.content_unit_id=u.id)') as b
						  (   file_uid  character varying,
							  unit_uid character varying)
					) b on (c.event_file_uid=b.file_uid)	
	
group by 
  unit_uid, 
  event_user_id,
  event_language, 
  event_user_agent_type,
  date_trunc('minute', event_timestamp)



)
