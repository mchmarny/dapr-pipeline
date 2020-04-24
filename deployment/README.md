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

Now just create a new service to expose the viewer app to external traffic. There are multiple ways to do that in Kubernetes but the simplest way is the expose command with NodePort as parameter. I'll create a proper laod balancer later. 

```shell
kubectl expose deployment/viewer --type="NodePort" --port 8083
```

And then export the dynamically asigned port to the viewer application 

```shell
export VIEWER_PORT=$(kubectl get services/viewer -o go-template='{{(index .spec.ports 0).nodePort}}')
```

## TODO

* Create a service to expose the viewer UI
* Document the expected results of the above commands 


