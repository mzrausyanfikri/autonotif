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
