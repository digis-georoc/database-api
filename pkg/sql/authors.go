// SPDX-FileCopyrightText: 2024 DIGIS Project Group
//
// SPDX-License-Identifier: BSD-3-Clause

package sql

const AuthorsQuery = `
SELECT * 
FROM odm2.people p
`

const AuthorByIDQuery = `
SELECT * 
FROM odm2.people p
WHERE p.personID = $1
`
