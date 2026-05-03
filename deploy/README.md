# Deployment

This folder contains the runtime compose file for images published by the CI/CD workflows.

## GHCR Compose

Set the required environment variables, then start the stack:

```powershell
$env:IMAGE_OWNER = "pcwto"
$env:IMAGE_TAG = "main"
$env:DJ_DEFAULT_ADMIN_PASSWORD = "change-this-admin-password"
$env:JWT_SECRET = "replace-with-at-least-32-random-characters"
$env:USER_JWT_SECRET = "replace-with-a-different-32-character-secret"
$env:APP_SECRET_KEY = "replace-with-another-32-character-secret"
docker compose -f deploy/docker-compose.ghcr.yml up -d
```

For a tagged release, set `IMAGE_TAG` to the tag name, for example `v1.0.0`.

The three image names are:

- `ghcr.io/<IMAGE_OWNER>/dujiao-next-server:<IMAGE_TAG>`
- `ghcr.io/<IMAGE_OWNER>/dujiao-next-user:<IMAGE_TAG>`
- `ghcr.io/<IMAGE_OWNER>/dujiao-next-admin:<IMAGE_TAG>`
