# Security Policy

## Reporting a vulnerability

Please report security issues privately. Do not open a public issue, and do not
include details in a pull request.

Use [GitHub's private vulnerability reporting](https://github.com/usekaneo/drim/security/advisories/new),
which notifies the maintainers directly and keeps the discussion private until a
fix ships. If you cannot use it, email [andrej@kaneo.app](mailto:andrej@kaneo.app)
instead.

Helpful things to include, as far as you have them:

- the drim version (`drim --version`) and the host OS
- the command you ran and what it generated
- what an attacker gains, and the access they need to do it

## Why this matters here

drim runs with elevated privileges, installs Docker, and writes a `.env` file
containing your database password and `AUTH_SECRET`. It is also installed by
piping a shell script from the network. Issues in any of those paths are worth
reporting even when they look minor.

## What to expect

- An acknowledgement within 3 days.
- An assessment of severity and scope within 7 days.
- A fix released as soon as it is ready.
- Credit in the advisory, unless you would rather stay anonymous.

## Supported versions

Fixes land in the latest release. drim can update itself with `drim update`.

## Scope

In scope: the drim CLI, `install.sh`, and the `docker-compose.yml`, `Caddyfile`
and `.env` it generates.

Out of scope: vulnerabilities in Docker, Caddy or PostgreSQL themselves, and
issues that require an already-compromised host. For vulnerabilities in Kaneo
itself, use the [Kaneo security policy](https://github.com/usekaneo/kaneo/security/policy).
