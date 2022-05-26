# ByteDance-BackendProject-Douyin
字节后端训练营-抖音demo项目

### 1. master分支里做的是grpc+gin（网关）+nacos做配置中心和服务发现
主分支完成了所有逻辑，以及服务的分离，使用了mysql，nacos（做服务发现和配置中心。因为需要更多配置才能用，所以不写这个版本的build了）

### 2. docker-compose分支做的是将项目部署在docker-compose上，没有考虑多个微服务就不做服务发现了（TODO：k8s；先用dns做负载均衡，打算用istio做服务治理和mesh）
这个版本大部分人可以直接运行但是mysql配置没有完成，也不难。具体操作切到这个分支里。
