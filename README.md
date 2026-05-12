# 🎨 图标空间站 • IconStation

> 打造属于自己的图标库

> 私有化图标库，集成多个图标源，支持一键复制使用；用户可自上传图标进行收集管理。

<img width="1280" height="640" src="https://github.com/user-attachments/assets/0bf7ba21-7e8a-4367-99a6-bd237287ee22" />

***

## 💬 讨论交流

欢迎加入企鹅群讨论交流：**1039270739**

***

## ✨ 功能特点

| 特性        | 说明                                                  |
| --------- | --------------------------------------------------- |
| 🗂️ 多源集成  | 集成多个图标源，无需四处寻找                                      |
| 📋 一键复制   | 支持直接使用，点击即可复制                                       |
| 📤 自上传    | 用户可自上传图标，个性化收集管理                                    |
| 💾 数据持久化  | 用户数据本地安全存储                                          |
| 🖥️ 全平台支持 | Windows amd64, Linux amd64/arm64, macOS amd64/arm64 |

***

<img width="1543" height="1028" alt="PixPin_2026-05-07_11-24-43" src="https://github.com/user-attachments/assets/f6833fae-c239-48b1-bb05-27263f73d216" />  

***

## 🖱️ 操作说明

### 图标预览

| 操作方式                | 说明                    |
| :------------------ | :-------------------- |
| 🖱️ **点击缩略图**       | 点击左右两侧的缩略图可切换上一张/下一张  |
| ⬅️➡️ **方向键**        | 前后翻页翻页  或者  预览上一张/下一张 |
| 📱 **滑动手势**         | 移动端：左滑下一张，右滑上一张       |
| ❌ **ESC键 / 点击任意位置** | 关闭预览模式                |

### 工具栏功能

| 按钮          | 功能说明                             |
| :---------- | :------------------------------- |
| 📤 **上传**   | 点击上传图标文件，支持多文件批量上传               |
| 📚 **图标库源** | 展开选择要显示的图标库                      |
| 🔍 **搜索框**  | 输入关键词实时过滤图标                      |
| 🔄 **模式**   | 选择加载模式（Github加速、CDN加速、Github原地址） |
| 🎨 **主题**   | 选择界面主题（默认、浅色、深色）                 |

***

## 📝 更新日志

| 版本         | 更新内容                     |
| :--------- | :----------------------- |
| **v0.2.0** | 支持自建分类,本地图标管理更方便；增大预览图面积 |
| **v0.1.7** | 完美适配移动端；支持设置本地图标是否公开浏览   |
| **v0.1.6** | 新增图标预览，支持方向键切换           |
| **更早之前**   | 横空出世，达成基础功能              |

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

<img src="https://github.com/user-attachments/assets/2b62ae3b-566a-4147-92ac-5eed098afd3a" />

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
