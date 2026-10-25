# vela-worker

Helm chart for the [Vela](https://go-vela.github.io/docs/) worker.

```shell
helm install vela-worker ./charts/vela-worker \
  --set server.addr=http://vela-server:8080 \
  --set server.secret=<shared-secret>
```

## Runtimes

- **kubernetes** (default): every build runs as a pod in `runtime.namespace` (default: the release namespace). The chart creates a ServiceAccount plus a Role and RoleBinding in that namespace. The kubernetes runtime needs a worker that includes the kubernetes runtime fixes from this branch.
- **docker**: builds run through the node's Docker socket (`runtime.docker.socketPath`), so worker pods need nodes that run Docker.

## Worker addresses

The worker runs as a StatefulSet behind a headless Service. Each replica registers with the server as `http://<pod>.<release>.<namespace>.svc:8080`, so its address stays the same across restarts and the server doesn't build up stale worker records.

## Server authentication

Set `server.secret` (or `server.existingSecret` / `server.existingSecretKey`) to use the shared `VELA_SERVER_SECRET`. Leave both empty to register each worker with a registration token instead.

## Pipeline pods template

The `PipelinePodsTemplate` CRD is installed from `crds/`. Set `runtime.podsTemplate.create=true` and `runtime.podsTemplate.template` to have the chart create one, or `runtime.podsTemplate.name` to use an existing template.

## Values

See [values.yaml](values.yaml). Anything not exposed there can be set with `extraEnv`.
