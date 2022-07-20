# Autonotif

## How to Run

### Local

#### Prerequisite
- [Docker](https://www.docker.com/get-started/)
- [Golang](https://go.dev/learn/) (if preferred)

#### 1. Prepare your config
Copy `config/config.yaml.sample` for example. Then adjust accordingly.
```
cp config/config.yaml.sample config/config.yaml
```

#### 2. Run autonotif

__Docker__
```
make run
```

__Golang__
```
make go-run
```

## Optional Usage

### Force Set Last ID
Make sure autonotif server already run.

__Force autonotif to start notify from specific id__
```
curl --request POST 'http://localhost:8080/force-last-id' \
--header 'chain: COSMOS' \
--header 'lastId: 73'
```
__Force autonotif to start notify from zero__
```
curl --request POST 'http://localhost:8080/force-last-id' \
--header 'chain: COSMOS' \
--header 'lastId: -1'
```

## Deployment
See guideline in `deployment/remote/README.md`

## Features

| Features                                      | Cosmoshub | Osmosis | Juno    |
| --------------------------------------------- | --------- | ------- | ------- | 
| Blockchain Governance Proposals Notification  | Supported | Backlog | Backlog |
| Sending to Telegram bot                       | Supported | Backlog | Backlog |
| Sending to Email                              | Backlog   | Backlog | Backlog |
| Sending to Slack bot                          | Backlog   | Backlog | Backlog |
| Sending to WhatsApp bot                       | Backlog   | Backlog | Backlog |
| Sending to Discord bot                        | Backlog   | Backlog | Backlog |
| Blockchain Cosmos Upgrade Plan Notification   | Backlog   | Backlog | Backlog |
| Check "Have <name> voted to the Proposal <ID> | Backlog   | Backlog | Backlog |
| Check "Have <name> Upgraded to Plan <VERSION> | Backlog   | Backlog | Backlog |

## Room for Refactors

| Refactors                                     | Status    |
| --------------------------------------------- | --------- |
| Time zone for Docker autonotif-scheduler      | Done      |
| Time zone for Docker autonotif-postgres       | Done      |
| PostgreSQL Docker Compose                     | Done      |
| PostgreSQL Initialize Table - Proposals       | Done      |


