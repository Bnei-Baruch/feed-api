truncate table dwh_dim_content_units;
insert into dwh_dim_content_units
(
  /** all content units from mdb **/	
	select * from dblink('mdb_conn',
    'select 
      cu.id as content_unit_id,
      cu.uid as content_unit_uid,
      cu.created_at as content_unit_created_at,
      properties->>''duration'' as content_unit_duration,
      cu.type_id as content_unit_type_id,
      ct.name as content_unit_type_name,
      cn.name as content_unit_name,
      cn.language as content_unit_language
    from content_units cu 
    join content_types ct on (cu.type_id=ct.id)
    join content_unit_I18n cn on (cu.id=cn.content_unit_id)') 
	as content_unit(
    content_unit_id integer,
    content_unit_uid character(8),
    content_unit_created_at timestamp with time zone,
    content_unit_duration bigint,
    content_unit_type_id bigint,
    content_unit_type_name character varying, 
    content_unit_name character varying,
    content_unit_language character varying
  )
)
