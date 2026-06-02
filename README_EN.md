# RustDesk API

[Chinese Doc](README.md)

This project uses Go to implement the RustDesk API and includes Web Admin and a Web Client.

[![version](https://img.shields.io/github/v/tag/HanGYvv/rustdesk-api?label=version)](https://github.com/HanGYvv/rustdesk-api/releases)
[![license](https://img.shields.io/github/license/HanGYvv/rustdesk-api)](LICENSE)
[![ci](https://github.com/HanGYvv/rustdesk-api/actions/workflows/ci.yml/badge.svg)](https://github.com/HanGYvv/rustdesk-api/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/golang-1.26.3-blue)](https://go.dev/)
[![Gin](https://img.shields.io/badge/gin-v1.9.0-lightBlue)](https://github.com/gin-gonic/gin)
[![Gorm](https://img.shields.io/badge/gorm-v1.25.7-green)](https://gorm.io/)
[![Swag](https://img.shields.io/badge/swag-v1.16.3-yellow)](https://github.com/swaggo/swag)
[![Go Report Card](https://goreportcard.com/badge/github.com/HanGYvv/rustdesk-api/v2)](https://goreportcard.com/report/github.com/HanGYvv/rustdesk-api/v2)

## Better with [RustDesk Server](https://github.com/HanYvv/rustdesk-server)

> [RustDesk Server](https://github.com/HanYvv/rustdesk-server) is forked from the official RustDesk Server repository and merged with the following features from [lejianwen/rustdek-server](https://github.com/lejianwen/rustdesk-server):
>
> 1. Solves the API connection timeout issue
> 2. Can force login before initiating a connection
> 3. Supports client websocket

## Features

- PC API
  - Personal API
  - Login
  - Address Book
  - Groups
  - Authorized login: supports `github`, `google` and `OIDC` login, supports `web admin` authorized login, and supports `LDAP` (tested with AD and OpenLDAP) if the API Server is configured
  - i18n
- Web Admin
  - User management
  - Device management
  - Address book management
  - Tag management
  - Group management
  - OAuth management
  - Configure LDAP by config file or environment variables
  - Login logs
  - Connection logs
  - File transfer logs
  - Quick access to web client
  - i18n
  - Share to guests through web client
  - Server control: some official simple commands [WIKI](https://github.com/HanGYvv/rustdesk-api/wiki/Rustdesk-Command)
- Web Client
  - Automatically obtain API server
  - Automatically obtain ID server and KEY
  - Automatically obtain address book
  - Guests can directly access the device through a temporary share link
- CLI
  - Reset admin password

## Usage

### API Service

The PC-side base interfaces are implemented. Personal version APIs are supported and can be enabled by configuring the `rustdesk.personal` file or the `RUSTDESK_API_RUSTDESK_PERSONAL` environment variable.

#### Login

![Login](docs/en_img/pc_login.png)

#### Address Book and Groups

| Address Book | Groups |
| --- | --- |
| ![Address Book](docs/en_img/pc_ab.png) | ![Groups](docs/en_img/pc_gr.png) |

### Web Admin

- The frontend and backend are separated to provide a user-friendly management interface, mainly for management and display. The frontend code is in [rustdesk-api-web](https://github.com/HanGYvv/rustdesk-api-web)
- The backend access address is `http://<your server>[:port]/_admin/`
- The administrator `username` during the initial installation is `admin`. The password will be printed in the console, and you can change it through the [command line](### CLI)

![Initial admin password](./docs/init_admin_pwd.png)

  1. Admin interface
  ![web_admin](docs/web_admin.png)

  2. Regular user interface
  ![web_user](docs/web_admin_user.png)

  3. Each user can have multiple address books, and address books can also be shared with other users

  4. Groups can be customized for easy management. Two types are supported for now: `shared group` and `regular group`

  5. You can directly open the web client for convenience; you can also share it with guests, and guests can directly access the device through the web client

  6. OAuth support:
  Currently, `GitHub`, `Google` and `OIDC` are supported. You need to create an `OAuth App` and configure it in the admin panel. For `Google` and `GitHub`, `Issuer` and `Scopes` do not need to be filled in; for `OIDC`, `Issuer` is required. `Scopes` is optional and defaults to `openid,profile,email`. Make sure you can obtain `sub`, `email` and `preferred_username`. Create the `GitHub OAuth App` in `Settings` -> `Developer settings` -> `OAuth Apps` -> `New OAuth App`, and the address is [Developer settings](https://github.com/settings/developers). Fill `Authorization callback URL` with `http://<your server[:port]>/api/oidc/callback`, for example `http://127.0.0.1:21114/api/oidc/callback`

  7. Login logs

  8. Connection logs

  9. File transfer logs

  10. Server control
  -`Simple mode`: some simple commands have been visualized and can be executed directly in the backend
  ![rustdesk_command_simple](./docs/rustdesk_command_simple.png)
  -`Advanced mode`: commands can be executed directly in the backend

  11. **LDAP support**: when LDAP is set on the API Server (tested with AD and LDAP), you can log in using user information from LDAP. If LDAP authentication fails, it falls back to local users.

### Web Client

  1. If you are already logged into the backend, the web client will log in automatically
  2. If you are not logged into the backend, click the login button in the upper-right corner, and the API server has already been configured
  3. After logging in, the ID server and KEY will be automatically synchronized
  4. After logging in, the address book will be automatically saved in the web client for convenient use

### Automated Documentation

Swag is used to generate API documentation, making it easier for developers to understand and use the API.

  1. Backend docs `<your server[:port]>/admin/swagger/index.html`

  2. PC docs `<your server[:port]>/swagger/index.html`
  ![api_swag](docs/api_swag.png)

### CLI

```bash
# Help
./apimain -h
```

#### Reset admin password

```bash
./apimain reset-admin-pwd <pwd>
```

## Installation and Running

### Related Configuration

  *[Config file](./conf/config.yaml)
  *Refer to the `conf/config.yaml` file and modify the related configuration
  *If `gorm.type` is `sqlite`, MySQL-related configuration is not required
  *If language is not set, the default is `zh-CN`

### Environment Variables

Environment variables correspond one-to-one with the configuration in `conf/config.yaml`. The variable name prefix is `RUSTDESK_API`
The following table does not list all entries. Please refer to the configuration in `conf/config.yaml`.

| Variable Name                                     | Description                                                                                                                      | Example                      |
| ------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | ---------------------------- |
| TZ                                                | Timezone                                                                                                                         | Asia/Shanghai                |
| RUSTDESK_API_LANG                                 | Language                                                                                                                         | `en`,`zh-CN`                 |
| RUSTDESK_API_APP_WEB_CLIENT                       | Whether to enable web-client; 1: enabled, 0: disabled; enabled by default                                                        | 1                            |
| RUSTDESK_API_APP_REGISTER                         | Whether to enable registration; `true`, `false`; default `false`                                                                 | `false`                      |
| RUSTDESK_API_APP_SHOW_SWAGGER                     | Whether Swagger docs are visible; `1` shown, `0` hidden; default `0` shown                                                       | `1`                          |
| RUSTDESK_API_APP_TOKEN_EXPIRE                     | Token validity duration                                                                                                          | `168h`                       |
| RUSTDESK_API_APP_DISABLE_PWD_LOGIN                | Whether to disable password login; `true`, `false`; default `false`                                                              | `false`                      |
| RUSTDESK_API_APP_REGISTER_STATUS                  | Default status for registered users; 1 enabled, 2 disabled, default 1                                                            | `1`                          |
| RUSTDESK_API_APP_CAPTCHA_THRESHOLD                | Captcha trigger count; -1 disabled, 0 always enabled, >0 enable after failed logins; default `3`                                 | `3`                          |
| RUSTDESK_API_APP_BAN_THRESHOLD                    | Ban IP trigger count; 0 disabled, >0 ban after failed logins; default `0`                                                        | `0`                          |
| -----ADMIN configuration-----                     |                                                                                                                                  |                              |
| RUSTDESK_API_ADMIN_TITLE                          | Admin title                                                                                                                      | `RustDesk Api Admin`         |
| RUSTDESK_API_ADMIN_HELLO                          | Admin welcome message; HTML is supported                                                                                         |                              |
| RUSTDESK_API_ADMIN_HELLO_FILE                     | Admin welcome message file; more convenient when the content is long. It overrides `RUSTDESK_API_ADMIN_HELLO`                    | `./conf/admin/hello.html`    |
| -----GIN configuration-----                       |                                                                                                                                  |                              |
| RUSTDESK_API_GIN_TRUST_PROXY                      | Trusted proxy IP list, separated by `,`, trusted by default                                                                      | 192.168.1.2,192.168.1.3      |
| -----GORM configuration-----                      |                                                                                                                                  |                              |
| RUSTDESK_API_GORM_TYPE                            | Database type `sqlite` or `mysql`, default `sqlite`                                                                              | sqlite                       |
| RUSTDESK_API_GORM_MAX_IDLE_CONNS                  | Maximum idle connections                                                                                                         | 10                           |
| RUSTDESK_API_GORM_MAX_OPEN_CONNS                  | Maximum open connections                                                                                                         | 100                          |
| RUSTDESK_API_RUSTDESK_PERSONAL                    | Whether to enable Personal API; 1 enabled, 0 disabled; enabled by default                                                        | 1                            |
| -----MYSQL configuration-----                     |                                                                                                                                  |                              |
| RUSTDESK_API_MYSQL_USERNAME                       | MySQL username                                                                                                                   | root                         |
| RUSTDESK_API_MYSQL_PASSWORD                       | MySQL password                                                                                                                   | 111111                       |
| RUSTDESK_API_MYSQL_ADDR                           | MySQL address                                                                                                                    | 192.168.1.66:3306            |
| RUSTDESK_API_MYSQL_DBNAME                         | MySQL database name                                                                                                              | rustdesk                     |
| RUSTDESK_API_MYSQL_TLS                            | Whether to enable TLS, optional values: `true`, `false`, `skip-verify`, `custom`                                                 | `false`                      |
| -----RUSTDESK configuration-----                  |                                                                                                                                  |                              |
| RUSTDESK_API_RUSTDESK_ID_SERVER                   | RustDesk ID server address                                                                                                       | 192.168.1.66:21116           |
| RUSTDESK_API_RUSTDESK_RELAY_SERVER                | RustDesk relay server address                                                                                                    | 192.168.1.66:21117           |
| RUSTDESK_API_RUSTDESK_API_SERVER                  | RustDesk API server address                                                                                                      | <http://192.168.1.66:21114>  |
| RUSTDESK_API_RUSTDESK_KEY                         | RustDesk key                                                                                                                     | 123456789                    |
| RUSTDESK_API_RUSTDESK_KEY_FILE                    | RustDesk key file                                                                                                                | `./conf/data/id_ed25519.pub` |
| RUSTDESK_API_RUSTDESK_WEBCLIENT_MAGIC_QUERYONLINE | Whether to enable the new online status query method in web client v2; `1` enabled, `0` disabled                                 | `0`                          |
| RUSTDESK_API_RUSTDESK_WS_HOST                     | Custom WebSocket Host                                                                                                            | `wss://192.168.1.123:1234`   |
| ----PROXY configuration-----                      |                                                                                                                                  |                              |
| RUSTDESK_API_PROXY_ENABLE                         | Whether to enable proxy: `false`, `true`                                                                                         | `false`                      |
| RUSTDESK_API_PROXY_HOST                           | Proxy address                                                                                                                    | `http://127.0.0.1:1080`      |
| ----JWT configuration----                         |                                                                                                                                  |                              |
| RUSTDESK_API_JWT_KEY                              | Custom JWT KEY; if empty, JWT is disabled. If `MUST_LOGIN` in `rustdesk-server` is not used, it is recommended to leave it empty |                              |
| RUSTDESK_API_JWT_EXPIRE_DURATION                  | JWT validity duration                                                                                                            | `168h`                       |

### Running

#### docker运行

<!-- markdownlint-disable MD029 -->

  1. Run directly with Docker. The configuration can be modified by mounting the config file `/app/conf/config.yaml`, or by overriding the configuration file with environment variables

  ```bash
  docker run -d --name rustdesk-api -p 21114:21114 \
  -v /data/rustdesk/api:/app/data \
  -e TZ=Asia/Shanghai \
  -e RUSTDESK_API_LANG=zh-CN \
  -e RUSTDESK_API_RUSTDESK_ID_SERVER=192.168.1.66:21116 \
  -e RUSTDESK_API_RUSTDESK_RELAY_SERVER=192.168.1.66:21117 \
  -e RUSTDESK_API_RUSTDESK_API_SERVER=http://192.168.1.66:21114 \
  -e RUSTDESK_API_RUSTDESK_KEY=<key> \
  HanGYvv/rustdesk-api
  ```

  2. Use `docker compose`, refer to [WIKI](https://github.com/HanGYvv/rustdesk-api/wiki)

<!-- markdownlint-enable MD029 -->

#### Download and run release directly

[Download](https://github.com/HanGYvv/rustdesk-api/releases)

#### Source installation

<!-- markdownlint-disable MD029 -->

1. Clone the repository

  ```bash
  git clone https://github.com/HanGYvv/rustdesk-api.git
  cd rustdesk-api
  ```

2. Install dependencies

  ```bash
  go mod tidy
  # Install swag; if you do not need to generate docs, you can skip it
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

3. Build the backend front-end; the front-end code is in [rustdesk-api-web](https://github.com/HanGYvv/rustdesk-api-web)

  ```bash
  cd resources
  mkdir -p admin
  git clone https://github.com/HanGYvv/rustdesk-api-web
  cd rustdesk-api-web
  npm install
  npm run build
  cp -ar dist/* ../admin/
  ```

4. Run

  ```bash
  # Run directly
  go run cmd/apimain.go
  # Or generate API and run using generate_api.go
  go generate generate_api.go
  ```

   > Note: when using `go run` or the compiled binary, the `conf` and `resources`
   > directories must exist in the current directory. If you run from another directory,
   > you can specify the absolute path with `-c` and `RUSTDESK_API_GIN_RESOURCES_PATH`, for example:
   >
   > ```bash
   > RUSTDESK_API_GIN_RESOURCES_PATH=/opt/rustdesk-api/resources ./apimain -c /opt/rustdesk-api/conf/config.yaml
   > ```

5. Build. If you want to compile it yourself, first go to the project root, then run `build.bat` on Windows or `build.sh` on Linux. The compiled executable files will be generated in the `release` directory. Run the compiled executable directly.

6. Open your browser and visit `http://<your server[:port]>/_admin/`, with default username/password `admin`. Please change the password in time.

<!-- markdownlint-enable MD029 -->

#### Running with the `s6` image

- The connection timeout issue has been solved
- You can force login before initiating a connection
- [rustdesk-server](https://github.com/HanGYvv/rustdesk-server)
- The `s6` build now uses `HanGYvv/rustdesk-server-s6:latest` as its base image

```yaml
networks:
  rustdesk-net:
    external: false
services:
  rustdesk:
    ports:
      - 21114:21114
      - 21115:21115
      - 21116:21116
      - 21116:21116/udp
      - 21117:21117
      - 21118:21118
      - 21119:21119
    image: HanGYvv/rustdesk-api-s6:latest
    environment:
      - RELAY=<relay_server[:port]>
      - ENCRYPTED_ONLY=1
      - MUST_LOGIN=N
      - TZ=Asia/Shanghai
      - RUSTDESK_API_RUSTDESK_ID_SERVER=<id_server[:21116]>
      - RUSTDESK_API_RUSTDESK_RELAY_SERVER=<relay_server[:21117]>
      - RUSTDESK_API_RUSTDESK_API_SERVER=http://<api_server[:21114]>
      - RUSTDESK_API_KEY_FILE=/data/id_ed25519.pub
      - RUSTDESK_API_JWT_KEY=xxxxxx # jwt key
    volumes:
      - /data/rustdesk/server:/data
      - /data/rustdesk/api:/app/data # mount the database
    networks:
      - rustdesk-net
    restart: unless-stopped
```

## Others

- [WIKI](https://github.com/HanGYvv/rustdesk-api/wiki)
- [Connection timeout issue](https://github.com/HanGYvv/rustdesk-api/issues/92)
- [Change client ID](https://github.com/abdullah-erturk/RustDesk-ID-Changer)
- [webclient source](https://hub.docker.com/r/keyurbhole/flutter_web_desk)

## Acknowledgements

Thanks to everyone who contributed!

[![Contributors](https://contrib.rocks/image?repo=HanGYvv/rustdesk-api)](https://github.com/HanGYvv/rustdesk-api/graphs/contributors)

## Thank you for your support! If this project helps you, please give it a ⭐️. Thank you
