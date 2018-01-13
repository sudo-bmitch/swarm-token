# Swarm Token
Retrieves the docker swarm mode token needed to join a remote node.
Uses a shared secret to retrieve the token and provides it over http.
This allows swarm workers to join without needing to know the token in
advance, which is useful for deployments with Terraform, Cloud Formation,
and any elastic scaling process.

## Quick start:

```
WORKER_KEY=123 MANAGER_KEY=abc docker-compose up

curl -sSL -H "X-Key: 123" http://localhost:8888/worker | jq -r .Token
curl -sSL -H "X-Key: abc" http://localhost:8888/manager | jq -r .Token
```

## Disclaimer
Do not expose this to the public internet, this should only be started on a 
private network that you trust. The Key and Token are transmitted over clear
text.

