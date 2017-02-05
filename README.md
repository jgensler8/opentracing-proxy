# OpenTracing Proxy

Proxy Service to determine dependencies from services that can't be instrumented or instrumented in the future.

It is a piece of middleware that can start spans for new requests and continue spans for existing requests.

## Example

Either compile and run locally or compile and run in a docker container.

```bash
GOOS=linux go build
docker-compose up
```

```bash
# Proxy running locally OR in container
curl --proxy http://localhost:9999 http://localhost:9090/hello/world
# Proxy running in container
curl --proxy http://localhost:9999 http://echoheaders:9090/hello/world 
```


```bash
# Open the Zipkin UI
open http://localhost:9411
```

![This is the "trace" view](/pictures/zipkin-ui-example.png)