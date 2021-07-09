insert into dwh_fact_play_units_by_minutes
(
  
select 
  event_unit_uid, 
  content_unit_type_name,
  event_end_minute, 
  max(a.event_end_id) as event_end_id_max,
  count(*) as event_count,
  NOW() as dwh_update_datetime

  from (
  /** all play events only - play join stop **/
  select * from dblink('chronicles_conn', '
    select b.id as event_end_id,
    a.data->>''unit_uid'' as event_unit_uid,
    date_trunc(''minute'', b.created_at) as event_end_minute
    from entries a 
    join entries b on (a.client_event_id=b.client_flow_id and b.client_event_type=''player-stop'')
    where 
	a.client_event_type=''player-play''
    /** for incremental loading **/ 
	and b.id >= ''$minutes-prev-read-id''
  ') as a(
    event_end_id character(27) COLLATE pg_catalog."POSIX", 
    event_unit_uid character varying,
    event_end_minute timestamp with time zone
   )   
  ) a
	
	
join

/** all content units from mdb **/

  dwh_dim_content_units c on (a.event_unit_uid=c.content_unit_uid
              /** remove the filter when there will be a language field to join to **/
              and c.content_unit_language='he'
              )

  group by event_unit_uid, event_end_minute, content_unit_type_name

	order by 5 desc
	
)



