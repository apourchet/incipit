# incipit
A boilerplate for a scalable Kubernetes/Docker backed infrastructure. The directory structure 
is meant to be simple and easily extensible.

# Features
- Quick bootstrapping of local k8s cluster
- Ingress controller with Path-Based-Routing and TLS support
- Redis single-node
- Etcd clustered
- Ingress with TLS and path-based-routing
- Environment-independent build targets
- Example of RPC server working with etcd/redis

# Planned
- Automatic Letsencrypt cert renewal
- Continuous Integration with one of CircleCI/Jenkins/TravisCI
- Automatic transparent intra-VPC encryption

# Goal
Provide a clean, robust example of a scalable backend that is also language-agnostic
and easy to develop on; while still benefitting from the microservice architecture.
