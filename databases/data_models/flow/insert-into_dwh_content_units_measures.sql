truncate table dwh_content_units_measures;
insert into dwh_content_units_measures

(

select 

event_unit_uid, 
count(event_end_id_max) as all_events_count,
count(distinct case when event_end_minute>= NOW() - (15* interval '1 minute') then event_user_id else null end) as unique_users_last10min_count,
count(distinct event_user_id) as unique_users_count,
NOW() as dwh_update_datetime

from public.dwh_fact_play_units_by_minutes f 

group by event_unit_uid  


)
