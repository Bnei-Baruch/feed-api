create TEMP table IF NOT EXISTS all_events_count_temp (event_unit_uid text, all_events_count bigint);
truncate table all_events_count_temp;
insert into all_events_count_temp (
select
  event_unit_uid, count(event_end_id_max) as all_events_count
from dwh_fact_play_units_by_minutes f 
group by event_unit_uid
);


create TEMP table IF NOT EXISTS unique_users_count_temp (event_unit_uid text,unique_users_count bigint);
truncate table unique_users_count_temp;
insert into unique_users_count_temp (
select event_unit_uid, sum(one) as unique_users_count
from (
		select distinct event_unit_uid, event_user_id, 1 as one from dwh_fact_play_units_by_minutes
	) a group by event_unit_uid
);


create TEMP table IF NOT EXISTS watching_now_temp (event_unit_uid text, unique_users_watching_now_count bigint);
truncate table watching_now_temp;
insert into watching_now_temp (
select  
event_unit_uid,
count(distinct last180min) as unique_users_watching_now_count
from (
  select
  event_unit_uid,
  case when event_end_minute>= NOW() - (180* interval '1 minute') then event_user_id else null end as last180min
  from dwh_fact_play_units_by_minutes f 
	 )a 
where last180min is not null
group by event_unit_uid
having count(distinct last180min)>=100
);


truncate table dwh_content_units_measures;
insert into dwh_content_units_measures (
	
  select a.event_unit_uid, all_events_count, unique_users_watching_now_count, unique_users_count, NOW() as dwh_update_datetime 
  from all_events_count_temp a 
	left join unique_users_count_temp u on (a.event_unit_uid=u.event_unit_uid)
	left join watching_now_temp w on (a.event_unit_uid=w.event_unit_uid)
);


drop table all_events_count_temp;
drop table unique_users_count_temp;
drop table watching_now_temp;
