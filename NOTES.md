# dapr-pipeline notes

Log for questions, friction, and comments I captured while I learned learning dapr and wrote the [dapr-pipeline](https://github.com/mchmarny/dapr-pipeline). These notes cover only local development and are limited to the HTTP API.

## tl;dr

* dapr has an intuitive HTTP API (I haven't tried the gRPC one yet) with easy to navigate documentation
* Using dapr framework functionality is easy but the relation of components to dapr's core functionality (e.g. state, topic publishing) is not always clear
* With a few exceptions (e.g. topic sub) dapr HTTP API provides pretty idiomatic developer experience. Language-specific libraries (see [godapr](https://github.com/mchmarny/godapr)) can provide even more natural DX

## notes

> Note not stack ranked, mostly chronological

### Invoke with payload vs string

AFAIK there is no way to pass payload reference in `dapr invoke`. It seems to requires string (e.g. `'{ "f1": "v1", "f2": "v2" }'`). As a developer, I'd like to this `curl` like option:

```shell
dapr invoke --app-id provider \
            --method query \
            --payload "@dir/payload.json"
```

### Components

The components are auto-generated on first run in local dev mode, which is nice, but it isn't always clear what is the relation between individual components/files and the common types of app functionality (e.g. topic publishing and `messagebus` as name). Also, there is the whole `bus` vs `queue` thing ;)

Here is more specific example. The viewer app only subscribes to a topic. Still, whenever I run the `viewer` app I get the `pubsub.yaml` and `statestore.yaml` components generated. Removing them seems to cause dapr to regenerate on the next run.

As a developere I'd like to see...

```shell
dapr components pubsub my-topic \
  --type pubsub.redis
  --metadata redisHost=localhost:6379
```

...or something similar that I'd run explicitly to create an already configured component.

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: my-topic
spec:
  type: pubsub.redis
  metadata:
  - name: redisHost
    value: localhost:6379
  - name: redisPassword
    value: ""
```


### Deep service path

This kind of nested methods (e.g. `/api/v1/query`) invocation in CLI works fine:

```shell
dapr invoke --app-id provider \
            --method /api/v1/query \
            --payload '{ "f1": "v1", "f2": "v2" }'
```

...but `REST` invocation gets messy

```shell
curl -d "@query/simple-query.json" \
     -H "Content-type: application/json" \
     "http://localhost:8000/v1.0/invoke/provider/method/%2Fapi%2Fv1%2Fquery"
```

### Topic Subscription

The method of creating subscription code (requiring `GET` route to `/dapr/subscribe`) is little invasive. I could see some need for dynamically creating new subscriptions so that's probably OK. Still, I would like to see an option to not have to create this route in my handler and use the control plane to create the subscription. For example:

```shell
dapr subscriptions create --topic processed \
                          --app-id viewer \
                          --method /api/v1.0/events
```

I can still subscribe to multiple topics, I do not have to ree-deploy or even re-start the app, and this would keep the developer code entirely idiomatic in situation when the target service wasn't even designed with dapr in mind.

### State store bootstrapping

On initial get to store that has not been yet configured there seems to be some error inconsistencies:

```shell
$: curl -v http://localhost:3500/v1.0/state/producer/qid-1dcd8a3205dfbd5725f2e6e5ec59df28
*   Trying ::1...
* TCP_NODELAY set
* Connection failed
* connect to ::1 port 3500 failed: Connection refused
...
< HTTP/1.1 401 Unauthorized
< Server: fasthttp
< Date: Sat, 11 Apr 2020 12:57:51 GMT
< Content-Type: application/json
< Content-Length: 80
...
* Connection #0 to host localhost left intact
{"errorCode":"ERR_STATE_STORE_NOT_FOUND","message":"state store name: producer"}
* Closing connection 0
```

Also, sometimes, I get `Payment Required` error

```shell
http://localhost:3500/v1.0/state/producer POST: 402 (Payment Required)
```

