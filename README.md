# Introduction
This is a skeleton application to demonstrate a basic Golang application suitable for deployment on the cloud.gov.au CloudFoundry environment that:

1. Uses a Postgresql database.
2. Uses a vendor directory.
3. Uses CircleCI with push-to-deploy.

This is very much a work-in-progress.

## Pre-work

Create the database:

```bash
cf create-service postgres shared go-helloworld-db
```

Bind our app to it - note, that since we're using postgresql, we want to pick a username, else we lose access if we rename the app:

```bash
cf bind-service go-helloworld go-helloworld-db -c '{"username":"gohelloworlduser"}'
```

(TODO: see if there's a better way to fix this)

Push it:

```bash
cf push
```

Once-off, visit the bootstrap page: <https://\<url\>/bootstrap>