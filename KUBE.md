# Kubernetes notes

## Kubectl cheatsheet
* `minikube start`
  > Local k8s implementation to test k8s locally
* `minikube stop`
* `minikube delete`
  > Sometimes necessary to clear local state before running start again
* `kubectl get nodes`
* `kubectl get pods`
* `kubectl get deployments`
* `kubectl get services`
* `kubectl scale deployments/<deployment_name> --replicas=<desired_number>`
  > Fun command but bot won't receive this kind of traffic
* `kubectl set image deployments/<deployment_name> <deployment_name>=<new_docker_image_name>`
  > This will perform a rolling update
* `kubectl rollout undo deployments/<deployment_name>`
  > Undo a rolling update back to previous state
* `minikube dashboard`

## Demo run
1. `minikube start`
2. `kubectl create deployment hello-minikube --image=k8s.gcr.io/echoserver:1.10`
  > Create deployment
3. `kubectl expose deployment hello-minikube --type=NodePort --port=8080`
  > Create service using the newly created deployment
4. `kubectl get pod`
  > Verify pod is launched and ready
5. `minikube service hello-minikube --url`
6. View the url to see status page
7. `kubectl delete services hello-minikube`
8. `kubectl delete deployment hello-minikube`
9. `minikube stop`
10. `minikube delete`

## To use local docker images
* `minikube docker-env`
* `eval $(minikube -p minikube docker-env)`

## Bot benefits
* Service recovery if bot goes down
* Rolling deployments (no downtime)

## Clusters
* Single master node
* 1 or many Nodes

## Nodes
* The workers in the cluster
* Runs a Kubelet to communicate with master
* Runs a container runtime

## Pods
* Logical abstraction that can represent a group of one or more containers
* Containers in a pod share the same IP address and ports
* Containers in a pod run in a shared-context
* In a pod you might run a persistent volume app and web server app together

## Deployments
* A configuration file that master will use to maintain state
* This is the self-healing part of k8s
* Node goes down? Machine failure? Master will re-run the config

## Services
* Allows applications to receive traffic
  1. ClusterIP - service is available within the same cluster
  2. NodePort - service is available outside the cluser
  3. LoadBalancer - creates fixed external IP
  4. ExternalName - exposes service by name (CNAME record)
* Services are exposed Deployments meaning traffic can come through
