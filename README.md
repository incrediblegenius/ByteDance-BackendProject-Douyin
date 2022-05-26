# ByteDance-BackendProject-Douyin
字节后端训练营-抖音demo项目

在这个分支上完成了服务在k8s中的改造。使用原生dns:///来做的grpc的负载均衡，ingress+nginx的暴露服务的负载均衡。

### TODO ： 加入istio做服务治理和和grpc负载均衡

## 默认大家都有docker环境，都有一个k8s的cluster，都会kubectrl **

首先git clone下来后切换分支到k8s：

`git checkout k8s`

分别执行各个微服务的deployment.yaml(七个文件夹下都有，镜像是已经在腾讯云上的公开仓库里)

`kubectl apply -f deployment.yaml`  x  7

我的镜像仓库是托管在腾讯云的，其中还使用了oss服务（本人开了个test用户，所以在一定时间内其他人也可以用，但自己部署的话可以根据自己的情况换成自己的，因此账户不保活）

关于服务暴露是使用了ingress+nginx，具体实现在参考网站 `https://kubernetes.github.io/ingress-nginx/deploy/#local-testing`。确实很简单，把quick start看几行就可以快速部署了，如果要详细部署要去把yaml拷出来自己改。

主要命令也在 `build.sh` 中，可以参考一下。

但其实关于数据库还没整完，后面会update自己的一个sql（或者切到 `Gin/model/CreateTable`下自己建表，mysql我没部署在docker上，因为我用mac m1开发的，mysql镜像和linux有些许不同，其实想要把mysql加到compose的其实也不难。
