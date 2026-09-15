# Runtime View

A browser client calling ignis produces two requests: a CORS preflight, then the real request. ignis answers both itself, whether or not a proxy sits in front of it.

## Step 1: preflight (OPTIONS)

Building Configurator calls `GET /api/v1/fields` from a page on `http://localhost:5173`. The target is a different origin, so the browser first sends an automatic `OPTIONS` preflight to ask permission.

ignis's CORS middleware matches the `Origin` against `ALLOWED_ORIGINS`, echoes it back alongside the methods and headers it accepts, and returns `204 No Content`. In the HTTPS environment Caddy forwards the `OPTIONS` through unchanged; it holds no CORS configuration of its own.

## Step 2: the real request

The browser sends `GET /api/v1/fields`. ignis checks the `Origin` again, adds `Access-Control-Allow-Origin` to the response, and returns the payload. Where Caddy is in front, it decrypts the request, forwards it to `ignis-app:8080` as plain HTTP over the internal Docker network, and passes the response back unchanged.

## Step 3: a request from an origin that is not allowed

The response is still produced, but without `Access-Control-Allow-Origin`, so the browser discards it before the calling page can read it. This is a browser-side rule enforced on behalf of that page. A server-to-server caller sends no `Origin` header and is unaffected by the list, which is why reachability is decided by the port mapping rather than by `ALLOWED_ORIGINS`.

## Why this matters

ignis's request handling is the same in both environments. Caddy adds TLS termination and nothing else, so moving between them changes the scheme and the port, never the behaviour of an endpoint.
