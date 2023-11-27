## Secure Thanos receive endpoint with proxy config

1. Use the htpasswd utility to create a file containing username and password pairs. This file will be used for basic authentication.

`htpasswd -cB /path/to/htpasswdfile username`

kubectl -n prometheus create secret generic thanos-receive-auth-secret --from-file=argocd/applications/thanos/htpassword

2. Configure Traefik Middleware for Basic Authentication and use the middleware in ingressroute config 
