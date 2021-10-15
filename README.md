# Gasnow 2
Clone of gasnow.org

## Quickstart

The application connects to infura to get data from its ethereum full node. In order to do so the following env variables need to be setup in a `.env` file:
```.env
INFURA_URL=https://mainnet.infura.io/v3
INFURA_PROJECT=<your_project_id>
```

To start the app during development
```bash
# run the app
make run

# test health endpoint
curl http://localhost:5000/health

# connect via websocket
wscat -c ws://localhost:5000/ws
```

### Build executable:
```bash
# build
make build

# Run the newly built executable
bin/server
```
