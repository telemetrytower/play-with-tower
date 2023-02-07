# MyService 服务 SLO 监控

本样例主要演示 MyService 指标收集，基于 Sloth 的 SLO 告警监控。

## 程序运行

-  下载代码到本地

```
git clone https://github.com/telemetrytower/play-with-tower.git
```

修改配置中的 token 为 telemetrytower 数据库真实值，主要修改 [grafana-provisioning-datasources.yaml#L27](https://github.com/telemetrytower/play-with-tower/blob/master/myservice/config/grafana/grafana-provisioning-datasources.yaml#L27) 、[grafana-provisioning-datasources.yaml#L15](https://github.com/telemetrytower/play-with-tower/blob/master/myservice/config/grafana/grafana-provisioning-datasources.yaml#L15)、[prometheus.yaml#L12](https://github.com/telemetrytower/play-with-tower/blob/master/myservice/config/prometheus.yaml#L12) 这三处。


- 运行程序

```
cd play-with-tower/myservice
docker-compose up -d
```

## 加载 rule 和 alertmanager 配置

这里我们主要使用 tower-tool 命令行工具进行 rule 和 alertmananger 配置管理，在使用之前通过环境变量注入 JTW TOKEN。

```
export TOWER_AUTH_TOKEN=token
```

- 加载 rules

我们主要加载 config/rules/myservice.yaml 文件的 rules，此内容是根据 slos/myservice.yaml 配置的 SLO 通过 sloth 自动生成的。

```
tower-tool rules load ./config/rules/myservice.yaml
```

加载完成后，可以使用`tower-tool rules list` 查看刚加载的所有告警分组。 

- 加载 alertmanager

使用以下命令进行 alertmanager 配置加载：

```
tower-tool alermanager load ./config/alertmanager/alertmanager.yaml ./config/alertmanager/default_template.tmpl
```

加载完成后，可以使用`tower-tool alermanager get` 查看刚加载的所有 alertmanager 配置。 

## Grafana 查看

最后我们可以访问 [http://localhost:3000](http://localhost:3000) 看到刚配置的 rules 和 Alertmananger。






