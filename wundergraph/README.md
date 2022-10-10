# 安装依赖
```shell
yarn
```


# 集成prisma 以及prismaClient
```shell
npm install prisma typescript ts-node @types/node --save-dev
npx prisma

npm install @prisma/client
npx prisma generate
```

# 启动后台
进入.wundergraph输入命令
```shell
wunderctl up --debug
```

# 启动web
在根目录下运行命令
```shell
npm run nextDev
```