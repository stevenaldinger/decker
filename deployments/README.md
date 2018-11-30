# `/deployments`

IaaS, PaaS, system and container orchestration deployment configurations and templates (docker-compose, kubernetes/helm, mesos, terraform, bosh).

# Docker Compose

[docker-compose.yml](docker-compose.yml) runs `decker` with two volumes, one for config files and one for generated reports.

# Kubernetes

## Job

[k8s-job.yaml](k8s-job.yaml) runs `decker` once as a [Job](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/) in a Kubernetes cluster.

## CronJob

[k8s-cronjob.yaml](k8s-cronjob.yaml) runs `decker` once every hour as a [CronJob](https://kubernetes.io/docs/tasks/job/automated-tasks-with-cron-jobs/) in a Kubernetes cluster.
