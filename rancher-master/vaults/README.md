
```sh
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update
helm install vault hashicorp/vault
kubectl exec -it vault-0 -- vault operator init // COPY the tokens and store.
kubectl exec -it vault-0 -- sh

$ vault operator unseal // PASTE unseal tokens
$ vault operator unseal // PASTE unseal tokens
$ vault operator unseal // PASTE unseal tokens
$ vault login $INITIAL_ROOT_TOKEN
$ vault secrets enable --version=2 --path=kv kv
$ vault kv put -mount=kv neuvector adminUser=test
$ vault kv put -mount=kv neuvector adminPass='password'
$ // create all the secret here or can create from the UI as well.
$ exit

kubectl apply -f vaults/vaultTokenSecret.yaml
kubectl apply -f ingresses/vault-ingress.yaml
kubectl apply -f clusterSecretStore.yaml
```
