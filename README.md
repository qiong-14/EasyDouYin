# EasyDouYin

## 抖音项目运行步骤

运行docker
```shell
docker-compose up
```
写入环境变量
```shell
source dy_secure_config.sh
```
### 视频推流的实现
基于`minio`和`redis`,我们用单元测试的方式插入一些基础数据, **执行之前, 请按照群里的要求配置好环境变量**

```shell
go test tests/random_data_test.go
```
快乐的测试吧 运行程序
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

