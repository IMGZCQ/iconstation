# 图标空间站 • iconstation
集成多个图标源直接使用
用户可自上传收集管理
<img width="1528" height="1318" alt="image" src="https://github.com/user-attachments/assets/11ea1704-db7f-4f82-9187-63c61f70ec16" />

使用方法
services:
  iconstation:
    container_name: iconstation      # 容器名称
    image: imgzcq/iconstation        # Docker镜像名称
    ports: [9168:9168]             # 访问端口:内部端口
    volumes: [./UserData:/app/UserData]  # 用户图标持久化目录
    restart: always               # 开机自启
