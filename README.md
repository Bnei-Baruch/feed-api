# Feed API

This service mainly provides content feed for mobile archive application (latest lessons and other content units).
This service reads MDB and have alocal synced version, it also reads Chronicles and keeps a window of chronicles always updating it.
The service also provides recommendations and servers views (content popularity).

## Install

```console
dep ensure
```

Migration tool
```
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz | tar xvz
```
