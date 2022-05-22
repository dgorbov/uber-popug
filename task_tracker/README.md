# Task tracker service
- Provides REST API to manage tasks in the system
- Requires keycloak as auth service
- All APIs are protected and requires header `Authorization: Bearer {access_token}`

## Required env variables
- `PORT`
- `KEYCLOAK_URL`
- `KAFKA_TOPIC_NEW_USER`
- `KAFKA_BROKERS_URL` coma separated values: `localhost:9091,localhost:9092,localhost:9093`.

## Tests
- `go test -v ./...`

## To run locally
- `docker-compose up`
- keycloak container will spawn with test realm from file `test/uber-popug-realm.json`
- You can enter keycloak admin panel at: `http://localhost:8080/admin/` with admin:admin
- `export PORT=8081 KEYCLOAK_URL=http://localhost:8080 KAFKA_TOPIC_NEW_USER=new-user KAFKA_BROKERS_URL=localhost:9091`
- `go build task_tracker`
- `go run task_tracker`

### Existing users in keycloak (username:password)
- popug:popug, has role `popug`
- admin:admin, has role `admin`
- manager:manager, has role `manager`

### To get access_token from keycloak
````
curl --location --request POST 'localhost:8080/realms/uber-popug/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=password' \
--data-urlencode 'scope=openid' \
--data-urlencode 'client_id=auth-proxy' \
--data-urlencode 'client_secret=yGSDYoV1XMEFIf6XMxpxFqZivMpJm70d' \
--data-urlencode 'password=popug' \
--data-urlencode 'username=popug'
````
### To get all task assigned to logged-in user (if any)
````
curl --location --request GET 'localhost:8081/tasks/my' \
--header 'Authorization: Bearer {access_token}'
````
### More API
Check out code in `api/controller.go`