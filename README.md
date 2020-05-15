# freebox_exporter

A Prometheus exporter for Freebox stats

## Cmds

`freebox_exporter`

## Flags

- `-endpoint`: Freebox API url (default http://mafreebox.freebox.fr)
- `-listen`: port for Prometheus metrics (default :10001)
- `-debug`: turn on debug mode
- `-fiber`: turn off DSL metric for fiber Freebox

## Preview

Here's what you can get in Prometheus / Grafana with freebox_exporter:

![Preview](https://user-images.githubusercontent.com/13923756/54585380-33318800-4a1a-11e9-8e9d-e434f275755c.png)

# How to use it

## Compiled binary

If you want to compile the binary, you can refer to [this document](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63) which explains how to do it, depending on your OS and architecture. Alternatively, you can use `./build.sh`.

You can also find the compiled binaries for MacOS, Linux (x86_64, arm64 and arm) and Windows in the release tab.

### Quick start

```
./freebox_exporter
```

### The following parameters are optional and can be overridden:

- Freebox API endpoint

```
./freebox_exporter -endpoint "http://mafreebox.freebox.fr"
```

- Port

```
./freebox_exporter -listen ":10001"
```

## Docker

### Quick start

```
docker run -d --name freebox-exporter --restart on-failure  -p 10001:10001 \
  saphoooo/freebox-exporter
```

### The following parameters are optional and can be overridden:

- Local token

Volume allows to save the access token outside of the container to reuse authentication upon an update of the container.

```
docker run -d --name freebox-exporter --restart on-failure  -p 10001:10001 \
  -e HOME=token -v /path/to/token:/token saphoooo/freebox-exporter
```

- Freebox API endpoint

```
docker run -d --name freebox-exporter --restart on-failure -p 10001:10001
  saphoooo/freebox-exporter -endpoint "http://mafreebox.freebox.fr"
```

- Port

```
docker run -d --name freebox-exporter --restart on-failure -p 8080:10001 \
  saphoooo/freebox-exporter
```

## Caution on first run

If you launch the application for the first time, you must allow it to access the freebox API.
- The application must be launched from the local network.
- You have to authorize the application from the freebox front panel.
- You have to modify the rights of the application to give it "Modification des réglages de la Freebox"

Source: https://dev.freebox.fr/sdk/os/
