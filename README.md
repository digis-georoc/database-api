# database-api

API to access the GEOROC database in the new ODM2 schema

## Structure

The database resources are available at `/api/v1/queries`.
You can get all resources of a type by GETing `/resource`, optionally providing pagination parameters `/resource?limit=10&offset=30`, where `limit` is the pagesize and `offset` is the pagenumber times the pagesize.
To get a specific resource by its identifier, GET `/resource/:identifier`.

## Documentation

To view the api documentation, open the route `/api/v1/docs/index.html` in your browser.

The api documentation is generated with [swagger](https://github.com/swaggo/swag).
For installation guides see [the documentation](https://github.com/swaggo/swag#getting-started). You may need to add the `GOPATH` to your `PATH` variable to be able to exectue `swag` from your commandline: `export PATH=$(go env GOPATH)/bin:$PATH`.

To generate the documentation files under **docs/**, execute

`swag fmt && swag init -d pkg/api/,pkg/api/handler/,pkg/model/ -g api.go`
