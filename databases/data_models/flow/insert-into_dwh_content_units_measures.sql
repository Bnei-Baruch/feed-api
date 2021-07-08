truncate table dwh_content_units_measures;
insert into dwh_content_units_measures
(
select event_unit_uid, content_unit_type_name, sum(event_count) as all_events_count,
count (distinct case when event_end_date = cast(NOW() as date) then event_user_id else null end) as unique_users_curday_count,
max(NOW() AT TIME ZONE '-3:00') as dwh_update_datetime
from dwh_fact_play_events_by_day_user
where content_unit_type_name IS NOT NULL
group by event_unit_uid, content_unit_type_name
)
  
