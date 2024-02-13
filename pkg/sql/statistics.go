// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package sql

const CountCitationsQuery = `
select count(c.citationid) as numCitations from odm2.citations c
`

const CountAnalysesQuery = `
select count(s.samplingfeatureid) as numAnalyses from odm2.samplingfeatures s
where s.samplingfeaturedescription = 'Batch'
`

const CountSamplesQuery = `
select count(s.samplingfeatureid) as numSamples from odm2.samplingfeatures s
where s.samplingfeaturedescription = 'Sample'
`

const CountResultsQuery = `
select count(r.resultid) as numResults from odm2.results r
`
