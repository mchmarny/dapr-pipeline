# client (simple dapr HTTP client)

dapr has gRPC and REST APIs. For `go`, there is and nice [gRPC SDK](https://github.com/dapr/go-sdk). The REST API is pretty simple to implement but you end up with a lot of redundant code and end up leaking the dapr API throughout application. I create this simple HTTP api to constraint the dapr API to one library and keep all usage idiomatic. 

> Warning, this is a simple library, it does not implements only the most common bits of dapr API.

## Usage

### Create Client 

### State

### Events

### Binding 


