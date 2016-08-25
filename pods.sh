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
  generateName: pause-
  labels:
    app: pause
    priority: "PRIORITY"
    valueMult: "VALUE"
    expectedDuration: "DURATION"
spec:
  containers:
  - image: gcr.io/google_containers/pause:2.0
    name: pause
    resources:
      requests:
        cpu: 300m
        memory: 100Mi
END
)

for n in $(seq 1 $MAX); do
    VALUE="$(( ( RANDOM % 10 )  + 1 ))"
    DURATION="$(( ( RANDOM % 1000 )  + 1 ))"
    echo "$TEMPLATE" | sed "s/PRIORITY/$n/g" | sed "s/VALUE/$VALUE/g" | sed "s/DURATION/$DURATION/g" | kubectl create -f - --validate=false
done
