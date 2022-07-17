# Deployment (Server)

## Steps
1. Did all step in `deployment/remote/README.md`
2. Run ssh to your server. Example:
```
ssh user@127.0.0.1:/home/user -p 80
```

## In your server
1. Prepare your config
Adjust `config/config.yaml`

1. Run autonotif
```
make run
```

3. Stop autonotif
```
make stop
```
