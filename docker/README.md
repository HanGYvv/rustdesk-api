# RustDesk API S6

RustDesk API S6 is a self-hosted RustDesk container image built on
[s6-overlay](https://github.com/just-containers/s6-overlay). It bundles the
RustDesk API service together with the `hbbs` and `hbbr` servers into a single
compact runtime with proper service supervision.

- **Source repository:** <https://github.com/HanGYvv/rustdesk-api>
- **Docker Hub:** `hangyvv/rustdesk-api-s6`
- **GHCR:** `ghcr.io/hangyvv/rustdesk-api-s6`

## Highlights

- All-in-one RustDesk deployment: API + `hbbs` + `hbbr` in one image
- Built with s6-overlay for reliable service supervision and health checks
- Automatic key-pair handling via the server's `/data` volume
- API configured through `RUSTDESK_API_*` environment variables

## Exposed ports

| Port        | Service | Purpose                          |
| ----------- | ------- | -------------------------------- |
| `21114/tcp` | API     | HTTP API, Web Admin, Web Client  |
| `21115/tcp` | hbbs    | NAT type test                    |
| `21116/tcp` | hbbs    | TCP hole punching / heartbeat    |
| `21116/udp` | hbbs    | ID registration / heartbeat      |
| `21117/tcp` | hbbr    | Relay service                    |
| `21118/tcp` | hbbs    | Web client support               |
| `21119/tcp` | hbbr    | Web client support               |

## Environment variables

The API is configured entirely through environment variables prefixed with
`RUSTDESK_API_` (dots in config keys become underscores, e.g. the config key
`rustdesk.id-server` maps to `RUSTDESK_API_RUSTDESK_ID_SERVER`). The `RELAY`
and `ENCRYPTED_ONLY` variables are passed through to the bundled `hbbs` /
`hbbr` server.

| Variable                             | Default                  | Description                                        |
| ---                                  | ---                      | ---                                                |
| `RUSTDESK_API_RUSTDESK_ID_SERVER`    | `192.168.1.66:21116`     | Public address of the ID server (`hbbs`)           |
| `RUSTDESK_API_RUSTDESK_RELAY_SERVER` | `192.168.1.66:21117`     | Public address of the relay server (`hbbr`)        |
| `RUSTDESK_API_RUSTDESK_API_SERVER`   | `http://127.0.0.1:21114` | API server URL advertised to clients               |
| `RUSTDESK_API_RUSTDESK_KEY`          | _(unset)_                | Server public key; falls back to `key-file`        |
| `RUSTDESK_API_RUSTDESK_KEY_FILE`     | `/data/id_ed25519.pub`   | Path to the server public key file                 |
| `RUSTDESK_API_GIN_MODE`              | `release`                | Gin mode: `release`, `debug` or `test`             |
| `RELAY`                              | `relay.example.com`      | Public address advertised for the relay (`hbbr`)   |
| `ENCRYPTED_ONLY`                     | `0`                      | Set to `1` to force encrypted connections (`-k _`) |

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
  -e RELAY=<relay-server>:21117 \
  -p 21114:21114 \
  -p 21115-21119:21115-21119 \
  -p 21116:21116/udp \
  -v rustdesk-data:/data \
  -v rustdesk-api-data:/app/data \
  hangyvv/rustdesk-api-s6:latest
```

### docker-compose

```yaml
services:
  rustdesk-api:
    image: hangyvv/rustdesk-api-s6:latest
    container_name: rustdesk-api
    environment:
      - TZ=Asia/Shanghai
      - RUSTDESK_API_RUSTDESK_ID_SERVER=192.168.1.66:21116
      - RUSTDESK_API_RUSTDESK_RELAY_SERVER=192.168.1.66:21117
      - RUSTDESK_API_RUSTDESK_API_SERVER=http://127.0.0.1:21114
      - RUSTDESK_API_RUSTDESK_KEY=123456789
      - RELAY=relay.example.com
    ports:
      - "21114:21114"
      - "21115-21119:21115-21119"
      - "21116:21116/udp"
    volumes:
      - rustdesk-data:/data
      - rustdesk-api-data:/app/data
    restart: unless-stopped

volumes:
  rustdesk-data:
  rustdesk-api-data:
```

## Data & keys

All persistent state lives under two volumes:

- `/data` — the server key pair (`id_ed25519` / `id_ed25519.pub`) and server
  data. The API reads the public key from here by default.
- `/app/data` — the API database and runtime data.

Back both volumes up to preserve your server identity and API data across
container recreation.
