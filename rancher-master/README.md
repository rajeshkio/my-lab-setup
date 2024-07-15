# 🚀 Rancher Local Cluster Deployment

Ensure to update the secrets with appropriate values. Dummy values are used here for demonstration purposes.

## Step 1: SSH to the First Node and Install K3s

```sh
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.26.7+k3s1 INSTALL_K3S_EXEC="server" sh -
```

## Step 2: Copy the `k3s.yaml` to Interact with the Cluster

```sh
cat /etc/rancher/k3s/k3s.yaml
```

## Step 3: Configure `KUBECONFIG`

Paste the `k3s.yaml` contents in your local system and export it as `KUBECONFIG`. If it points to `localhost:6443`, modify the `clusters.cluster.server` URL to the correct server URL.

## Step 4: Deploy Cert-Manager and Rancher

```sh
kubectl create ns cattle-system
helm install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.13.1 --set installCRDs=true
helm install rancher rancher-latest/rancher -n cattle-system -f rancher-master/rancher-helm-values.yaml --create-namespace
kubectl -n cattle-system get certificaterequest
kubectl -n cattle-system get issuer
kubectl -n cattle-system get certificate
kubectl -n cattle-system get order
kubectl -n cattle-system get secret
kubectl -n cattle-system get pods
```

At this point, Rancher should be running with cert-manager managing the Rancher-generated certificates. Get the login details and confirm if everything is working as expected.

```sh
kubectl get secret --namespace cattle-system bootstrap-secret -o go-template='{{.data.bootstrapPassword|base64decode}}'
echo https://rancher.rajesh-kumar.in/dashboard/?setup=$(kubectl get secret --namespace cattle-system bootstrap-secret -o go-template='{{.data.bootstrapPassword|base64decode}}')
```

## Step 5: Create Ingress Resource to Access Rancher

```sh
kubectl apply -f ingresses/rancher-ingress.yaml
```

I have hosted this on a personal domain; hence, I created a wildcard certificate for my domain. I also created a secret to store the certificates.

## Step 6: Store the Secret as YAML

```sh
kubectl -n cert-manager get secret rancher-tls-cert -o yaml | yq 'del(.metadata.creationTimestamp, .metadata.resourceVersion, .metadata.selfLink, .metadata.uid, .metadata.managedFields, .metadata.namespace, .metadata.annotations)' | kubectl apply -n cattle-system -f -
```

## Step 7: Store Secrets with Vault and External-Secrets

Look at the README.md under vaults/ section.

> **Note:** `rancher-master/rancher-tls-secret.yaml` needs to be created for each namespace and on each cluster where ingressroute is created. This file isn't part of the git repo for security issues.

## 🆕 Adding a New Remote Cluster

### Step 1: Deploy Kubernetes (K3s, RKE2, or RKE)

#### RKE:

[RKE Installation Guide](https://rke.docs.rancher.com/installation)

Use `rancher-master/cluster.yml` to deploy an RKE cluster. Make necessary changes, or run `rke config` to create a `cluster.yaml`.

#### K3s:

[K3s Installation Guide](https://docs.k3s.io/installation/configuration#configuration-with-install-script)

```sh
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.26.7+k3s1 INSTALL_K3S_EXEC="server" sh -
```

### Step 2: Ensure Traefik is Deployed

For RKE:

```sh
helm upgrade --install traefik traefik/traefik -n kube-system -f applications/home-cluster-rke/traefik-config.yaml
```

For K3s:

SSH to the master node and create `/var/lib/rancher/k3s/server/manifests/traefik-config.yaml`

```yaml
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

### Step 3: Deploy External-Secret

```sh
helm install external-secrets external-secrets/external-secrets -n external-secrets --create-namespace
```

### Step 4: Create Vault Secret and ClusterSecretStore

```sh
kubectl apply -f rancher-master/vaults/vaultTokenSecret.yaml
kubectl apply -f rancher-master/clusterSecretStore.yaml
```

### Step 5: Create Rancher TLS Secret in App Namespaces

```sh
kubectl -n app-ns apply -f rancher-master/rancher-tls-secret.yaml
```

### Step 6: Expose Application with IngressRoute and Use Rancher TLS Secret for TLS

### Step 7: Ensure Everything is Working

```sh
kubectl get clustersecretstore
kubectl get externalsecret -A
kubectl get ingressroute -A
```

### Create ImagePullSecrets for DockerHub and Use It with the Applications

```sh
kubectl -n cattle-neuvector-system create secret generic regcred --from-file=.dockerconfigjson=/home/rajesh/.docker/config.json --type=kubernetes.io/dockerconfigjson
```
