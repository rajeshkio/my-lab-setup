### Argocd Workflow

![Argocd Workflow](https://github.com/rajeshkio/my-lab-setup/blob/fbffd5fd2d0f80f8b826e1739092a78e7a27dc3d/argocd-helm-setup/argocd-github-implementation.png)

### Deploy ArgoCD

$ kubectl create ns argo

$ kubectl apply -f argocd-helm-setup/externalSecretDexConfig.yaml -f argocd-helm-setup/externalSecretSlackConfig.yaml -f argocd-helm-setup/externalSecretWebhook.yaml -f argocd-helm-setup/argo-ingress.yaml

$ helm upgrade --install argo -n argo argo/argo-cd --values argocd-helm-setup/values.yaml

$ kubectl -n argo get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d | xargs

https://github.com/argoproj/argo-cd/blob/master/docs/faq.md

### Enable Auth0 with github

1 - Create github oAuth app and get clientId and secret. ( Hint: settings --> Developer options --> oAuth). I have created an organisation and created oAuth from the organisation settings.

2 - Create a secret named `argocd-secret` in `argo` namespace.

3 - Include the client and secret id variable from the secret to argocd-cm config-map. (https://argo-cd.readthedocs.io/en/stable/operator-manual/user-management/)

4 - Update the argocd-rbac-cm with new role and policies for this github user. (https://argo-cd.readthedocs.io/en/stable/operator-manual/rbac/)


For multi Rancher cluster see this : https://gist.github.com/janeczku/b16154194f7f03f772645303af8e9f80

kubectl apply -f rancher-master/argocd-multi-rancher-cluster-secret.yaml

metrics exporter issue on ARM64 https://github.com/argoproj/argo-helm/issues/2233

Add anotations to appproject crd
kh -n argo edit appproject // This will enable notifications for all the apps deployed in the project

annotations:
	notifications.argoproj.io/subscribe.on-deleted.slack: argocd-alerts
	notifications.argoproj.io/subscribe.on-deployed.slack: argocd-alerts
        notifications.argoproj.io/subscribe.on-health-degraded.slack: argocd-alerts
	notifications.argoproj.io/subscribe.on-sync-failed.slack: argocd-alerts
	notifications.argoproj.io/subscribe.on-sync-status-unknown.slack: argocd-alerts
	notifications.argoproj.io/subscribe.on-sync-succeeded.slack: argocd-alerts

