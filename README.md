# EasyDouYin

## 抖音项目服务端简单示例

具体功能内容参考飞书说明文档

运行docker
```shell
docker-compose up
```
运行程序
```shell
go build && ./EasyDouYin
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
