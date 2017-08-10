#!/bin/bash
{{#with cfg.postgres}}
PGPORT="{{port}}"
PGHOST="{{host}}"
PG_SUPERUSER="{{superuser_name}}"
PG_SUPERUSER_PASSWORD="{{superuser_password}}"
PG_SUPERUSER_URI="postgres://{{superuser_name}}:{{superuser_password}}@{{host}}:{{port}}/postgres"
{{/with}}