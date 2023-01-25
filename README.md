# KO Vault Service Example

## Vault

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 0.0.0.0:8200
export VAULT_ADDR=http://0.0.0.0:8200
vault login root
colima start -c 6 -m 16 -k
export EXTERNAL_VAULT_ADDR=192.168.5.2
```

## KO

```bash
colima start -c 6 -m 16 -k
export DOCKER_HOST="unix://${HOME}/.colima/default/docker.sock"
ko apply -f .
k logs -l app=hello -f
k run httpie --image=alpine/httpie --rm -it --restart=Never -- hello-service:8080/Barney
```

## References

- <https://developer.hashicorp.com/vault/tutorials/kubernetes/kubernetes-external-vault>
- <https://github.com/hashicorp/hello-vault-go>
