# freebox_exporter
A prometheus exporter for freebox stats

## flags
- `-version`: freebox API version (default v6)
- `-endpoint`: freebox API url (default http://mafreebox.freebox.fr)
- `-listen`: port for prometheus metrics (default :10001)

# how to use it

## Compiled binary

### Quickstart

```
./freebox_exporter
```

### The following parameters are optional and can be superseed:

- Freebox API version

```
./freebox_exporter -version "v6"
```

- Freebox API endpoint

```
./freebox_exporter -endpoint "http://mafreebox.freebox.fr
```

- Port

```
./freebox_exporter -listen ":10001"
```

## Docker 

### Quickstart
  
```
docker run -d --name freebox-exporter --restart on-failure  -p 10001:10001 \
  saphoooo/freebox-exporter
```

### The following parameters are optional and can be superseed:

- Local token

```
docker run -d --name freebox-exporter --restart on-failure  -p 10001:10001 \
  -e HOME=token -v /path/to/token:/token saphoooo/freebox-exporter
```

- Freebox API version

```
docker run -d --name freebox-exporter --restart on-failure  -p 10001:10001 \
  saphoooo/freebox-exporter -version "v6"
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
  
Source: https://dev.freebox.fr/sdk/os/login/
