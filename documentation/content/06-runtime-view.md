# Runtime View

A browser client calling ignis through the proxy produces two requests: a CORS preflight, then the real request. This section traces both, plus a rejected request.

## Step 1: preflight (OPTIONS)

Building Configurator calls `GET /api/v1/fields` from a page on `http://localhost:5173`. The target is a different origin (`https://localhost`) and the request carries a custom header (`X-Api-Key`), so the browser first sends an automatic `OPTIONS` preflight to ask permission. A preflight never carries the custom header itself.

Caddy matches the preflight on method alone, before the API-key check, and answers it: it echoes the allowed origin, lists the allowed methods and headers, and returns `204 No Content`. The preflight never reaches ignis.

## Step 2: the real request

The browser now sends `GET /api/v1/fields` with `X-Api-Key` attached. Caddy checks the key. On a match, it forwards the request to `ignis-app:8080` over the internal Docker network as plain HTTP (encryption already ended at Caddy). ignis handles it as if no proxy existed: it checks the `Origin` against `ALLOWED_ORIGINS` and adds `Access-Control-Allow-Origin` to the response. Caddy passes the response back unchanged.

## Step 3: rejected request

A request with a missing or wrong `X-Api-Key` hits Caddy's fallback rule and gets `403 Forbidden`. It never reaches ignis. ignis publishes no host port, so the proxy is the only way in.

## Why this matters

ignis behaves identically with or without the proxy in front of it. Caddy adds the TLS, the key check, and the preflight response; ignis's code knows nothing about any of it.
