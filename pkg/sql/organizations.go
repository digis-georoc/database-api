package sql

const GetOrganizationNamesQuery = `
select distinct o.organizationname as name
from odm2.organizations o
`
