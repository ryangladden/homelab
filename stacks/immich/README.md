# Environment Variables

This document lists the environment variables required to configure each service.

## Immich

Configure these variables in the Immich `.env` file.

| Variable                   |     Required     | Default / Example                           | Description                                              |
| -------------------------- | :--------------: | ------------------------------------------- | -------------------------------------------------------- |
| `UPLOAD_LOCATION`          |        Yes       | `/opt/container-data/immich/mount/library`  | Location where the Immich photo/video library is stored. |
| `DB_DATA_LOCATION`         |        Yes       | `/opt/container-data/immich/postgres`       | Location where PostgreSQL data is stored.                |
| `THUMB_LOCATION`           |        Yes       | `/opt/container-data/immich/thumbs`         | Location where generated thumbnails are stored.          |
| `ENCODED_VIDEO_LOCATION`   |        Yes       | `/opt/container-data/immich/encoded-videos` | Location where transcoded/encoded videos are stored.     |
| `EXTERNAL_LIB_LOCATION`    |        Yes       | `/opt/container-data/immich/mount`          | Host location exposed to Immich for external libraries.  |
| `TZ`                       |        No        | `Etc/UTC`                                   | Timezone using a valid TZ database identifier.           |
| `IMMICH_VERSION`           |        Yes       | `release`                                   | Immich Docker image version/tag to deploy.               |
| `DB_PASSWORD`              |      **Yes**     | —                                           | PostgreSQL password. **Must be set.**                    |
| `IMMICH_API_KEY`           | Power Tools only | —                                           | Immich API key used by Immich Power Tools.               |
| `IMMICH_TELEMETRY_INCLUDE` |        No        | `all`                                       | Optional Immich telemetry configuration.                 |


### Minimum Required Configuration

At minimum, make sure `DB_PASSWORD` is populated before starting the stack:

```env
DB_PASSWORD=<your-secure-database-password>
```

If Immich Power Tools is enabled, also provide:

```env
IMMICH_API_KEY=<your-immich-api-key>
```

