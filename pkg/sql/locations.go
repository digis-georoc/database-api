package sql

const LocationNamesByLevelQuery = `
select distinct sg.locationname as name
from odm2.sitegeolocations sg
left join odm2.geolocations g on g.geolocationid = sg.geolocationid
where right(g.locationhierarchy::varchar, 3) = $1
`
