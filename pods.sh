#!/bin/bash
# first arguement defines the MAX, with a default if not set

MAX="${1:-9}"

TEMPLATE=$( cat <<-END
apiVersion: v1
kind: Pod
metadata:
  annotations:
    "scheduler.alpha.kubernetes.io/name": priority
    "k8s_priority": "PRIORITY"
  generateName: pausePRIORITY-
  labels:
    app: pause
    priority: "PRIORITY"
spec:
  containers:
  - image: gcr.io/google_containers/pause:2.0
    name: pause
    resources:
      requests:
        cpu: 400m
        memory: 100Mi
END
)

for n in $(seq 1 $MAX); do
    echo "$TEMPLATE" | sed "s/PRIORITY/$n/g" | kubectl create -f - --validate=false
done
