# EasyDouYin

## 抖音项目运行步骤
### 依赖下载
windows
```shell
set GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
```
运行docker
```shell
docker-compose up
```
写入环境变量
* linux or macos 直接执行以下命令
```shell
source dy_secure_config.sh
```
* windows 必须要用**cmd终端**，运行以下bat程序
```cmd
dy_secure_config.bat
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

### 用户注册
基于Hertz.JWT鉴权的注册，注册后调用了mw.JwtMiddleware.LoginHandler进行登录
[JWT认证 | CloudWeGo](https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/middleware/jwt/)
```shell
curl --request POST 'http://localhost:8080/douyin/user/register/?username=704788475&password=111111'
# {"status_code":1,"status_msg":"user already exits","token":""}

curl --request POST 'http://localhost:8080/douyin/user/register/?username=readygo&password=111111'
# {"status_code":0,"status_msg":"login success","user_id":10,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzYwNDU0NzAsImlkZW50aXR5IjoicmVhZHlnbyIsIm9yaWdfaWF0IjoxNjc2MDQxODcwfQ.3G05OinRGLYDGlsDz5zt4XJX4UnjW6XnILRk1SvK2gM"}
```
### 用户登录
```shell
curl --request POST 'http://localhost:8080/douyin/user/login/?username=readygo11&password=111111'
# {"status_code":401,"status_msg":"record not found","token":""}

 curl --request POST 'http://localhost:8080/douyin/user/login/?username=readygo&password=111111'
# {"status_code":0,"status_msg":"login success","user_id":10,"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzYwNDU1ODYsImlkZW50aXR5IjoicmVhZHlnbyIsIm9yaWdfaWF0IjoxNjc2MDQxOTg2fQ.BIMU_OS2CLrmmN1vrW0XWkFwaPPu5gPtViBAnw-lXK4"}
```

### 鉴权后的ping
```shell
curl --location --request GET 'localhost:8080/douyin/ping' --header 'Authorization: Bearer ($token)'
例：
curl --location --request GET 'localhost:8080/douyin/ping' --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzYwNDU1ODYsImlkZW50aXR5IjoicmVhZHlnbyIsIm9yaWdfaWF0IjoxNjc2MDQxOTg2fQ.BIMU_OS2CLrmmN1vrW0XWkFwaPPu5gPtViBAnw-lXK4'
# {"message":"user_id:4"}
```