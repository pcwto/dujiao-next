# VPS API Deployment

This directory contains the VPS runtime compose file for the backend API image published by CI/CD.

The public `user` and `admin` frontends are deployed separately on Vercel. This compose file runs only the Go API container and connects it to Supabase Postgres.

## Published GHCR Package

Use the backend package:

- `ghcr.io/<IMAGE_OWNER>/dujiao-next-server:<IMAGE_TAG>`

The frontend packages may still be published by CI/CD, but they are not part of the VPS runtime when Vercel hosts the frontends.

## Tag Policy

The publish workflow produces these tag families:

- `main` / `master`: mutable branch tags for rolling environments
- `vX.Y.Z`: release tags for human-friendly deployments
- `sha-<commit>`: immutable build-oriented tags for automation and rollback
- `latest`: published only from the repository default branch for convenience pulls

For production deployments, prefer `vX.Y.Z` or a digest-backed deployment flow. `main`, `master`, and `latest` can move over time.

## Prepare Environment

Copy the template on the VPS and fill in real values:

```bash
cp deploy/.env.example deploy/.env.prod
```

Required values:

- `IMAGE_OWNER`
- `IMAGE_TAG`
- `DATABASE_DSN`
- `CORS_ALLOWED_ORIGINS`
- `APP_SECRET_KEY`
- `JWT_SECRET`
- `USER_JWT_SECRET`
- `DJ_DEFAULT_ADMIN_PASSWORD`

Supabase should be configured as PostgreSQL:

```env
DATABASE_DRIVER=postgres
DATABASE_DSN=postgresql://postgres.<project-ref>:<password>@aws-0-<region>.pooler.supabase.com:5432/postgres?sslmode=require
```

For a long-running VPS API, prefer Supabase's session pooler on port `5432`.
Use the transaction pooler on port `6543` mainly for short-lived/serverless workloads.

`CORS_ALLOWED_ORIGINS` must match the exact Vercel origins that will call the API, for example:

```env
CORS_ALLOWED_ORIGINS=https://your-shop.vercel.app,https://your-admin.vercel.app
CORS_ALLOW_CREDENTIALS=true
```

Use different random values for `APP_SECRET_KEY`, `JWT_SECRET`, and `USER_JWT_SECRET`; each should be at least 32 characters.

## Start API

```bash
docker compose --env-file deploy/.env.prod -f deploy/docker-compose.ghcr.yml up -d
```

## GitHub Actions VPS Deployment

`Deploy VPS API` runs after `Publish Docker Image` completes on `main` or `master`. It deploys the immutable `sha-<commit>` tag produced by the publish workflow, updates the VPS env file, restarts only the `server` service, and checks the public health endpoint.

Configure these GitHub repository secrets in `pcwto/dujiao-next`:

- `VPS_HOST`: VPS IP or DNS name
- `VPS_PORT`: SSH port, optional; defaults to `22` when empty
- `VPS_USER`: SSH user
- `VPS_SSH_KEY`: private key with access to the deployment directory
- `VPS_COMPOSE_DIR`: deployment root, for example `/opt/dujiao-next-server`
- `PUBLIC_API_HEALTH_URL`: public health URL, for example `https://mmmslai.win/dujiao-next/health`

The deploy workflow auto-detects these existing server layouts under `VPS_COMPOSE_DIR`:

- compose file: `compose/docker-compose.yml`, `docker-compose.yml`, or `deploy/docker-compose.ghcr.yml`
- env file: `env/server.env`, `.env.prod`, or `deploy/.env.prod`

It only updates `IMAGE_OWNER` and `IMAGE_TAG` in the env file. Database, SMTP, payment, JWT, and admin secrets must remain in the VPS env file or in the admin console.

## SMTP and Payment Operations

SMTP is configured from the admin console at `Settings -> SMTP`. Use the same values as the env template only when you want the container default to be email-enabled:

- `EMAIL_ENABLED=true`
- `EMAIL_HOST`, `EMAIL_PORT`
- `EMAIL_USERNAME`, `EMAIL_PASSWORD`
- `EMAIL_FROM`, `EMAIL_FROM_NAME`

Stripe and PayPal channel secrets are stored as payment-channel configuration in the database, not in Git:

- Stripe: `provider_type=official`, `channel_type=stripe`, `interaction_mode=redirect`, `payment_method_types=card`
- PayPal: `provider_type=official`, `channel_type=paypal`, `interaction_mode=redirect`
- Stripe webhook: `https://mmmslai.win/dujiao-next/api/v1/payments/webhook/stripe`
- PayPal webhook: `https://mmmslai.win/dujiao-next/api/v1/payments/webhook/paypal`

Use sandbox/test credentials first. Do not commit SMTP, Stripe, PayPal, database, or SSH credentials.

## Verify

Check the rendered configuration without printing real secrets in logs:

```bash
docker compose --env-file deploy/.env.prod -f deploy/docker-compose.ghcr.yml config --quiet
```

Check runtime status:

```bash
docker compose --env-file deploy/.env.prod -f deploy/docker-compose.ghcr.yml ps
curl -fsS http://127.0.0.1:${SERVER_PORT:-8080}/health
```

## Vercel Frontends

Set the same API base URL in both Vercel projects:

```env
VITE_API_BASE_URL=https://api.your-domain.com
```

For the admin Vercel project, also set `VITE_ADMIN_PATH` only if the admin app is deployed under a subpath. Leave it empty for a normal root deployment.

The frontend clients append `/api/v1` themselves, so do not include `/api/v1` in `VITE_API_BASE_URL`.
