# Crosscutting Concepts

## TLS termination

ignis's code has no TLS, no certificates, no HTTPS server. That makes TLS a property of a deployment rather than of the service, and the two Compose environments are the two answers: `environment/https` puts Caddy in front to hold the certificate and decrypt, `environment/http` leaves it to whatever already terminates TLS upstream.

Where Caddy is used, the hop from proxy to app is plain HTTP, so it is only as safe as the network it crosses. That is fine while the network is one we control (a single Docker network). It stops being fine if that hop ever crosses infrastructure we don't control, which is why mutual TLS on that hop is listed as future work.

## Trust follows the CA, not the certificate

A client doesn't trust a certificate on its own. It trusts the certificate authority (CA) that signed it, and that trust then covers every certificate the same CA signs. This is why per-service proxies scale: sign each one's certificate with the same CA (a public CA for internet-facing services, a shared internal CA for service-to-service), and a client establishes trust once, not once per service.

Locally this is proven directly: `caddy trust` trusts the local CA once, and every certificate it signs afterward is trusted, including from a freshly recreated container.

## The network is the trust boundary

ignis checks no credential. It is deployed on an internal network reachable only over VPN, behind a platform that has already authenticated the end user, so an access check in front of ignis would gate traffic that is already gated. See ADR-004.

What decides who can reach ignis is therefore its port mapping. `ignis-db` publishes nothing in either environment; `ignis-app` publishes nothing in the HTTPS environment and binds to the host loopback in the HTTP one. Widening that binding is the deliberate act that exposes the service.

## CORS is a browser rule, not an access control

A cross-origin request from a page triggers a preflight `OPTIONS` before the real request. ignis answers it from its own middleware, matching `Origin` against `ALLOWED_ORIGINS` and advertising the methods and headers it accepts, and applies the same check to the real request.

This is enforced by the browser on behalf of the page it loaded. A caller that is not a browser sends no `Origin` header and ignores the response headers, so `ALLOWED_ORIGINS` neither admits nor blocks it. The list protects a user's session in a browser tab; it is not a substitute for deciding reachability at the network layer.

## Open items

- Mutual TLS on the proxy-to-app hop, for deployments where that hop leaves infrastructure we control.

- A caller credential, should ignis ever be deployed somewhere the network is not the trust boundary.
