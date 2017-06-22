# Introduction

This is a skeleton application to demonstrate a basic Golang application suitable for deployment on the cloud.gov.au CloudFoundry environment that:

1. Uses a Postgresql database.
1. Uses a vendor directory.
1. Uses CircleCI with push-to-deploy.

This is very much a work-in-progress.

## Pre-work

Create the database:

```bash
cf create-service postgres shared go-helloworld-db
```

Update `manifest.yml` to include a reference to our database service:

```yaml
  services:
  - go-helloworld-db
```

Push it (recommended to use the [Blue/Green deployer plugin for CF](https://github.com/bluemixgaragelondon/cf-blue-green-deploy)):

```bash
cf blue-green-deploy go-helloworld
```

(or, if you prefer downtime: `cf push`)

Once-off, visit the bootstrap page: <https://your.site/bootstrap>