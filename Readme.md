# build docker
docker build -t zokypesch/go_example .

# login to docker
docker login

username: zokypesch
password: sepertibiasa

# push to cloud 

docker push zokypesch/go_example

# run docker local

docker ps -a | grep "zokypesch/go_example"
docker stop b15222b1753e
docker run -d -it -p 8000:8000  zokypesch/go_example

# create cluster
gcloud container clusters create belajar-ci --zone asia-east1-a
# delete cluster
gcloud container clusters delete belajar-ci --zone asia-east1-a


# deploy
kubectl apply -f kube/namespace.yml
kubectl apply -f kube/deployment.yml
kubectl apply -f kube/ingress.yml

# get deployment, svc
kubectl -n staging get deployment,pod,svc,endpoints,pvc,ingress
kubectl describe ingress goexample -n staging
kubectl describe pod <pod> -n staging

# see ip public
gcloud container clusters list
35.229.221.221

# get ingress
kubectl get ingress goexample -n staging

# go a log
kubectl logs goexample-559c8bd7cd-4x6b5 -n staging 

kubectl set image deployment goexample goexample=zokypesch/go_example:latest --record -n staging

# to get local cluser dns 
kubectl exec goexample-5c6445974-4j7v9 cat /etc/resolv.conf -n staging
 result: staging.svc.cluster.local

 if you want to connect directly please add service (see a deployment.yml) ex: goexample 
 so final: goexample.staging.svc.cluster.local

# exec 
kubectl exec -it goexample-5c6445974-k467h -n staging -- /bin/bash

exit
for exit exec
