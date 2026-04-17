# envoy-cli

A CLI tool for managing and syncing `.env` files across environments with secret masking support.

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envoy-cli/releases).

## Usage

```bash
# Sync .env file to a target environment
envoy-cli sync --env production --file .env

# List all variables with secrets masked
envoy-cli list --mask-secrets

# Pull environment variables from a remote source
envoy-cli pull --env staging --output .env.staging

# Diff two environment files
envoy-cli diff .env .env.production

# Validate a .env file for syntax errors
envoy-cli validate --file .env
```

### Example

```bash
$ envoy-cli list --file .env --mask-secrets

DATABASE_URL=postgres://user:****@localhost:5432/mydb
API_KEY=****
DEBUG=true
PORT=8080
```

## Configuration

Create an `envoy.yaml` file in your project root to define environments and sync targets:

```yaml
environments:
  production:
    source: s3://my-bucket/envs/production.env
  staging:
    source: s3://my-bucket/envs/staging.env
```

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

## License

This project is licensed under the [MIT License](LICENSE).
