# KO Vault Agent Service Example

This is an example of using the Vault agent for `auto_auth` and to inject secrets to a configuration file `config.txt`.
There are two applications `hello` and `reader`, each uses a different approach.

The `hello` app assumes the Vault agent injector has been deployed and configured. It relies on annotations to trigger
the creation of init and sidecar Vault agent container.

The `reader` app explicitly added init and sidecar containers to the deployment manifest and does not require
the Vault agent injector to be deployed.

Both apps use `ConfigMap`'s to configure the Vault agents to write the Vault token and the templated `config.txt` file.
Each app has its own corresponding secrets path that it is allowed to `read` and `list`.

The following alias is used throughout the examples.

```bash
alias k=kubectl
```

## Vault

Run a development Vault server with a specified dev root token.

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 0.0.0.0:8200
```

In a separate shell go about starting a kubernetes cluster with [colima](https://github.com/abiosoft/colima) and configure
kubernetes authorization in Vault for that cluster.

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
VAULT_AUTH_JWT=$(kubectl get secret vault-auth-token --output='go-template={{ .data.token }}' | base64 --decode)
KUBE_CA_CERT=$(kubectl config view --raw --minify --flatten --output='jsonpath={.clusters[].cluster.certificate-authority-data}' | base64 --decode)
KUBE_HOST=$(kubectl config view --raw --minify --flatten --output='jsonpath={.clusters[].cluster.server}')
vault write auth/kubernetes/config \
     token_reviewer_jwt="$VAULT_AUTH_JWT" \
     kubernetes_host="$KUBE_HOST" \
     kubernetes_ca_cert="$KUBE_CA_CERT"
```

Setup a policy and role for `hello` in Vault. Test that `$VAULT_AUTH_JWT` can be used request a token
for that role.

```bash
# Hello setup
vault policy write hello - <<EOF
path "secret/data/hello/*" {
  capabilities = ["read", "list"]
}
EOF
vault kv put secret/hello/config username='helloUser' password='helloSecret' ttl='30s'
vault write auth/kubernetes/role/hello \
     bound_service_account_names=vault-auth \
     bound_service_account_namespaces=default \
     policies=hello \
     ttl=24h
http POST :8200/v1/auth/kubernetes/login jwt=$VAULT_AUTH_JWT role=hello
```

Setup a policy and role for `reader` in Vault. Test that `$VAULT_AUTH_JWT` can be used request a token
for that role.

```bash
# Reader setup
vault policy write reader - <<EOF
path "secret/data/reader/*" {
  capabilities = ["read", "list"]
}
EOF
vault kv put secret/reader/config username='readerUser' password='readerSecret' ttl='30s'
vault write auth/kubernetes/role/reader \
     bound_service_account_names=vault-auth \
     bound_service_account_namespaces=default \
     policies=reader \
     ttl=24h
http POST :8200/v1/auth/kubernetes/login jwt=$VAULT_AUTH_JWT role=reader
```

## KO

[`ko`](https://github.com/ko-build/ko) is used to conveniently build and deploy each of the applications.

Configure `ko` to use local mode.

```bash
export KO_DOCKER_REPO=ko.local
```

Set `DOCKER_HOST` to point to `colima`.

```bash
export DOCKER_HOST="unix://${HOME}/.colima/default/docker.sock"
```

### Hello Service

Deploy the `hello` application and use `rollout` to wait for it to deploy completely.

```bash
ko apply -f hello
k rollout status deployment/hello-deployment
```

The logs can be inspected.

```bash
k logs -l app=hello -c hello -f --tail=-1
k logs -l app=hello -c vault-agent-init -f --tail=-1
k logs -l app=hello -c vault-agent -f --tail=-1
```

The `hello` service can be invoked.

```bash
k run httpie --image=alpine/httpie --rm -it --restart=Never -- hello-service:8080/hello/config
```

### Reader Service

Deploy the `reader` application and use `rollout` to wait for it to deploy completely.

```bash
ko apply -f reader
k rollout status deployment/reader-deployment
```

The logs can be inspected.

```bash
k logs -l app=reader -c reader -f --tail=-1
k logs -l app=reader -c vault-agent -f --tail=-1
k logs -l app=reader -c vault-agent-init -f --tail=-1
```

The `reader` service can be invoked.

```bash
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
