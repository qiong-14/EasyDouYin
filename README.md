# EasyDouYin

## 抖音项目服务端简单示例

具体功能内容参考飞书说明文档

运行docker
```shell
docker-compose up
```
写入环境变量
```shell
source dy_secure_config.sh
```
运行程序
```shell
make all
```
## 注册用户
请求
```shell
curl --request POST 'http://localhost:8080/douyin/user/register/?username=testname&password=12345'
```
注册成功
```shell
{"status_code":0,"user_id":3,"token":""}
```
注册失败
```shell
{"status_code":1,"status_msg":"user already exits","token":""}
```

### 视频推流的实现

> 基于`minio`和`redis`

1. 启动`mysql`, 并执行[init.sql](pkg/configs/sql/init.sql), 初始化数据库
2. 我们用单元测试的方式插入一些基础数据, 执行[random_data_test.go](tests/random_data_test.go), 执行之前, 请按照群里的要求配置好环境变量
3. 快乐的测试吧