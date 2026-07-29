#!/bin/sh
# Derives ALLOWED_ORIGINS_PATTERN (the regexp the Caddyfile's header_regexp
# matcher needs) from ALLOWED_ORIGINS (the plain comma list ignis's own Go
# CORS middleware reads) before exec'ing into the real command. Exporting it
# here, rather than duplicating this derivation inline in each compose
# file's entrypoint, means it's also reachable from `caddy adapt`/`validate`
# run manually inside the container, not just from the process this script starts.
set -eu
: "${IGNIS_API_KEY:?IGNIS_API_KEY is not set}"
ALLOWED_ORIGINS_PATTERN="^($(echo "${ALLOWED_ORIGINS:-}" | sed -e 's/[.]/\./g' -e 's/,/|/g'))$"
export ALLOWED_ORIGINS_PATTERN
exec "$@"
