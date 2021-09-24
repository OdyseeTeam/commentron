# Commentron 

This is the commenting system for odysee.com. The system uses a JSON RPC server.
 
## Run it locally

### Goland

![GoLand configuration](goland-config.png)

Make sure you setup the configuration to use env files. 

`MYSQL_DSN_RO="lbry-ro:lbry@tcp(localhost:3306)/commentron"`

`MYSQL_DSN_RW="lbry-rw:lbry@tcp(localhost:3306)/commentron"`

`SDK_URL="https://api.lbry.tv/api/v1/proxy"`

`SOCKETY_TOKEN="sockety_token" #If you want to integrate directly with sockety locally`

I put the `IS_TEST=true` in the configuration but it could be in the `.env` file. 

### MySQL 5.7

Install MySQL 5.7 and create a database named `commentron` then adjust the DSN
in the `.env` file locally with the username and password. 

Then once mysql is running you can hit play. 

## Key packages

These are the key packages where most work is completed.

`commentapi` -  This contains all the client API information 
`server/service` - Contains all the API implementations
`migration` - Contains the migrations to the database
`http_requests` - This has http requests for testing different APIs. Best used as templates





