# SamoyedQQBot

## 部署

### 环境变量

* `BOT_TOKEN` 必须设置，对应[开发文档](https://bot.q.qq.com/wiki/develop/api/#bot-token)中的配置项
* `ONLINE` 选填，`true` 时使用正式环境 API，否则默认沙箱环境

### 运行

1. 直接运行：

```bash
> BOT_TOKEN="Your token" make run
```

2. 用 Docker 运行

```bash
> Docker built . -t qqbot
> docker run --env BOT_TOKEN="Your token" qqbot 
```

## 功能

1. 自我介绍

<img src="./assets/IMG_2437.jpg" alt="IMG_2437" style="zoom: 33%;" />

2. 单词接龙

<img src="./assets/IMG_2438.jpg" alt="IMG_2438" style="zoom:25%;" />

3. 错误提示：错误单词 & 错误开头

<img src="./assets/IMG_2439.jpg" alt="IMG_2439" style="zoom:25%;" />

4. 重新开始游戏

<img src="./assets/IMG_2440.jpg" alt="IMG_2440" style="zoom:25%;" />

## 设计

