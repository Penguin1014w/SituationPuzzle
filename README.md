# 海龟汤游戏

一个使用Go和React开发的海龟汤游戏，集成了ChatGPT API来验证答案。

## 功能特点

- 预设多个海龟汤谜面
- 用户可以选择谜面进行猜测
- 使用ChatGPT API验证答案
- 限制猜测次数为5次
- 简洁美观的用户界面

## 技术栈

- 后端：Go
- 前端：React
- API：OpenAI ChatGPT

## 安装和运行

### 后端

1. 进入backend目录：
```bash
cd backend
```

2. 设置环境变量：
```bash
export OPENAI_API_KEY=你的OpenAI API密钥
```

3. 运行后端服务器：
```bash
go run main.go chatgpt.go
```

### 前端

1. 进入frontend目录：
```bash
cd frontend
```

2. 安装依赖：
```bash
npm install
```

3. 启动开发服务器：
```bash
npm start
```

## 使用说明

1. 在首页选择想要猜测的谜面
2. 输入你的答案
3. 系统会使用ChatGPT验证你的答案是否正确
4. 你有5次猜测机会
5. 猜对或用完猜测次数后可以返回选择其他谜面

## 注意事项

- 确保已设置正确的OpenAI API密钥
- 后端服务器默认运行在8080端口
- 前端开发服务器默认运行在3000端口 