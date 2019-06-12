in a first terminal

UPSTREAM=https://play.dhis2.org/2.32.0 go run main.go

then in a second terminal

curl -v http://127.0.0.1:1330/api/organisationUnits -u admin:district

curl -v http://127.0.0.1:1330/api/organisationUnits?filter=name:ilike:oil -u admin:district


problematic

curl  -L -v http://127.0.0.1:1330