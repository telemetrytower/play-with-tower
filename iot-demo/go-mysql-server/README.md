## go-mysql-server-proxy
此项目启发自https://github.com/dolthub/go-mysql-server。  
本项目旨在提供一个适配sql协议的接入层协议，方便原有的基于sql语义的写入端，做较少的改动，将sql数据写入其他存储系统。  
本项目暂不提供sql的持久化存储功能。目前的底层存储，基于Prometheus生态构建，如果您需要其他存储，如sqllite等，可以参考本工程进行扩展即可。  

## 编译方法
对于linux，直接执行./build.sh
对于windows,执行 go build main.go

## 使用方法
【编译好的可执行程序，在release目录下，包含linux/windows版本。】  
【启动位置不限，在可执行程序目录下，创建好您的数据库同名的目录即可。】
以运动数据为例，进行说明。假设我们手环的运动数据，有如下格式：
```
stepinfo 
(
usrid  varchar(255) ,  
sportmode int ,
starttime int ,
steps int,
calo int,
wendu float,
)
```

上述数据字段分别表示：用户id,运动模式，运动数据上报时间点（1min粒度），运动步数，消耗卡路里，当时温度。
### 对于linux
windows部分也试用该部分说明，只是操作有稍微差异，可根据实际进行调整。
#### 启动
创建“sqldata”文件目录。
```
./main "sqldata/sport"
```
因为sql客户端连接时，需要指定数据库名称，上述“sqldata/sport”即表示该工具创建的数据库名称。如果您有多个数据库，启动多个工具即可。
默认用户名，root，密码为空。
#### 创建表
注意，只需要创建1次即可，创建表是为了让该sql工具能解析原生的sql语义。创建一次后，无需再分库分表。
如果客户端根据存储系统协议，不再使用sql语义时，可参考我们其他方案，将数据直接写入后端存储，不再需要此工具。
```
mysql --host=127.0.0.1 --port=3306 -u root sqldata/sport -e "CREATE TABLE stepinfo
(
usrid  varchar(255) ,
sportmode int ,
starttime int ,
steps int,
calo int,
wendu float,
);"
```
#### 指定关键字段
对于上述运动数据，如果用户id，运动模式确定，那么在某个时间点，他的步数、卡路里、温度，就是确定的。所以，我们需要区分用户id,运动模式，这两个"关键字段"。步数、卡路里、温度我们称为"业务字段"
```
curl "ip:3400/v1/tableflag?table=stepinfo&keyflags=usrid,sportmode&keytime=starttime"
```
使用上述curl命令，给我们的工具发送一个信息，告诉他：usrid,sportmode作为keyflags；starttime作为keytime。其他业务字段比如steps,我们的工具会自动识别，无需再做处理。所以当您的设备有多个"业务字段"的时候，只需要提供"关键字段"

#### 写入数据
```
mysql --host=127.0.0.1 --port=3306 -u root sqldata/sport -e "INSERT INTO stepinfo (usrid, sportmode, starttime, steps,calo,wendu) VALUES ('3',1,1681959664,100,300,28.7), ('5',1,1681959664,100,300,28.7), ('6',2,1681959664,100,300,28.7);"
```
如上就是普通的sql数据写入了，可以单条写入，也可以批量写入。

### 对于windows
#### 启动
创建“sqldata”文件目录。
cmd终端启动
```
main.exe  "sqldata/sport"
```
本项目可直接在window64机器运行，运行时，不需要go环境。

注意如果需要在windows编译，要安装go环境。可自行安装，也可以参考博客，https://www.cnblogs.com/GreenForestQuan/p/14411115.html 
```
go build main.go
```
#### 获取sql客户端工具
如果您已经有sql链接工具，可以直接使用root登录即可。如果您需要单独sql命令行，可以下载https://dev.mysql.com/downloads/mysql/ 解压后，将xxx/mysql-8.0.33-winx64\bin添加到环境系统环境变量即可。
#### 创建表
说明见上文linux部分。
安装好mysql客户端工具后，执行cmd命令行执行
```
mysql --host=127.0.0.1 --port=3306 -u root sqldata/sport -e "CREATE TABLE stepinfo14 ( usrid varchar(255),sportmode int,starttime int,steps int,calo int,wendu float);"
```
#### 指定关键字段
说明见上文linux部分。
本地浏览器访问：
```
 "http://ip:3400/v1/tableflag?table=stepinfo&keyflags=usrid,sportmode&keytime=starttime"
```
使用您自己本机ip,替换上述地址即可。

## 数据查询

目前数据导入底层存储后，可直接在如下平台查看数据，  
地址：https://grafana.telemetrytower.com/d/yDRMq0EVk/testboard?orgId=3  
账号：test  
密码：test  



