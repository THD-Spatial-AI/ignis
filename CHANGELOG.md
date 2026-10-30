# Changelog

Notable changes to ignis. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/). Releases up to
v0.8.0 are described in their tag messages.

## [Unreleased]

### Changed

- `environment/docker-compose.yml` replaces the compose files under
  `environment/http/` and `environment/https/`. `--build` builds from the
  checkout, `--pull always --no-build` runs the release named by
  `IGNIS_IMAGE_TAG`, and `--profile https` adds the Caddy proxy.
- The Compose project is `ignis` and the services are `db`, `api`,
  `build-db` and `proxy`, so the API container is `ignis-api-1`. On
  `tentacron-net` the API answers to `ignis` and, until callers move,
  `ignis-http`.
- `env/` and `caddy/` move to `environment/`, and the two `.env.example`
  files merge into `environment/.env.example`. The data volume of a stack
  started from the old files is not reused: re-seed with `up`, or copy the
  old `db-data` volume across.
