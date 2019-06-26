# small test
in a first terminal
```
UPSTREAM=https://play.dhis2.org/2.32.0 go run main.go
```
then in a second terminal

```
curl -v http://127.0.0.1:1330/api/organisationUnits -u admin:district
curl -v http://127.0.0.1:1330/api/organisationUnits?filter=name:ilike:oil -u admin:district
```

problematic

```
curl  -L -v http://127.0.0.1:1330
```

# current pending request via stats endpoint

only ongoing request


start the reverse proxy
```
   UPSTREAM=https://play.dhis2.org go run main.go
```
watch the stats endpoint
```
   watch -n 1 'curl -s "http://127.0.0.1:3000/stats" -u admin:district | jq .'
```
do from time to time request via the proxy 
```
   time curl "http://127.0.0.1:3000/2.32.0/api/organisationUnits?paging=false&fields=:all" -u admin:district -o /dev/null
```
the watch screen should from time to time display some content

in browser http://127.0.0.1:3000/2.32.0/dhis-web-dashboard/#/
login with admin district
then issue a lot CTRL-SHIFT-R

