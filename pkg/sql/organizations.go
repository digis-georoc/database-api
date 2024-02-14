// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package sql

const GetOrganizationNamesQuery = `
select distinct o.organizationname as name
from odm2.organizations o
`
