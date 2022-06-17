# ByteDance-BackendProject-Douyin
字节后端训练营-抖音demo项目

### 1. master分支里做的是grpc+gin（网关）+nacos做配置中心和服务发现
主分支完成了所有逻辑，以及服务的分离，使用了mysql，nacos（做服务发现和配置中心。因为需要更多配置才能用，所以不写这个版本的build了）

### 2. docker-compose分支做的是将项目部署在docker-compose上，没有考虑多个微服务就不做服务发现了
这个版本大部分人可以直接运行但是mysql配置没有完成，也不难。具体操作切到这个分支里。

### 3. k8s分支里完成了k8s部署，istio，ingress组件的使用，负载均衡是在原本的dns负载均衡上的修改，做了轮训更新设置（todo：有机会写个监听）
