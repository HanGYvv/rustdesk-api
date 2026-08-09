# RustDesk API

RustDesk API is a minimal self-hosted RustDesk container image. It runs only
the RustDesk API service (`apimain`) and relies on a separately deployed
`hbbs` / `hbbr` server for ID registration and relay.

- **Source repository:** <https://github.com/HanGYvv/rustdesk-api>
- **Docker Hub:** `hangyvv/rustdesk-api`
- **GHCR:** `ghcr.io/hangyvv/rustdesk-api`

> Looking for an all-in-one image that also bundles `hbbs` and `hbbr`? Use the
> s6-overlay variant: `hangyvv/rustdesk-api-s6`.

## Highlights

- Self-hosted RustDesk API service with Web Admin and Web Client
- Minimal Alpine-based image containing only the API service
- API configured through `RUSTDESK_API_*` environment variables

## Exposed ports

| Port        | Service | Purpose                          |
| ----------- | ------- | -------------------------------- |
| `21114/tcp` | API     | HTTP API, Web Admin, Web Client  |

## Environment variables

The API is configured entirely through environment variables prefixed with
`RUSTDESK_API_` (dots in config keys become underscores, e.g. the config key
`rustdesk.id-server` maps to `RUSTDESK_API_RUSTDESK_ID_SERVER`).

| Variable                             | Default                  | Description                                 |
| ---                                  | ---                      | ---                                         |
| `RUSTDESK_API_RUSTDESK_ID_SERVER`    | `192.168.1.66:21116`     | Public address of the ID server (`hbbs`)    |
| `RUSTDESK_API_RUSTDESK_RELAY_SERVER` | `192.168.1.66:21117`     | Public address of the relay server (`hbbr`) |
| `RUSTDESK_API_RUSTDESK_API_SERVER`   | `http://127.0.0.1:21114` | API server URL advertised to clients        |
| `RUSTDESK_API_RUSTDESK_KEY`          | _(unset)_                | Server public key; falls back to `key-file` |
| `RUSTDESK_API_RUSTDESK_KEY_FILE`     | `/data/id_ed25519.pub`   | Path to the server public key file          |
| `RUSTDESK_API_GIN_MODE`              | `release`                | Gin mode: `release`, `debug` or `test`      |

Any other config key can be overridden the same way, for example
`RUSTDESK_API_APP_REGISTER` or `RUSTDESK_API_GORM_TYPE`.

## Quick start

```bash
docker run -d \
  --name rustdesk-api \
  -e RUSTDESK_API_RUSTDESK_ID_SERVER=<id-server>:21116 \
  -e RUSTDESK_API_RUSTDESK_RELAY_SERVER=<relay-server>:21117 \
  -e RUSTDESK_API_RUSTDESK_API_SERVER=http://<host>:21114 \
  -e RUSTDESK_API_RUSTDESK_KEY=<public-key> \
  -p 21114:21114 \
  -v rustdesk-api-data:/app/data \
  hangyvv/rustdesk-api:latest
```

### docker-compose

```yaml
services:
  rustdesk-api:
    image: hangyvv/rustdesk-api:latest
    container_name: rustdesk-api
    environment:
      - TZ=Asia/Shanghai
      - RUSTDESK_API_RUSTDESK_ID_SERVER=192.168.1.66:21116
      - RUSTDESK_API_RUSTDESK_RELAY_SERVER=192.168.1.66:21117
      - RUSTDESK_API_RUSTDESK_API_SERVER=http://127.0.0.1:21114
      - RUSTDESK_API_RUSTDESK_KEY=123456789
    ports:
      - "21114:21114"
    volumes:
      - rustdesk-api-data:/app/data
    restart: unless-stopped

volumes:
  rustdesk-api-data:
```

## Data

The API database and runtime data live under `/app/data`. Back this volume up
to preserve your API data across container recreation.
