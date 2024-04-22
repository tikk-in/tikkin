<p align="center">
  <img width="300" height="300" src="assets/logo.png">
</p>

# Tikk.in

Tikk.in is a simple, headless, lightweight, and blazing fast URL shortener.
It is built using Go, PostgreSQL, and Redis (optional).

## Features

- **Headless**: No mandatory frontend, just a REST API. (A simple, optional UI will be provided in the future)
- **Lightweight**: Tikk.in is built using Go, which makes it lightweight and fast.
- **Link Expiry**: You can set an expiry date for the shortened URL.
- **Custom Alias**: You can set a custom alias for the shortened URL.
- **Configuration as Code**: Tikk.in can be fully configured using environment variables.

## Installation

### Docker compose

It will build the Docker image and run tikkin + PostgreSQL
```bash
docker compose up
```

### Docker image

```bash
docker pull tikkin/tikkin:latest
```

### From source

Clone the repository and build
```bash
go mod download
go build -o tikkin

# Run the binary
./tikkin --admin-password=<your-pass> --config=./your-config.yml
```

## Configuration

Tikk.in requires a configuration file to run. Check out the [example configuration file](example.config.yml) for more
information.

The configuration file path can be provided using the `CONFIG_PATH` environment variable or by using the `--config`
flag.

Configuration options table:

| Name                  | Description                            | Flag           | Environment Variable | Default Value         | Required |
|-----------------------|----------------------------------------|----------------|----------------------|-----------------------|----------|
| config path           | The path to the configuration file.    | --config       | CONFIG_PATH          | `./config.yml`        | true     |
| `server.port`         | The port on which the server will run. |                |                      | `3000`                | true     |
| `server.jwt.secret`   | The database configuration.            | jwt-secret     | SERVER_JWT_SECRET    | `changemeplease`      | true     |
| `db.host`             | The database host.                     |                |                      | `localhost`           | true     |
| `db.port`             | The database port.                     |                |                      | `5432`                | true     |
| `db.user`             | The database user.                     |                |                      | `tikkin`              | true     |
| `db.password`         | The database password.                 |                |                      | `tikkin`              | true     |
| `db.database`         | The database name.                     |                |                      | `tikkin`              | true     |
| `db.connections`      | The number of database connections.    |                |                      | `10`                  | true     |
| `site.name`           | The site name.                         |                |                      | `Tikk.in`             | false    |
| `site.url`            | The site URL.                          |                |                      | `https://example.com` | false    |
| `admin.email`         | The root admin email.                  |                |                      |                       | false    |
| `admin.password`      | The root admin password.               | admin-password | ADMIN_PASSWORD       |                       | false    |
| `links.length`        | The length of the shortened URL.       |                |                      | `6`                   | true     |
| `email.enabled`       | Enable email notifications.            |                |                      | `true`                | false    |
| `email.smtp.host`     | The SMTP host.                         |                |                      |                       | false    |
| `email.smtp.port`     | The SMTP port.                         |                |                      | `587`                 | false    |
| `email.smtp.username` | The SMTP username.                     |                |                      |                       | false    |
| `email.smtp.password` | The SMTP password.                     | smtp-password  | SMTP_PASSWORD        |                       | false    |
| `email.from`          | The email from address.                |                |                      | `noreply@example.com` | false    |



