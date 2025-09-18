#!/bin/bash

go build -o bin/actors ./cmd/actors/
go build -o bin/map_reduce ./cmd/map_reduce/
go build -o bin/monolithic ./cmd/monolithic/
go build -o bin/persistent_tables ./cmd/persistent_tables/
go build -o bin/pipeline ./cmd/pipeline/
go build -o bin/quarantine ./cmd/quarantine/
go build -o bin/things ./cmd/things/