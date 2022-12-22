package sql

const GeoJSONQuery = `
select
(s.latitude, s.longitude) as coordinates,
count(distinct s.samplingfeatureid) as num_samplingfeatureids,
array_agg(distinct s.samplingfeatureid) as samplingfeatureids,
array_agg(distinct s.latitude) as lat ,
array_agg(distinct s.longitude) as long,
array_agg(distinct s.setting) as setting,
array_agg(distinct toplevelloc.locationname) as loc1,
array_agg(distinct secondlevelloc.locationname) as loc2,
array_agg(distinct thirdlevelloc.locationname) as loc3,
array_agg(distinct s.sitedescription) as land_or_sea
from odm2.sites s
left join 
(
	select sg.samplingfeatureid, sg.locationname
	from odm2.sitegeolocations sg
	left join odm2.geolocations g on g.geolocationid = sg.geolocationid
	where right(g.locationhierarchy::varchar, 3) = '100' --toplevel
) toplevelloc on toplevelloc.samplingfeatureid = s.samplingfeatureid 
left join
(
	select sg.samplingfeatureid, sg.locationname
	from odm2.sitegeolocations sg
	left join odm2.geolocations g on g.geolocationid = sg.geolocationid
	where right(g.locationhierarchy::varchar, 3) = '200' --second level
) secondlevelloc on secondlevelloc.samplingfeatureid = s.samplingfeatureid 
left join
(
	select sg.samplingfeatureid, sg.locationname
	from odm2.sitegeolocations sg
	left join odm2.geolocations g on g.geolocationid = sg.geolocationid
	where right(g.locationhierarchy::varchar, 3) = '300' --third level
) thirdlevelloc on thirdlevelloc.samplingfeatureid = s.samplingfeatureid
group by (s.latitude, s.longitude)
`
