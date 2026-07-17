# Security Policy

## Supported Versions

ignis does not yet follow a formal release cadence. Only the latest commit on `main` is supported — please update before reporting an issue.

## Reporting a Vulnerability

Please report security vulnerabilities privately, not through a public GitHub issue.

Use [GitHub's private vulnerability reporting](https://github.com/THD-Spatial-AI/ignis/security/advisories/new) (Security tab → **Report a vulnerability**). This opens a private advisory thread with the maintainers — the report stays hidden from the public repository until a fix is out.

You should hear back within a week. If the report is valid, we'll work with you on a fix and coordinate disclosure timing before anything is made public.

## Scope

A few things about ignis's design that are **intentional, documented limitations**, not vulnerabilities to report:

- **The API key is a prototype-stage credential** (see the [arc42 architecture doc](https://thd-spatial-ai.github.io/ignis/documentation.pdf), ADR-004). It is not a substitute for real authentication and must not be relied on in a production deployment.
- **ignis has no authentication of its own.** It is designed to run behind a reverse proxy on a private network, never exposed directly to the internet. Deploying it without a proxy in front is a misconfiguration, not a vulnerability in ignis itself.

If you find a genuine issue within that design — for example, a way to bypass the reverse proxy's checks, an injection vulnerability, or a way to access data outside the caller's intended scope — please report it as above.
