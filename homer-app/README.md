## Deploy Homer-app

helm -n homer-app install homer-app djjudas21/homer -f homer-app/homer-helm-values.yaml --create-namespaces 

kubectl -n homer-app get pods

kubectl -n homer-app apply -f homer-app/ingress-route.yaml 

kubectl -n homer-app apply -f rancher-master/rancher-tls-secret.yaml 

