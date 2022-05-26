# ByteDance-BackendProject-Douyin
字节后端训练营-抖音demo项目

在这个分支上完成了docker-compose中的改造。

####（默认大家都有docker环境）

首先git clone下来后切换分支到docker-compose：
`git checkout docker-compose`

我的docker-compose在docker网络下通信的，所以先要新建个docker network（不想写build.sh）
`sudo docker network create douyin-net`

我的镜像仓库是托管在腾讯云的，其中还使用了oss服务（本人开了个test用户，所以在一定时间内其他人也可以用，但自己部署的话可以根据自己的情况换成自己的，因此账户不保活）

这时候可能找不到仓库，但我试过，先pull单独一个镜像下来后，再compose就能解析到了（可选）
`sudo docker pull ccr.ccs.tencentyun.com/douyin-test/user-srv:v1.0`
然后再compose
`sudo docker-compose -f docker-compose.yaml up`

然后差不多就可以了，服务暴露在 `http://localhost:8080` 上

但其实关于数据库还没整完，后面会update自己的一个sql，mysql我没部署在docker上，因为我用mac m1开发的，mysql镜像和linux有些许不同，其实想要把mysql加到compose的其实也不难。
