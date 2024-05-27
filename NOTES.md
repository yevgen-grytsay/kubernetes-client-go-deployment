

```sh
read -s PAT
export PAT
envsubst < templates/repo-https-secret.tpl.yaml > templates/repo-https-secret.yaml

```

## Prerequisites
- Додати репозиторій в ArgoCD (вже має бути)
- Сервіс-акаунт з усіма необхідними доступами

## TODO
- Додати кластер в ArgoCD програмно


## Resources
- [client-go | dynamic-create-update-delete-deployment](https://github.com/kubernetes/client-go/blob/v0.30.1/examples/dynamic-create-update-delete-deployment/main.go)
- [client-go | out-of-cluster-client-configuration](https://github.com/kubernetes/client-go/blob/v0.30.1/examples/out-of-cluster-client-configuration/main.go)
