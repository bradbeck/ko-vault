# KO Service Example

## KO

```bash
colima start -c 6 -m 16 -k
export DOCKER_HOST="unix://${HOME}/.colima/default/docker.sock"
ko apply -f .
k logs -l app=hello -f
k run httpie --image=alpine/httpie --rm -it --restart=Never -- hello-service:8080/Barney
```
