# Play with tower

此案例主要通过收集 Prometheus 实例指标用于展示遥测塔数据库服务的基本使用。

## 如何运行

1. 通过[账号申请](https://docs.telemetrytower.com/yao-ce-ta/application) 获得授权数据库的 Authorization Token。
2. 分别替换 config/agent/文件中 `xxx` 内容为申请的账号 Token。 
3. 执行 `docker-compose up -d` 启动服务。 
4. 通过[数据查询](https://docs.telemetrytower.com/yao-ce-ta/query)进行数据读取等操作。