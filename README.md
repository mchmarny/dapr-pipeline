# dapr-tweet-processing-pipeline

Example of Twitter event processing pipeline


## Questions/Friction/Comments

### Invoke with payload vs string

AFAIK there is no way to pass payload reference, requires string (e.g. `'{ "f1": "v1", "f2": "v2" }'`). Like to do this:

```shell
dapr invoke --app-id provider \
            --method query \
            --payload "@dir/payload.json"
```

### Deep service path

This kind of nested methods (e.g. `/v1/query`) works fine:

```shell
dapr invoke --app-id provider \
            --method /v1/query \
            --payload '{ "f1": "v1", "f2": "v2" }'
```

But `curl` invocation gets messy

```shell
curl -d "@query/simple-query.json" \
     -H "Content-type: application/json" \
     "http://localhost:8000/v1.0/invoke/provider/method/%2Fv1%2Fquery"
```

### State store bootstrapping

On initial get to store that has not been yet configured there seems to be some error inconsistency

```shell
$: curl -v http://localhost:3500/v1.0/state/producer/qid-1dcd8a3205dfbd5725f2e6e5ec59df28
*   Trying ::1...
* TCP_NODELAY set
* Connection failed
* connect to ::1 port 3500 failed: Connection refused
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 3500 (#0)
> GET /v1.0/state/producer/qid-1dcd8a3205dfbd5725f2e6e5ec59df28 HTTP/1.1
> Host: localhost:3500
> User-Agent: curl/7.64.1
> Accept: */*
>
< HTTP/1.1 401 Unauthorized
< Server: fasthttp
< Date: Sat, 11 Apr 2020 12:57:51 GMT
< Content-Type: application/json
< Content-Length: 80
<
* Connection #0 to host localhost left intact
{"errorCode":"ERR_STATE_STORE_NOT_FOUND","message":"state store name: producer"}* Closing connection 0
```

Also, on `POST` I sometimes get `Payment Required` error?

```shell
http://localhost:3500/v1.0/state/producer POST: 402 (Payment Required)
```

### State data

The [doc](https://github.com/dapr/docs/tree/master/howto/stateful-replicated-service) doesn't seem clear. What's the thinking behind state data on post being an array? Given the value is an arbitrary type users can set set it to array and avoid needing to create and parse slices on the more common use-case of single item.