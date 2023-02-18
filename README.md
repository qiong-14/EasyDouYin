# EasyDouYin

### 目录结构

```
EasyDouyin 
├── /biz/ 路由相关
|   ├── /handler/ 路由处理函数
|   ├── /resp/    响应结构体
|   ├── /router/  路由注册
├── /configs/ 配置文件
├── /constants/ 常量
├── /dal/ 数据库操作
├── /middleware/ 中间件相关
│   ├── jwt/ 鉴权
│   ├── minio/ 文件储存系统
├── /service/ 服务层 (暂时用不上)
├── /tools/ 公用函数
├── .gitignore
├── go.mod
├── main.go
├── README.md
├── Makefile
```

## 抖音项目运行步骤

下载安装`ffmpeg`, 用于视频封面的获取

```shell
## for macos
brew install ffmpeg
## for windows, you need to fix env variables(PATH) to include ffmpeg
winget install ffmpeg
## for ubuntu, debian, w.r.t. apt
sudo apt-get install -y ffmpeg
## others
todo
```

下载项目依赖

```shell
go mod tidy
```

运行docker

```shell
docker-compose up
```

写入环境变量

```shell
source dy_secure_config.sh
```

* windows 必须要用**cmd终端**，运行以下bat程序

```shell
dy_secure_config.bat
```

### 视频推流的实现

**启动docker后运行一次即可**

基于`minio`和`redis`,我们用单元测试的方式插入一些基础数据, **执行之前, 请按照群里的要求配置好环境变量**

```shell
go apitest tests/random_data_test.go
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
## you need replace `($token)` into real token 
curl --location --request GET 'localhost:8080/douyin/ping' --header 'Authorization: Bearer ($token)'
## for example:
curl --location --request GET 'localhost:8080/douyin/ping' --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzYwNDU1ODYsImlkZW50aXR5IjoicmVhZHlnbyIsIm9yaWdfaWF0IjoxNjc2MDQxOTg2fQ.BIMU_OS2CLrmmN1vrW0XWkFwaPPu5gPtViBAnw-lXK4'
# {"message":"user_id:4"}
```

### 关于Redis的使用

See in [Redis使用文档](./UseRedis.md)
