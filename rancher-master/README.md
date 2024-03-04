### Rancher Local cluster Deployment

Please ensure and update the secrets with appropriate values. I have used dummy values for obvious reasons.

  1 - SSH to the first node and Install k3s

      $ curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.26.7+k3s1 INSTALL_K3S_EXEC="server" sh -
    
  2 - Once installed, copy the k3s.yaml to interact with the cluster
      
      $ cat /etc/rancher/k3s/k3s.yaml

  3 - Paste the k3s.yaml contents in your local system and export as KUBECONFIG. If it is pointing to localhost:6443 you might need to modify `clusters.cluster.server` url. Use the correct Server URL.

  4 - Now run the below Steps to deploy cert-manager and rancher:

      $ kubectl create ns cattle-system
      $ helm install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.13.1 --set installCRDs=true
      $ helm install rancher rancher-latest/rancher -n cattle-system -f rancher-master/rancher-helm-values.yaml --create-namespace
      $ kubectl -n cattle-system get certificaterequest
      $ kubectl -n cattle-system get issuer
      $ kubectl -n cattle-system get certificate
      $ kubectl -n cattle-system get order
      $ kubectl -n cattle-system get secret
      $ kubectl -n cattle-system get pods

  At this point you should have Rancher running with certmanager managing the rancher generated certificates. Get the login details and confirm if everything is working as expected.

      $ kubectl get secret --namespace cattle-system bootstrap-secret -o go-template='{{.data.bootstrapPassword|base64decode}}'
      $ echo https://rancher.rajesh-kumar.in/dashboard/?setup=$(kubectl get secret --namespace cattle-system bootstrap-secret -o go-template='{{.data.bootstrapPassword|base64decode}}')

  5 - Create ingress resource to access the Rancher.

      $ kubectl apply -f ingresses/rancher-ingress.yaml    

  I have hosted this on personal domain hence I created a wildcard certificate for my domain. I also created a secret to store the certificates.

  6 - Let's store the secret as YAML as I have used the same certificate to host all the applications.
      
      $ kubectl -n cert-manager get secret rancher-tls-cert -o yaml | yq 'del(.metadata.creationTimestamp, .metadata.resourceVersion, .metadata.selfLink, .metadata.uid, .metadata.managedFields, .metadata.namespace,   .metadata.annotations)' | kubectl apply -n cattle-system -f -

  7 - I have used Vault to store all my secrets and have used external-secrets to fetch the secrets.

      $ helm repo add hashicorp https://helm.releases.hashicorp.com
      $ helm install vault hashicorp/vault
      $ kubectl exec -it vault-0 vault operator init // COPY the tokens and store.
      $ kubectl exec -it vault-0 -- vault login $INITIAL_ROOT_TOKEN
      $ kubectl apply -f vaults/vaultTokenSecret.yaml
      $ kubectl apply -f ingresses/vault-ingress.yaml
      $ kubectl apply -f clusterSecretStore.yaml


NOTE: rancher-master/rancher-tls-secret.yaml needs to be created for each namespace and on each cluster where ingressroute is created. This file isn't part of git repo for security issues.


Adding New remote cluster

1 - Deploy kubernetes using one of k3s, rke2 or rke

RKE: 

https://rke.docs.rancher.com/installation

Use rancher-master/cluster.yml to deploy an rke cluster make changes as required, or you can run rke config to create a cluster.yaml

K3s:

https://docs.k3s.io/installation/configuration#configuration-with-install-script
```SHELL
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.26.7+k3s1 INSTALL_K3S_VERSION=v1.26.7+k3s1 INSTALL_K3S_EXEC="server" sh -
```
2 - Ensure traefik is deployed. If not deploy it. k3s deploys it by default while rke doesn't in the case deploy traefik.

`For RKE:`
```SHELL
helm upgrade --install traefik traefik/traefik -n kube-system -f applications/home-cluster-rke/traefik-config.yaml
```
`For k3s:`
SSh to the master node and create /var/lib/rancher/k3s/server/manifests/traefik-config.yaml

```SHELL
apiVersion: helm.cattle.io/v1
kind: HelmChartConfig
metadata:
  name: traefik
  namespace: kube-system
spec:
  valuesContent: |-
    additionalArguments:
    - "--serversTransport.insecureSkipVerify=true"
    deployment:
      kind: DaemonSet
    logs:
      access:
        enabled: true
```

3 - Deploy external-secret

```SHELL 
 helm install external-secrets external-secrets/external-secrets -n external-secrets --create-namespace
```

4 - Create vault secret, clustersecretstore
 
```SHELL 
kubectl apply -f rancher-master/vaults/vaultTokenSecret.yaml
kubectl apply -f rancher-master/clusterSecretStore.yaml
 ```

5 - Create rancher-tls-secret in app namespaces

 ```SHELL
 kubectl -n app-ns apply -f rancher-master/rancher-tls-secret.yaml
 ```
6 -  Expose application with ingressroute and use rancher-tls-secret for tls secret.

7 - Ensure everything is working

```SHELL
kubectl get clustersecretstore
kubectl get externalsecret -A
kubectl get ingressroute -A
```
