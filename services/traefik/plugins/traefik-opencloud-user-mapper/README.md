# Traefik OpenCloud User Mapper

Traefik middleware that maps a trusted authenticated username (for example TinyAuth's `Remote-User`) to the corresponding OpenCloud user UUID using the OpenCloud Graph API, then overwrites `X-Remote-User` before forwarding the request.

## Security model

Use this middleware **after** your authentication middleware. It fails closed if the authenticated username is absent, the OpenCloud user cannot be resolved, or the OpenCloud API is unavailable. Any inbound `X-Remote-User` is deleted and replaced with the OpenCloud UUID.

## Traefik environment

```yaml
environment:
  OC_MAPPER_USERNAME: ${OC_MAPPER_USERNAME}
  OC_MAPPER_APP_TOKEN: ${OC_MAPPER_APP_TOKEN}
```

The App Token may contain spaces; Docker Compose environment values preserve them.

## Local plugin installation

Static `traefik.yaml`:

```yaml
experimental:
  localPlugins:
    opencloud-user-mapper:
      moduleName: github.com/ryangladden/traefik-opencloud-user-mapper
```

Mount this repository at:

```text
/plugins-local/src/github.com/ryangladden/traefik-opencloud-user-mapper
```

Example Compose mount:

```yaml
volumes:
  - /opt/traefik/plugins/traefik-opencloud-user-mapper:/plugins-local/src/github.com/ryangladden/traefik-opencloud-user-mapper:ro
```

Dynamic configuration:

```yaml
http:
  middlewares:
    tinyauth:
      forwardAuth:
        address: http://tinyauth:3000/api/auth/traefik
        authResponseHeaders:
          - Remote-User

    opencloud-user-mapper:
      plugin:
        opencloud-user-mapper:
          opencloudURL: http://opencloud:9200
          usernameHeader: Remote-User
          outputHeader: X-Remote-User
          serviceUsernameEnv: OC_MAPPER_USERNAME
          appTokenEnv: OC_MAPPER_APP_TOKEN
          cacheTTL: 1h
          requestTimeout: 3s
```

Apply in this order:

```yaml
traefik.http.routers.radicale.middlewares: tinyauth@file,opencloud-user-mapper@file
```

Radicale remains:

```ini
[auth]
type = http_x_remote_user
```

## OpenCloud URL

If Traefik and OpenCloud share a Docker network, prefer the internal URL (for example `http://opencloud:9200`) to avoid a public round trip. The Graph lookup is:

```text
GET /graph/v1.0/users/{Remote-User}
```

using the configured OpenCloud service username and App Token via HTTP Basic authentication.
