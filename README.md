# drim

One-click deployment tool for self-hosted [Kaneo](https://kaneo.app) instances.

All you need. Nothing you don't.

![Demo](demo/demo.gif)

drim generates your `docker-compose.yml`, `Caddyfile` and `.env`, installs Docker if it is
missing, and brings the stack up, so a fresh server goes from nothing to a running Kaneo in
one command.

## Install

```bash
curl -fsSL https://assets.kaneo.app/install.sh | sh
```

Then deploy:

```bash
drim setup
```

Or do both at once, with automatic HTTPS for your domain:

```bash
curl -fsSL https://assets.kaneo.app/install.sh | sh -s -- --setup --domain=kaneo.example.com
```

## Documentation

Full documentation lives in the Kaneo docs:

- [Deploy with drim](https://docs.kaneo.app/core/installation/drim) covers every command,
  the install script options, configuration, and requirements
- [Migrating to drim](https://docs.kaneo.app/core/installation/drim-migration) covers moving
  an existing deployment onto drim, and updating an older one to the current
  single-container layout
- [Environment variables](https://docs.kaneo.app/core/installation/environment-variables)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Bugs and feature requests go in
[issues](https://github.com/usekaneo/drim/issues); vulnerabilities go to
[SECURITY.md](SECURITY.md) so they can be handled privately.

## License

MIT. See [LICENSE](LICENSE).
