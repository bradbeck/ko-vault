# KO Vault Agent Service Example

```bash
alias k=kubectl
```

## Vault

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 0.0.0.0:8200
```

```bash
export VAULT_ADDR='http://0.0.0.0:8200'
vault login root
colima start -c 6 -m 16 -k
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update
helm install vault hashicorp/vault --set "injector.externalVaultAddr=http://external-vault:8200"
k apply -f external-vault.yaml
k apply -f service-acct.yaml
k run httpie --image=alpine/httpie --rm -it --restart=Never -- external-vault:8200
vault auth enable kubernetes
TOKEN_REVIEW_JWT=$(kubectl get secret internal-app-token --output='go-template={{ .data.token }}' | base64 --decode)
KUBE_CA_CERT=$(kubectl config view --raw --minify --flatten --output='jsonpath={.clusters[].cluster.certificate-authority-data}' | base64 --decode)
KUBE_HOST=$(kubectl config view --raw --minify --flatten --output='jsonpath={.clusters[].cluster.server}')
vault write auth/kubernetes/config \
     token_reviewer_jwt="$TOKEN_REVIEW_JWT" \
     kubernetes_host="$KUBE_HOST" \
     kubernetes_ca_cert="$KUBE_CA_CERT" \
     issuer="https://kubernetes.default.svc.cluster.local"
vault policy write hello - <<EOF
path "secret/data/hello/*" {
  capabilities = ["read", "list"]
}
EOF
vault kv put secret/hello/config username='appuser' password='suP3rsec(et!' ttl='30s'
vault write auth/kubernetes/role/hello \
     bound_service_account_names=internal-app \
     bound_service_account_namespaces=default \
     policies=hello \
     ttl=24h
http POST :8200/v1/auth/kubernetes/login jwt=$TOKEN_REVIEW_JWT role=hello
```

## KO

### Hello Service

```bash
export DOCKER_HOST="unix://${HOME}/.colima/default/docker.sock"
ko apply -f hello
k logs -l app=hello -c hello -f --tail=-1
k logs -l app=hello -c vault-agent-init -f --tail=-1
k logs -l app=hello -c vault-agent -f --tail=-1
k run httpie --image=alpine/httpie --rm -it --restart=Never -- hello-service:8080/hello/config
```

### Reader Service

```bash
export DOCKER_HOST="unix://${HOME}/.colima/default/docker.sock"
ko apply -f reader
k logs -l app=reader -c reader -f --tail=-1
k logs -l app=reader -c vault-agent -f --tail=-1
k logs -l app=reader -c vault-agent-init -f --tail=-1
k run httpie --image=alpine/httpie --rm -it --restart=Never -- reader-service:8080/vault/secrets/config.txt
```

## References

- <https://developer.hashicorp.com/vault/tutorials/kubernetes/kubernetes-external-vault>
- <https://github.com/hashicorp/hello-vault-go>
- <https://developer.hashicorp.com/vault/docs/platform/k8s/injector/annotations>
- <https://developer.hashicorp.com/vault/docs/auth/kubernetes>
- <https://developer.hashicorp.com/vault/docs/agent/autoauth>
- <https://developer.hashicorp.com/vault/tutorials/kubernetes/agent-kubernetes>
- <https://developer.hashicorp.com/vault/tutorials/kubernetes/agent-kubernetes#start-vault-agent-with-auto-auth>
- <https://developer.hashicorp.com/vault/docs/platform/k8s/injector/examples>
- <https://developer.hashicorp.com/vault/docs/platform/k8s/injector/examples#configmap-example>
- <https://github.com/MikeDafi/Nielsen-Internship---React-Flask-Hashicorp-Vault/tree/main/vaulthelm/templates>
- <https://github.com/openlab-red/hashicorp-vault-for-openshift/tree/master/examples/golang-example>
- <https://github.com/ConsenSys/quorum-key-manager-helm/blob/main/templates/config-agents.yaml>
- <https://medium.com/ww-engineering/working-with-vault-secrets-on-kubernetes-fde381137d88>
