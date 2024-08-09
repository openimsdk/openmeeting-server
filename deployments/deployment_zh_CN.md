# OpenMeeting Server 源码部署指南

### 1. 下载源码

```bash
git clone https://github.com/openimsdk/openmeeting-server.git && cd openmeeting-server
```


### 2. 部署相关依赖组件(Etcd, MongoDB, Redis, LiveKit)
```bash
# 安装依赖组件
docker compose up -d

# 检查相关依赖组件是否正常运行
docker ps
```

### 3. 设置外部IP
```bash
Modify the `url` in `config/live.yml` to `ws://external_IP:17880` or a domain name.
```

### 4. 初始化
第一次编译前，linux/mac平台下执行：
```bash
bash bootstrap.sh
```

windows执行
```bash
bootstrap.bat
```

### 5. 编译以及运行
```bash
mage && mage start
```






