// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package sql

const SamplingTechniquesQuery = `
select distinct a.annotationtext as name
from odm2.annotations a 
where a.annotationcode = 'g_samples.samp_technique'
`

const MaterialsQuery = `
select distinct a.annotationtext as name
from odm2.annotations a 
where a.annotationcode = 'g_batches.material'
`

const InclusionTypesQuery = `
select distinct a.annotationtext as name
from odm2.annotations a 
where a.annotationcode = 'g_inclusions.inclusion_type'
`
