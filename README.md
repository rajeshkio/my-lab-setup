🌟 My Kubernetes Lab Setup 🌟

Welcome to my Kubernetes lab setup! This document provides a comprehensive overview of the architecture and technologies I've implemented.

🖥️ Clusters Overview

🌐 Rancher Local Cluster
- 🖥️ Node Configuration: 1 node (16GB RAM, 8 CPU)
- 🐳 Kubernetes Distribution: K3s
- 🛠️ Operating System: Ubuntu 22.04
- ☁️ Hosting: Cloud-hosted
  
🌐 Rancher Imported Tenant (Cloud)
- 🖥️ Node Configuration: 2 nodes (8GB RAM, 4 CPU each)
- 🐳 Kubernetes Distribution: K3s
- 🛠️ Operating System: Ubuntu 22.04
- ☁️ Hosting: Cloud-hosted

🏠 Rancher Imported Tenant (Home)
- 🖥️ Node Configuration:
  - 1 node (8GB RAM, 4 CPU) - Laptop
  - 1 node (8GB RAM, 4 CPU) - Raspberry Pi 5
- 🐳 Kubernetes Distribution: K3s
- 🛠️ Operating System: Ubuntu 23.10

🗺️ Architecture Diagram

🔍 Diagram Summary

The architecture diagram illustrates the following components and their interactions:

👤 User Requests: Initiated by the user and forwarded by the Nginx Reverse Proxy.

🔀 Nginx Reverse Proxy: Directs user requests to one of the three clusters (Cluster 1, Cluster 2, Cluster 3).

🌐 Clusters:
- Cluster 1:
  - Deployed services include Vault.
  - Node Exporter collects and forwards metrics to Prometheus 1.
- Cluster 2:
  - Deployed services include NeuVector and Grafana.
  - Node Exporter collects and forwards metrics to Prometheus 2.
- Cluster 3:
  - Deployed services include ArgoCD.
  - Node Exporter collects and forwards metrics to Prometheus 3.
    
- 📊 Prometheus Instances: Each cluster has its own Prometheus instance (Prometheus 1, Prometheus 2, Prometheus 3) that scrapes and collects metrics.
- 🔗 Thanos: Aggregates metrics from all Prometheus instances and serves as a datasource for Grafana.
- 📈 Grafana: Visualizes metrics collected by Prometheus and Thanos.

🛠️ Tech Stack
- 🚀 Continuous Delivery: ArgoCD
- 🌍 Ingress Controller: Traefik
- 🔀 Reverse Proxy / Load Balancer: Nginx
- 📊 Monitoring: Prometheus, Grafana, Thanos
- 📜 Logging: Loki (In progress)
- 🔒 Security: NeuVector

