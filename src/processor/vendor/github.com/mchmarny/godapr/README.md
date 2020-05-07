# godapr (simple dapr HTTP client)

dapr has gRPC and REST APIs. For `go`, there is the auto-generated [gRPC SDK](https://github.com/dapr/go-sdk) that covers the complete spectrum of dapr API. Developers can also implement their own HTTP calls to the REST API. When invoking the dapr REST APIs there usually is lot's of redundant code builting request and aprsing response. I create this simple HTTP api to constraint the dapr API to one library.

> Warning, this library implements only the most common parts of dapr API (state, pubsub, and binding). 

## Usage

To use `godapr` first get the library

```shell
go get github.com/mchmarny/godapr
```

### Create Client

To use `godapr` library in your code, first import it

```go
import dapr "github.com/mchmarny/godapr"
```

Then create a `godapr` client with the `dapr` server defaults

```go
client := dapr.NewClient()
```

or if you need to specify non-default dapr port

```go
client := dapr.NewClientWithURL("http://localhost:3500")
```

> consider getting the dapr server URL from environment variable

### State

#### Get Data

To get state data you can either use the client defaults ("strong" Consistency, "last-write" Concurrency)

```go
data, err := client.GetData("store-name", "record-key")
```

Or define your own state options

```go
opt := &StateOptions{
    Consistency: "eventual",
    Concurrency: "first-write",
}

data, err := client.GetDataWithOptions("store-name", "record-key", opt)
```

#### Save Data

Similarly with saving state, assuming you have your own person object for example

```go
person := &Person{
    Name: "Example John",
    Age: 35,
}
```

you can either use the defaults

```go
err := client.SaveData("store-name", "record-key", person)
```

Or define your own state data object

```go
data := &StateData{
    Key: "id-123",
    Value: person,
    Options: &StateOptions{
        Consistency: "eventual",
        Concurrency: "first-write",
    },
}

err := client.Save("store-name", data)
```

### Events

To publish events to a topic you can just

```go
err := client.Publish("topic-name", person)
```

### Binding

Similarly with binding you can use the `Send` method

```go
err := client.Send("binding-name", person)
```


## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [Apache v2 License](./LICENSE)
