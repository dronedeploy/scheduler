# scheduler

Demos the ability to have a priority queue scheduler for kubernetes

## Usage

install and run minikube (or other kubernetes cluster, must have no auth for this demo to work)
```
minikube start
```
Setup a proxy connection to your cluster
```
> kubectl proxy
Starting to serve on 127.0.0.1:8001
```

### Create some properly labeled pods

```
bash pods.sh 20
```

The pods should be in a pending state:

```
kubectl get pods
```

### Build the scheduler

Run the build script
```
bash build
```
there is also a dockerized version

### Run the priority updater
```
cd prioritizer && bash build && ./prioritizer
```

### Run the Scheduler

Run the priority queue scheduler:

```
scheduler
```
```
2016/08/19 11:16:25 Starting custom scheduler...
... does stuff
2016/08/19 11:16:35 Shutdown signal received, exiting...
2016/08/19 11:16:35 Stopped reconciliation loop.
2016/08/19 11:16:35 Stopped scheduler.
```

> Notice the highest priority pods are scheduled first

## Run the Scheduler on Kubernetes

```
kubectl create -f deployments/scheduler.yaml
```
``` 
deployment "scheduler" created
```
