# 🎨 图标空间站 • IconStation

> 打造属于自己的图标库

> 私有化图标库，集成多个图标源，支持一键复制使用；用户可自上传图标进行收集管理。

![IconStation Preview](https://github.com/user-attachments/assets/f0f9e6f1-a011-42ed-b8de-b97e5af33bed)

***

## ✨ 功能特点

| 特性       | 说明               |
| -------- | ---------------- |
| 🗂️ 多源集成 | 集成多个图标源，无需四处寻找   |
| 📋 一键复制  | 支持直接使用，点击即可复制    |
| 📤 自上传   | 用户可自上传图标，个性化收集管理 |
| 💾 数据持久化 | 用户数据本地安全存储       |

***

## 🚀 快速开始

### 方式一：飞牛应用 FPK 下载安装

📦 下载 FPK 安装包：  

[![点击下载 FPK](https://img.shields.io/badge/Download-FPK-ff69b4?style=for-the-badge)](https://fndesk.imcq.top/?url=dl\&at=GitHUb)

***

### 方式二：飞牛三方商店 Fndepot

🏪 直接在 Fndepot 商店 搜索 **"图标空间站"** 下载安装

***

### 方式三：Docker Compose 部署

#### 1️⃣ 创建 docker-compose.yml 文件

```yaml
services:
  iconstation:
    container_name: iconstation         # 容器名称
    image: imgzcq/iconstation           # Docker镜像名称
    ports: [9168:9168]                  # 访问端口:内部端口
    volumes: [./UserData:/app/UserData] # 用户图标持久化目录
    restart: always                     # 开机自启
```

#### 2️⃣ 启动服务

```bash
docker-compose up -d
```

#### 3️⃣ 访问服务

🌐 启动成功后，访问 **<http://localhost:9168>**

***

## 📁 目录结构

| 路径                   | 说明          |
| -------------------- | ----------- |
| `./UserData`         | 用户图标持久化存储目录 |
| `./UserData/icons`   | 用户上传的图标目录   |
| `./UserData/pw.json` | 登录密码配置文件    |

***

## 🔧 端口说明

|  主机端口  |  容器端口  | 说明         |
| :----: | :----: | ---------- |
| `9168` | `9168` | Web 服务访问端口 |

***

## 🔒 安全维护

| 操作          | 方法                            |
| ----------- | ----------------------------- |
| 🔑 重置密码     | 删除 `UserData/pw.json` 文件后重启服务 |
| 🗑️ 备份或移植图标 | 备份 `UserData/icons` 目录下所有文件   |

***

## ☕ 支持项目

如果你觉得这个项目对你有帮助，欢迎打赏支持！

![打赏](https://github.com/user-attachments/assets/7c58ea66-ca1a-4d65-8b26-cce89bde813c)

***

![3](https://github.com/user-attachments/assets/6b17ae00-a25a-4095-bad1-d3f734c27c2d)

![2](https://github.com/user-attachments/assets/88eedbb8-7f92-4aa9-b49a-d2a90f8f5464)

## ![1](https://github.com/user-attachments/assets/fb015080-dc63-4a85-80fd-5f6ad25f5630)

## 📝 开源协议

本项目基于 MIT 协议开源，欢迎 Star ⭐ 和贡献！
