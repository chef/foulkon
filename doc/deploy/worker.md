# Deploy Foulkon Worker

 You have to specify configuration file using flag -config-file. Using binary file command is `worker -config-file=/path/config.toml`
 
## Deploy with docker
Then, you can run the docker image, mounting (-v) a config.toml inside the container (you could also make a custom Dockerfile with "ADD my-custom-conf.toml /my-custom-conf.toml").
E.g. 
 ```
 docker run -v /home/myuser/foulkon/config.toml:/config.toml tecsisa/foulkon-worker -config-file=/config.toml
 ```
 
## Worker configuration file 
 This config file is a TOML file that has several parts:
 
### [server] 
| Server   | Server config properties              | Values                     | Default | Optional |
|----------|---------------------------------------|----------------------------|---------|----------|
| host     | Worker's hostname.                    | `localhost`                |         | No       |
| port     | Worker's port.                        | `8000`                     |         | No       |
| certfile | Absolute path for public certificate. | `/etc/secrets/public.pem`  |         | Yes      |
| keyfile  | Absolute path for private key.        | `/etc/secrets/private.pem` |         | Yes      |

__Note:__ Don't use Foulkon worker without certificate in production.

### [admin] 
| Admin user | Admin user configuration | Values     | Default | Optional |
|------------|--------------------------|------------|---------|----------|
| username   | Admin user name.         | `admin`    |         | No       |
| password   | Admin user password.     | `password` |         | No       |

__Note:__ Use a strong password for admin user in production.

### [logger] 
| Logger | Logger configuration properties.                        | Values                                                | Default   | Optional                    |
|--------|---------------------------------------------------------|-------------------------------------------------------|-----------|-----------------------------|
| type   | Type of logger to use.                                  | `file`, `default`                                     | `default` | Yes                         |
| level  | Log level.                                              | `debug`, `info`, `warning`, `error`, `fatal`, `panic` | `info`    | Yes                         |
| dir    | Full path where log file is. It won't be autogenerated. | `/tmp/foulkon.log`                                    |           | No if logger type is `file` |

### [database]
| Database | Database configuration | Values     | Default | Optional |
|----------|------------------------|------------|---------|----------|
| type     | Database backend type  | `postgres` |         | No       |

#### [database.postgres]
| PostgreSQL     | PostgreSQL configuration properties                          | Values                                                                 | Default | Optional |
|----------------|--------------------------------------------------------------|------------------------------------------------------------------------|---------|----------|
| datasourcename | Connection datasource including user, password and database. | `postgres://foulkon:password@localhost:5432/foulkondb?sslmode=disable` |         | No       |
| idleconns      | Idle connection number.                                      | `10`                                                                   | 5       | Yes      |
| maxopenconns   | Max open connection number.                                  | `20`                                                                   | 20      | Yes      |
| connttl        | Timeout for conenctions                                      | `200`                                                                  | 300     | Yes      |
 
### [authenticator]
| Authenticator | Authenticatior connector configuration properties        | Values | Default | Optional |
|---------------|----------------------------------------------------------|--------|---------|----------|
| type          | Type of connector that will be used. Only `oidc` at now. | `oidc` |         | No       |

#### [authenticator.oidc]
| OIDC      | OpenID Connect authenticatior connector configuration properties | Values                        | Default | Optional |
|-----------|------------------------------------------------------------------|-------------------------------|---------|----------|
| issuer    | Full url for token issuer.                                       | `https://accounts.google.com` |         | No       |
| clientids | List of allowed clients separated by `;`.                        | `clientId1;clientId2`         |         | No       |