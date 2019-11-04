# Documentation

For installation and usage, please [visit our docs](https://go-vela.github.io/docs).

**If you haven't already, please see the [Vela server documentation](https://github.com/go-vela/server/blob/master/.github/DOCS.md) to create the services necessary for executing builds locally.**

## Services

If you followed [the instructions from the contributing guide](CONTRIBUTING.md/#getting-started), you should have 2 services running as Docker containers on your machine:

* worker_1
* worker_2

### Worker

The `worker_1` and `worker_2` services are running the exact same code. The [docker-compose](../docker-compose.yml) file is already setup to connect will the other services (`redis` and `postgres`) as well as the Vela server service you can create from the [getting started section](CONTRIBUTING.md/#getting-started).

### Vault

The `vault` service hosts the Vault store used for integration with one of Vela's secret engine implementations.

## API

Coming soon!

## CLI

Coming soon!

## Executing Builds

In order to execute builds on your local machine, you'll also need to create a Vela server to push the workloads to the `redis` queue.

To create the server, you can follow the  [documentation](https://github.com/go-vela/server/blob/master/.github/DOCS.md) found in the [go-vela/server](https://github.com/go-vela/server) repository.
