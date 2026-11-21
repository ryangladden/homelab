# Homelab GitOps

Docker Compose services managed with [Komodo](https://komo.do/) and stored in this repository. Stack definitions target the Komodo server `gertrude` and use the linked repository `homelab`.

## Repository layout

```text
syncs/       Komodo Resource Sync definitions
stacks/      Komodo stack resource definitions (TOML)
builds/      Komodo build resource definitions (TOML)
services/    Compose files and service-specific configuration
```

A stack TOML under `stacks/` generally points Komodo at a Compose project under `services/` using `run_directory`. For example, `stacks/traefik.toml` deploys `services/traefik/compose.yaml`.

The current resource sync is `syncs/homelab.toml`. It links the `homelab` repository and scans the `stacks` directory. The build definitions in `builds/` are **not** included in that sync path; manage or sync those resources separately in Komodo if they are intended to be driven from this repository.

## Making a change

1. Edit the relevant Compose/config file under `services/` and, if needed, its stack definition under `stacks/`.
2. For a new Komodo-managed stack, add a stack resource TOML in `stacks/` and set its `linked_repo`, `server`, and `run_directory` consistently with the existing definitions.
3. Commit and push the change to the linked Git repository.
4. In Komodo, refresh/sync the `homelab-stacks` Resource Sync and review the resulting resource changes before applying them. Deploy the relevant stack and confirm its health and logs.

Treat production deployments as stateful operations: review image/tag changes, migrations, and data-volume impact before deploying. A Git revert restores configuration, but does not necessarily roll back application data or database migrations.

## Configuration and secrets

Stack-level environment values are declared in each stack TOML's `[stack.config].environment` multiline string. Values such as `[[HOST_DOMAIN]]` and `[[POSTGRES_PASSWORD]]` are Komodo variable substitutions; they are not Compose syntax. Define the referenced variables in the appropriate Komodo variable scope before syncing or deploying. Secret values should be stored as Komodo secrets/variables, not committed to Git.

Compose files also use `${VARIABLE}` interpolation. These values are supplied to Compose through the stack environment (and, for some services, a generated or separately managed `.env` file). When adding a variable, verify that it is defined in the right place and that required values cannot silently resolve to empty strings. Do not commit live `.env` files, credentials, tokens, private keys, or generated secrets.

## Host requirements

The Compose projects assume a prepared Docker host. Before deployment, check the relevant service configuration for:

- External Docker networks such as `proxy`, `docker-socket-proxy`, `radicale-internal`, and `egress`; Compose will not create networks declared `external: true`.
- Host bind-mount directories under paths such as `/opt/data`, `/opt/state`, and `/mnt`, including ownership, permissions, available space, and backup coverage.
- Hardware dependencies such as GPU devices, `/dev` mappings, or Docker GPU support.
- Required DNS, firewall, and ingress configuration for published ports and Traefik routes.

Requirements vary by stack; consult its Compose file and stack environment before deploying it to a different host.

## Validation and operations

Before applying changes, validate the relevant Compose project with Docker Compose using the same environment values Komodo will provide. For example, from a service directory:

```sh
docker compose config --quiet
```

This requires the stack's interpolation variables to be available locally. Be careful when rendering full Compose output: it can contain interpolated secrets. Never paste that output into an issue or commit it.

After deployment, check the Komodo stack status, container health, and logs. Back up persistent data and service-specific configuration—not just this Git repository—and periodically test that backups can be restored. Keep a record of host prerequisites and recovery steps for critical services.

## Adding a service

1. Create `services/<service>/compose.yaml` and keep service-specific configuration alongside it.
2. Add `stacks/<service>.toml` with a unique Komodo resource name and a `run_directory` of `services/<service>`.
3. Add required non-secret settings and Komodo variable references to the stack environment; provision referenced variables in Komodo.
4. Add appropriate health checks, restart behavior, network boundaries, persistent storage, and documented backup requirements.
5. Check image update policy, port exposure, and any privileged host access (especially Docker socket access) before deployment.
6. Sync and review the resource in Komodo, then deploy and verify it.
