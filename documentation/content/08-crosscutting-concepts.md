# Crosscutting Concepts

## TLS termination

ignis's code has no TLS, no certificates, no HTTPS server. Encryption ends at the proxy: it holds the certificate, decrypts the request, and forwards it as plain HTTP over the internal Docker network. Any service behind the proxy gets this for free: certificates live in one place, not in every service.

The internal hop from proxy to app is plain HTTP, so it is only as safe as the network it crosses. That is fine while the network is one we control (a single Docker network). It stops being fine if that hop ever crosses infrastructure we don't control, which is why mutual TLS on that hop is listed as future work.

## Trust follows the CA, not the certificate

A client doesn't trust a certificate on its own. It trusts the certificate authority (CA) that signed it, and that trust then covers every certificate the same CA signs. This is why per-service proxies scale: sign each one's certificate with the same CA (a public CA for internet-facing services, a shared internal CA for service-to-service), and a client establishes trust once, not once per service.

Locally this is proven directly: `caddy trust` trusts the local CA once, and every certificate it signs afterward is trusted, including from a freshly recreated container.

## The API-key gate

Caddy checks for a valid `X-Api-Key` before forwarding to ignis. ignis knows nothing about this key: it never reads it.

This is a prototype stand-in, not production auth. The current caller, Building Configurator, runs in the browser, so any key it sends is visible in its page source to anyone with developer tools. A static key shipped to the browser stops casual access but keeps no real secret. It stands in for the real caller in the target architecture, an orchestration layer running as a backend service, which can hold a credential the client never sees. What replaces the static key is settled in ADR-4.

## CORS: two concerns kept apart

A cross-origin request with a custom header triggers a preflight `OPTIONS` before the real request. ignis's own CORS middleware only knows `Content-Type`, correctly, since `X-Api-Key` is not part of ignis's API; it exists only because of the proxy. Teaching ignis about it would tie ignis's code to its deployment.

So Caddy answers the preflight itself, advertising the headers and origins it accepts. The real request still passes ignis's own origin check (`ALLOWED_ORIGINS`) when it arrives. Each layer handles its own concern.

## Open items

- Mutual TLS on the proxy-to-app hop, for deployments where that hop leaves infrastructure we control.

- Replacing the static API key with a credential issued by the real orchestration layer, once it exists.
