# Roadmap

## Key Goals

1. **E2E Testing for Main Connectors**
2. **Helm Deployment Documentation**
3. **K8s Operator for Multi-Transfer Deployments**

---

## 1. E2E Testing for Main Connectors

### Objective:
Set up comprehensive **end-to-end tests** in the CI pipeline for the following main connectors:
- **Postgres**
- **MySQL**
- **Clickhouse**
- **Yandex Database (YDB)**
- **YTsaurus (YT)**

### Steps:
- [ ] Configure test environments in CI for each connector.
- [ ] Design E2E test scenarios covering various transfer modes (snapshot, replication, etc.).
- [ ] Automate test execution for all supported connectors.
- [ ] Set up reporting and logs for test failures.

### Milestone:
Achieve **fully automated E2E testing** across all major connectors to ensure continuous integration stability.

---

## 2. Helm Deployment Documentation

### Objective:
Provide detailed documentation on deploying the transfer engine using **Helm** on Kubernetes clusters.

### Steps:
- [ ] Create Helm chart for easy deployment of the transfer engine.
- [ ] Write comprehensive **Helm deployment guide**.
    - [ ] Define key parameters for customization (replicas, resources, etc.).
    - [ ] Instructions for various environments (local, cloud).
- [ ] Test Helm deployment process on common platforms (GKE, EKS, etc.).

### Milestone:
Enable seamless deployment of the transfer engine via Helm with clear and accessible documentation.

---

## 3. Kubernetes Operator for Multi-Transfer Deployments

### Objective:
Develop a **Kubernetes operator** to manage multiple data transfers, simplifying the process for large-scale environments.

### Steps:
- [ ] Define CRD (Custom Resource Definitions) for transfer configurations.
- [ ] Implement operator logic for scaling and managing multi-transfer deployments.
- [ ] Add support for monitoring, scaling, and error recovery.
- [ ] Write user documentation for deploying and managing transfers via the operator.

### Milestone:
Provide a scalable solution for managing multiple data transfers in Kubernetes environments with an operator.

---

## Summary

- **Q2-Q3**: Focus on **E2E testing** for core connectors.
- **Q3**: Publish **Helm deployment** documentation and final testing.
- **Q3-Q4**: Develop and release the **Kubernetes operator** for multi-transfer management.

This roadmap aims to enhance testing, simplify deployment, and provide advanced scalability options for the transfer engine.
