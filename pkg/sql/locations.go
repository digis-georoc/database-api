package sql

const FirstLevelLocationNamesQuery = `
select distinct sg.locationname as name
from odm2.sitegeolocations sg
left join odm2.geolocations g on g.geolocationid = sg.geolocationid
where right(g.locationhierarchy::varchar, 3) = '100'
`

const SecondLevelLocationNamesQuery = `
select distinct secondlevelloc.locationname as name
from odm2.sites s 
left join
(
	select sg.samplingfeatureid, sg.locationname
	from odm2.sitegeolocations sg
	left join odm2.geolocations g on g.geolocationid = sg.geolocationid
	where right(g.locationhierarchy::varchar, 3) = '100' -- first level
) firstlevelloc on firstlevelloc.samplingfeatureid = s.samplingfeatureid 
left join
(
	select sg.samplingfeatureid, sg.locationname
	from odm2.sitegeolocations sg
	left join odm2.geolocations g on g.geolocationid = sg.geolocationid
	where right(g.locationhierarchy::varchar, 3) = '200' --second level
) secondlevelloc on secondlevelloc.samplingfeatureid = s.samplingfeatureid 
where firstlevelloc.locationname = $1
`

const ThirdLevelLocationNamesQuery = `
select distinct thirdlevelloc.locationname as name
from odm2.sites s 
left join
(
	select sg.samplingfeatureid, sg.locationname
	from odm2.sitegeolocations sg
	left join odm2.geolocations g on g.geolocationid = sg.geolocationid
	where right(g.locationhierarchy::varchar, 3) = '100' -- first level
) firstlevelloc on firstlevelloc.samplingfeatureid = s.samplingfeatureid 
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
where firstlevelloc.locationname = $1
and secondlevelloc.locationname = $2
`
