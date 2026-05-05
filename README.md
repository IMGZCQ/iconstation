# 图标空间站 • IconStation

集成多个图标源，直接使用；支持用户自上传收集管理。

<img width="1214" height="1271" alt="image" src="https://github.com/user-attachments/assets/f0f9e6f1-a011-42ed-b8de-b97e5af33bed" />


## 功能特点

- 集成多个图标源，无需四处寻找
- 支持直接使用，一键复制
- 用户可自上传图标，个性化收集管理
- 数据持久化存储

## 快速开始

### 环境要求

- Docker
- Docker Compose

### 安装部署

```yaml
services:
  iconstation:
    container_name: iconstation         # 容器名称
    image: imgzcq/iconstation           # Docker镜像名称
    ports: [9168:9168]                  # 访问端口:内部端口
    volumes: [./UserData:/app/UserData] # 用户图标持久化目录
    restart: always                     # 开机自启
```

### 启动命令

```bash
docker-compose up -d
```

### 访问服务

启动成功后，访问 <http://localhost:9168>

## 目录说明

| 目录           | 说明          |
| ------------ | ----------- |
| `./UserData` | 用户图标持久化存储目录 |

## 端口映射

| 主机端口 | 容器端口 | 说明       |
| ---- | ---- | -------- |
| 9168 | 9168 | Web 访问端口 |

## 维护说明

- 重置密码删除`UserData` 目录下 `pw.json` 即可
- 用户上传的图标存储在 `UserData/icons` 目录中

