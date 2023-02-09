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
