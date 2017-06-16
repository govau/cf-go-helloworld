#!/usr/bin/env bash

set -e
set -o pipefail

# TODO - using the password in this manner means that anyone else running `ps` on the host might see it.
# Doesn't look there are any better alternatives for now

cf api $CF_API
cf auth $CF_USER $CF_PASSWORD
cf target -o $CF_ORG -s $CF_SPACE
cf push
