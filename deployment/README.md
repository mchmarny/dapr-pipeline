# Kubernetes Deployment

This document will overview the `dapr-pipeline` demo deployment into Kubernetes. For illustration purposes, all commands in this document will based on Microsoft Azure. dapr supports wide array of state and pubsub backing services across multiple Cloud and on-prem deployments. So if you have a Kubernates cluster soemwhere else, you can substitute:

* [Azure Table Storage](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-create?tabs=azure-portal) as state backing service with anyone of [these](https://github.com/dapr/docs/tree/master/howto/setup-state-store)
* [Azure Service Bus](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-quickstart-topics-subscriptions-portal) as pubsub backing service with anyone of [these](https://github.com/dapr/docs/tree/master/howto/setup-pub-sub-message-broker) 

## Prerequisite

* [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest)

## Configuration

Also, to simplify all the scripts in this doc, set a few `az` defaults:

```shell
az account set --subscription <name or id>
az configure --defaults location=<location> group=<your resource group>
```

## Cluster Setup (optional)

If you don't already have one, you can create Kubernates cluster on Azure with all the necessary add-ons usign tihs command:

```shell
az aks create --name daprdemo \
              --kubernetes-version 1.15.10 \
              --enable-managed-identity \
              --vm-set-type VirtualMachineScaleSets \
              --node-vm-size Standard_F4s_v2 \
              --enable-addons monitoring,http_application_routing        \
              --generate-ssh-keys
```

## Component-backing services 

Assuming you have a Kubernates cluster and `kubectl` CLI configure to connect you will need to setup the `dapr` components and their backing services:


### State

To configure `dapr` state component we will use Azure Table Storage. To set it up you can follow [these instructions](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-create?tabs=azure-portal). Once finished, you will need to cofigure the Kubernates secrets to hold the Azure Table Storage account information:

```shell
kubectl create secret generic pipeline-state \
  --from-literal=account-name='' \
  --from-literal=account-key=''
```

To deploy the `dapr` state component configured for the above set up service

```shell
kubectl apply -f component/state.yaml
```

### PubSub

To configure `dapr` pubsub component we will use Azure Service Bus. To set it up you can follow [these instructions](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-quickstart-topics-subscriptions-portal). Once finished, you will need to configure the Kubernates secrets for Azure Service Bus connection string information. 


```shell
kubectl create secret generic pipeline-bus \
  --from-literal=connection-string=''
```

To deploy the `dapr` pubsub topic components for the above set up service

```shell
kubectl apply -f component/processed.yaml -f component/tweet.yaml
```

### Binding 

To configure `dapr` binding component we will use a simple service offered by thingspeak.com that does not require any additional configuration. 

```shell
kubectl apply -f component/alert.yaml
```

## Deploy pipeline 

Before deploying the actual pipeline you will have to create one more secret, the Twitter API secrets for `producer`. You can get these by registering a Twitter application [here](https://developer.twitter.com/en/apps/create).


```shell
kubectl create secret generic pipeline-twitter \
  --from-literal=access-secret: '' \
  --from-literal=access-token: '' \
  --from-literal=consumer-key: '' \
  --from-literal=consumer-secret: ''
```

Once the `pipeline-twitter` secret is created, you are ready to deploy the entire pipeline (`producer`, `processor`, `viewer`

```shell
kubectl apply -f producer.yaml -f processor.yaml -f viewer.yaml
```

### Exposign viewer UI

To expose the viewer application extertnally, create you will need to create Kubernetes `service` and `ingress` by applying the [route.yaml](./viewer-route.yaml)

```shell
kubectl apply -f viewer-route.yaml
```

> Note, you will have to change the ingress host rule to DNS you can actually control. I own `things.io` so in this case I created an `A` record to point to the ingress IP. 

```yaml
rules:
  - host: dapr.thingz.io
 ```

You can find the IP address assigned to the viewer ingress on your cluster using:

`kubectl get ingress viewer`

Now you should be able to access the demo UI using the DNS defined in your `ingress` (e.g. dapr.thingz.io)

## Invoking query

To submit query similar to the way described in the local developemnt demo, you will have to forward the local port to the `producer-dapr` service.

```shell
kubectl port-forward service/producer-dapr 8080:80
```

> Exposign the producer service externally like we did with the viewer is not recomanded as that would enable anyone in the world to submit queries and use your Twitter API credits

Once forwarded, you can execute queries like this: 

```shell
curl -d '{ "query": "serverless OR faas OR dapr", "lang": "en" }' \
     -H "Content-type: application/json" \
     "http://localhost:8080/v1.0/invoke/producer/method/query"
```

If everything went OK, you should see the tweets with sentiment score appear on the UI. 
