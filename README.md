
# My Lab Setup

The project is an overview of how built my lab server for kubernetes.

Three clusters:

    1 - Rancher local cluster :
      - 1 node (16GB, 8 CPU)
      - K3s
      - Ubuntu 22.04
      - Cloud hosted 
    2 - Rancher imported tenant
      - 2 nodes (8GB 4 CPU each)
      - K3s
      - Ubuntu 22.04
      - Cloud hosted
    3 - Rancher imported tenant setup at the home
      - 1 node (8GB 4CPU) Laptop
      - 1 node (8GB 4 CPU) Raspberry 5
      - Ubuntu 23.10
      - K3s

![My-Lab-Setup-Arch](https://github.com/rk280392/my-lab-setup/blob/e5dba8b764f3218ef555532deeb8ff51646e916f/my-lab-setup.png)

## Tech Stack

Continuous Delivery: ArgoCD

Ingress controller: Traefik

Reverse Proxy / Load Balancer: Nginx

Monitoring: Prometheus, Grafana, Thanos

Logging: Loki (In progress)

Security: Neuvector 
