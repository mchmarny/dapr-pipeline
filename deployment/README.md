# Kubernetes Deployment

Assuming you have `kubectl` installed and configure to connect to your cluster you will need to setup the necessary secrets:

> Note, dapr supports wide array of state and pubsub backing services across multiple Cloud and on-prem deployments. This document will use [Azure Table Storage](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-create?tabs=azure-portal) for state, and [Azure Service Bus](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-quickstart-topics-subscriptions-portal) for pubsub but you can easily substitute these using any of the components listed [here](https://github.com/dapr/docs/tree/master/howto).

## Component-backing services 

### Azure Table Storage

To set up Azure Table Storage itself follow the instructions [here](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-create?tabs=azure-portal)

```shell
kubectl create secret generic pipeline-state \
  --from-literal=account-name='' \
  --from-literal=account-key=''
```

> TODO: add expected return from command and way to validate 

### Azure Service Bus

To set up Azure Service Bus itself follow the instructions [here](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-quickstart-topics-subscriptions-portal)


```shell
kubectl create secret generic pipeline-bus \
  --from-literal=connection-string=''
```

> TODO: add expected return from command and way to validate 


### Deploy components

```shell
kubectl apply -f deployment/components
```

> TODO: add expected return from command and way to validate 


## Deploy pipeline 

Before deploying the actual pipeline you will have to create a secret to enable the `producer` to query Twitter API. You can get these by registering a Twitter application [here](https://developer.twitter.com/en/apps/create).


```shell
kubectl create secret generic pipeline-twitter \
  --from-literal=access-secret: '' \
  --from-literal=access-token: '' \
  --from-literal=consumer-key: '' \
  --from-literal=consumer-secret: ''
```

> TODO: add expected return from command and way to validate 

One the `pipeline-twitter` twitter is created, you are ready to deploy the entire pipeline (`producer`, `processor`, `viewer`

```shell
kubectl apply -f deployment/
```

> TODO: add expected return from command and way to validate 

### Exposign viewer UI

To expose the viewer application extertnally, create a `service` and an `ingress` by applying the `route.yaml`

```shell
kubectl apply -f viewer-route.yaml
```

> Note, you will have to change the ingress host rule to DNS you can actually control. I manage `things.io` so in this case I created an `A` record to point to the ingress IP. 

```yaml
rules:
  - host: dapr.thingz.io
 ```

You can find the IP address assigned to the viewer ingress on your cluster using:

`kubectl get ingress viewer`

Now you should be able to access the demo UI using the DNS defined in your `ingress` (e.g. dapr.thingz.io)

