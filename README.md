This project is used to been a supervisor, which can supervise OS provision progress.

# Design Concept

Main task is to resolve one issue.
Auto detect when the cluster was created succeeded by any OS provision tools.

# How to build
## Load package

### Get dep tool

```sh
go get -u github.com/golang/dep/cmd/dep
```

### Update package

```sh
dep ensure
```

### Build opps package

```sh
make build_linux
```
