# Flightdeck

## Local development with webpack-dev-server (hot-reload)

### Start webpack-dev-server in web folder

```shell
npm run start
```

### Start the go backend

Create a `Namespace`

```shell
kubectl create ns flightdeck
```

Create a `ServiceAccount` and `ClusterRoleBinding`

```shell
kubectl apply -f deploy/admin.user.yaml
```

Generate a bearer token for the `admin-user`

```shell
kubectl -n flightdeck create token admin-user
```

```shell
go run .
```

Now open your browser `localhost:8000` and copy the token and paste it into the `Bearer token` field on the login screen.

Click the `Sign in` button and that's it. You are now logged in as an admin.

## Local development with go server

### Start go backend (with proxy) in api folder

```shell
go run . -fileserverpath=../web/build -fileserver=true -addr=0.0.0.0:8000
```

## License

[Apache License 2.0](https://github.com/autovia/flightdeck/blob/master/LICENSE)

----
_Copyright [Autovia GmbH](https://autovia.de)_
