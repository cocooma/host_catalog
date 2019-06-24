
Introduction to the Application
------------------------------

The application is a gRPC REST API server which expose endpoints to allow accessing and manipulating `host data`. <br/> 
The operations that our endpoint will allow are:
* Creating a new host.
* Updating an existing host.
* Deleting an existing host.
* Fetching an existing host.
* Fetching a list of hots.

### API Specification
* Application health check GET /health,
* Create a new host in response to a valid POST request at /host,
* Update a host in response to a valid PUT request at /host/{id},
* Delete a host in response to a valid DELETE request at /host/{id},
* Fetch a host in response to a valid GET request at /host/{id}, and
* Fetch a list of hosts in response to a valid GET request at /hosts.

### Prerequisites

* Docker
* Golang
* PostgreSQL

### Dependencies

* mux – The Gorilla Mux router.
* pq – The PostgreSQL driver

### Application Scaffolding

```
┌── app.go
├── main.go
├── main_test.go
└── model.go
```

### Database Structure
    
* id – the primary key in this table,
* name – the name of the host.
* ip – the host ipv4 address.

```
CREATE TABLE hosts
(
    id SERIAL,
    name TEXT UNIQUE NOT NULL,
    ip VARCHAR(15) NOT NULL,
    CONSTRAINT hosts_pkey PRIMARY KEY (id)
)
```

##### In order to simplify the build process Makefile is used as a wrapper 
### Test the Application

```
make test
```

The test process make use of the upstream ` golang:1.12` and `circleci/postgres:latest `docker containers.<br/>
Technicality only `docker` and `make` should be available on the build environment.  

### Build the Application

```
make
```

##### The build process contains 3 stesp:
 * Runes Unit tests.
 * Compile build application binary.
 * Packages the compiled binary to the latest alpine docker container.

The final product is a docker container called `host_catalog:latest`

### Run up the appliction with docker-compose

```
docker-compose up -d
```

It will spin up the application container and the Postgres database.

### Application run time config parameters
The following environment variables are required by the application. 

```
APP_DB_HOST
APP_DB_USERNAME
APP_DB_PASSWORD
APP_DB_NAME
```

### The application in action
#### Create a host entry 

```
# curl  -X POST --data '{"name":"host_1","ip":"123.123.123.123"}' http://127.0.0.1:8000/host
{"id":1,"name":"host_1","ip":"123.123.123.123"}%

```
#### Update host details

```
# curl  -X PUT --data '{"name":"host_2","ip":"221.221.221.221"}' http://127.0.0.1:8000/host/2
{"id":2,"name":"host_2","ip":"221.221.221.221"}%
```

####  Fetch a host details

```
# curl  http://127.0.0.1:8000/host/2
{"id":2,"name":"host_2","ip":"221.221.221.221"}%
```
####  Fetch all host details

```
# curl -s  http://127.0.0.1:8000/hosts | jq .
[
  {
    "id": 1,
    "name": "host_1",
    "ip": "123.123.123.123"
  },
  {
    "id": 2,
    "name": "host_2",
    "ip": "221.221.221.221"
  }
]
```
#### Delete host entry

```
# curl -X DELETE http://127.0.0.1:8000/host/1
{"result":"success"}%
```
