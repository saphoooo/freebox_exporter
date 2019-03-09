# freebox_exporter
A prometheus exporter for freebox stats

# how to use it

## Standalone

```
./freebox_exporter -version "v6" -endpoint "http://mafreebox.freebox.fr -listen ":10001"
```

## Docker 
  
```
docker run -d --name freebox-exporter --restart always -p 10001:10001 saphoooo/freebox-exporter
```

## flags
- `-version`: freebox API version (default v6)
- `-endpoint`: freebox API url (default http://mafreebox.freebox.fr)
- `-listen`: port for prometheus metrics (default :10001)

## first run
If you launch the application for the first time, you must allow it to access the freebox API.
- The application must be launched from the local network.
- You have to authorize the application from the freebox front panel.
- You have to modify the rights of the application to give it "Modification des réglages de la Freebox"
  
Source: https://dev.freebox.fr/sdk/os/login/
