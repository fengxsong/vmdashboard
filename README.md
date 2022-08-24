# vmdashboard

Super light wight Kubernetes APIserver proxy. Since nodejs frontend service has server-side issue, inspired by [tekton dashboard](https://github.com/tektoncd/dashboard), we create a proxy service to wrap HTTP calls to APIserver. In the meantime, builtin authn/authz mechanisms are keep, and provide extension interface with external auth systems.

> Most codebase are from official [kubernetes](https://github.com/kubernetes/kubernetes) repository.

## Develop

```bash
git clone https://github.com/fengxsong/vmdashboard
cd vmdashboard
make tidy
```

## Build

```bash
# Binary
make bin
# docker image
make docker-build
```

## Usage

```bash
make bin
# print flags
./build/_output/bin/vmdashboard --help
```

## TODO

## Contribution

Issues and contributions are welcome!
