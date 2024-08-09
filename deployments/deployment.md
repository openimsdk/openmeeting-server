# OpenMeeting Server Source Code Deployment Guide

### 1. Downloading the Source Code

```bash
git clone https://github.com/openimsdk/openmeeting-server.git && cd openmeeting-server
```


### 2. Deploying Related Dependencies Component (Etcd, MongoDB, Redis, LiveKit)
```bash
# install dependencies component
docker compose up -d

# Checking if Related Dependencies components are Running Properly
docker ps
```

### 3. Set external IP
```bash
Modify the `url` in `config/live.yml` to `ws://external_IP:17880` or a domain name.
```

### 4. Initialization
Before the first compilation, execute on Linux/Mac platforms:
```bash
bash bootstrap.sh
```
On Windows execute:
```bash
bootstrap.bat
```

### 5. compile and run
```bash
mage && mage start
```






