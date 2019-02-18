# gateway

API网关程序，用于实现多个服务访问控制，精简单个服务的功能。主要功能参见: [主要功能](#主要功能)

### 此项目关闭，迁移到[jademperor](https://github.com/jademperor)

## 安装使用

1. 下载二进制文件
2. 下载对应前端打包文件

## Todos

- [x] HTTP反向代理及负载均衡
- [x] API缓存支持POST，PUT等带有body的请求
- [ ] 横向扩展master-slave 模式
- [ ] 权限管理（RBAC）
- [x] 流量控制（令牌桶算法）
- [x] 插件模式（支持动态关闭与开启）
- [x] 插件配置及时更新，不需重启网关加载
- [ ] 支持Docker部署
- [ ] 添加更多的测试代码，hah

## 开发环境

* MongoDB
* Go1.11.1
* Node v1.10.1 (npm: 6.4.1)

## 主要功能

### 1. 代理
#### 1.1 URI直接代理
#### 1.2 URI组合代理
#### 1.3 Server代理
