# GophKeeper

## Implemented

- Registration
- Authentication

## TODO

- Server connection retry

## Server configuration

Supports env, flags and config files. Priority in descending order: env, flag, config file.

| Env Variable        | CLI Flag    | Field               | Description                                     | Default     |
| ------------------- | ----------- |---------------------| ----------------------------------------------- | ----------- |
| `GRPC_PORT`         | `-p`        | `grpc_port`         | gRPC server port (e.g. `:8080`)                 | none        |
| `ENABLE_TLS`        | `-tls`      | `enable_TLS`        | Enable TLS (requires `CERT_PATH` and `PK_PATH`) | `false`     |
| `CERT_PATH`         | `-tls-cert` | `cert_path`         | Path to TLS certificate file                    | none        |
| `PK_PATH`           | `-tls-pk`   | `pk_path`           | Path to TLS private key file                    | none        |
| `CONFIG`            | `-c`        | `config`            | Path to configuration file                      | none        |
| `FILE_STORAGE_PATH` | `-s`        | `file_storage_path` | Path to file storage directory                  | `./storage` |
| `DATABASE_DSN`      | `-d`        | `database_dsn`      | Database connection string (DSN)                | none        |
| `MASTER_PASSWORD`   | —           | `master_password`   | Master password for encryption/decryption       | none        |
| `DEBUG`             | `-debug`    | `debug`             | Enable debug mode                               | `false`     |

## Client configuration 

Supports only flags.

| Flag    | Type   | Default          | Description                  |
| ------- | ------ | ---------------- | ---------------------------- |
| `-a`    | string | `localhost:8080` | Server address (`host:port`) |
| `-tls`  | bool   | `false`          | Enable TLS                   |
| `-cert` | string | `""`             | Path to TLS certificate file |
