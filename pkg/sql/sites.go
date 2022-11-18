package sql

const SitesQuery = `
SELECT * FROM odm2.sites
`

const SitesByCoordsQuery = `
SELECT * FROM odm2.sites s
WHERE s.latitude >= $1 AND s.latitude <= $2 AND s.longitude >= $3 AND s.longitude <= $4
`
