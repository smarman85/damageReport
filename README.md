#pagerDuty

```bash
# Use the link as a ref. to set the pagerDuty api token and add this variable to your bash_profile
export PAGER_DUTY_TOKEN="https://support.pagerduty.com/docs/using-the-api"
```

##run
```
$ go run app.go
```
# damageReport

## Dockerize this:
```
#Build:
$ docker build -t damagereport .

#Run:
$ docker run -e PAGER_DUTY_TOKEN="$(echo $PAGER_DUTY_TOKEN)" --rm damagereport
```
