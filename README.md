# MyDB Server

Ready to use database server include mydb package. It allows you to operate with simple commands with the Redis protocol (Thanks Antirez).

- Redcon init (all redis clients supported)
- Use my-db package (support boltdb and sniper db) (https://github.com/uretgec/mydb)
- Dockerfiles include
- Docker-compose file includes
- Build commands init (Go build, File generator)
- One service multiple db server support
- Go file generator commands include (go to installation section)

NOTES:
> This is not go package.

## Commands
```
// Main Commands
quit
ping -> string
set bucketName id "json_stringfy_data or simple value" -> string
get bucketName id -> string
mget bucketName id1 id2 id3 ... -> map[string]string
list bucketName cursor perpage -> map[string]string
prevlist bucketName cursor perpage -> map[string]string (reverse list command)
exists bucketName key -> bool
vexists bucketName value -> bool
del bucketName id -> int64
bstats bucketName -> int
backup path/ filename -> OK
restore path/ filename -> OK

// Broadcast Commands
publish channel message
subscribe channel
psubscribe channel
```

## Installation

1 Open mycli/.generator.dev.yaml (sample generator configration file) and create yours yaml file to own configrations

```
service-name: mystorage // your main service name. Recommended not change :)
databases: // your database name list (one or multiple name supports)
  - devstore
  - mystore

// every database name configration options
devstore:
  sync-interval: 10
  redis-addr: localhost
  redis-port: 6379
  db-storage: boltdb # boltdb|sniper
  db-name: devstore
  db-display-name: DevStore
  db-path: ../data/db/devstore/
  db-read-only: false
  container-name: devstorage
  bucket-list:
    - "urls"
    - "items"
  index-list:
    - "urls"
    - "items"

mydbstore:
  sync-interval: 10
  redis-addr: localhost
  redis-port: 6380
  db-storage: sniper # boltdb|sniper
  db-name: mydbstore
  db-display-name: MyDbStore
  db-path: ../data/db/mydbstore/
  db-read-only: false
  container-name: mydbstorage
  bucket-list:
    - "urls"
    - "items"
    - "options"
  index-list:
    - "urls"
    - "items"

```

2 Run the file generator command (change .generator.dev.yaml to yours)

```
go run . generator --config .generator.dev.yaml >> ../change.log
```

- env files (configration file)
- service files (unit file)
- storage file
- store files
- docker-compose file

has to be created.

3 Quick run the storage service

```
cd mystorage && go run . --env-file .devstore.env
```

or try docker-compose option

First run builder command
```
export SERVICE_BUILD=$(date '+%Y%m%d%H%M') && export SERVICE_COMMIT_ID=$(git describe --always) &&  docker-compose -f docker-compose-builder.yml build --compress --progress plain
```

after run docker-compose for build containers

```
export SERVICE_BUILD=$(date '+%Y%m%d%H%M') && export SERVICE_COMMIT_ID=$(git describe --always) &&  docker-compose -f docker-compose.yml build --compress --progress plain
```

finally up all containers

```
docker-compose up
```

close and clear all running containers
```
docker-compose down
```

## Build

Run this command to build go binaries and voala (change .build.dev.yaml to yours)

```
go run . build go --config .build.dev.yaml >> ../build.log
```

If has no error, go to build folder.

Ready to use. Good luck.

## TODO
- Add test files
- Add new examples

## Links

Bbolt (https://github.com/etcd-io/bbolt)

Sniper (https://github.com/recoilme/sniper)