truncate table dwh_content_units_measures;
insert into dwh_content_units_measures

(

select 

event_unit_uid, 
content_unit_type_name, 
count(event_end_id_max) as all_events_count,
count(distinct case when event_end_minute>= event_end_minute_max - (10* interval '1 minute') then event_user_id else null end) as unique_users_last10min_count,
NOW() as dwh_update_datetime

from public.dwh_fact_play_units_by_minutes f join dwh_dim_content_units d on (f.event_unit_uid=d.content_unit_uid)
cross join (select max(event_end_minute) as event_end_minute_max from dwh_fact_play_units_by_minutes) fm

group by event_unit_uid, content_unit_type_name

)
