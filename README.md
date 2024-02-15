<!--
SPDX-FileCopyrightText: 2024 DIGIS Project Group

SPDX-License-Identifier: BSD-3-Clause
-->

# Database-API

API to access the GEOROC2.0 database in the new ODM2-based schema.

Current version: **0.4.1**

This api is currently in a testing phase.

## Usage

For a basic liveness check, the route `/api/v1/ping` can be used. If the api is up and running, it will answer with http status 200 and the string "pong".
Documentation of the available resources and routes can be found at `/api/v1/docs/index.html`.

The database resources are available under `/api/v1/queries/`. All routes with the `/queries/`-prefix are secured via access key.
A valid access key has to be provided in the custom header `DIGIS-API-ACCESSKEY` for each request to a secured route.
To obtain an access key, refer to the [Get Access](#get-access) part of this README.

Most routes support basic pagination with the query-parameters `limit` and `offset`.

## Get Access

To access the api, a personal access token is needed.
Access tokens are generated on demand.
Please note that the api is currently in a testing phase, so bugs or service outages can happen.
If you want to provide feedback and participate in testing the api, please [contact us](digis-info@uni-goettingen.de).

## License

[![REUSE status](https://api.reuse.software/badge/github.com/digis-georoc/database-api)](https://api.reuse.software/info/github.com/digis-georoc/database-api)

The source code is licensed under the BSD-3-Clause license; generated documentation files are licensed under CC0-1.0.
The licensing of this repository is [reuse compliant](https://reuse.software/).
License headers were added with the reuse-helper-tool:
`docker run --rm --volume ./:/data fsfe/reuse annotate --copyright="DIGIS Project Group" --license="BSD-3-Clause" --recursive ./*`

## Operation

### Setup

For running the api, a Dockerfile is provided. No build-time parameters are needed.
The api needs database configuration to start. The database configuration must be provided in a file named `database-config.txt` with the contents as follows:

```json

{
    "DB_USER":"",
    "DB_PASSWORD":"",
    "SSH_USER":"",
    "SSH_PASSWORD":""
}

```

The ssh config is optional.

To be able to access the secured routes, access keys must be provided as a second configuration file named `accesskeys.txt` with contents formatted as follows:

```json

{
    "<KEY_NAME1>": "<KEY1>",
    "<KEY_NAME2>": "<KEY2>"
}

```

The configuration files must be mounted in the container under the path `/vault/secrets/` (e.g. `docker run -v <absolute-path-to-config-files-on-host>:/vault/secrets/ digis-api`).

### Update Documentation

The api documentation is generated with [swagger](https://github.com/swaggo/swag).
For installation guides see [the documentation](https://github.com/swaggo/swag#getting-started). You may need to add the `GOPATH` to your `PATH` variable to be able to exectue `swag` from your commandline: `export PATH=$(go env GOPATH)/bin:$PATH`.

To generate the documentation files under **docs/**, execute

`swag fmt && swag init -d pkg/api/,pkg/api/handler/,pkg/model/ -g api.go`

### Update Version

To update the api version, change the version number in the following places:

- The start of this README
- handler.go/version method
- api.go swaggo comment
- push a new tag with the new version number to the repositories main branch (attention: triggers a staging (test) release!)
