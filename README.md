# gocon-k8s-server
Kubernetes上で動作させるGoアプリケーションのサンプルリポジトリです。

------------------

## Requirement

### Docker
[install](https://docs.docker.com/v17.09/engine/installation/#docker-cloud)

------------------

# Usage 

## Docker

### Dockerイメージのビルド

```cassandraql
$ docker image build ./app
```

### Dockerコンテナの起動
```cassandraql
 docker run -d -p 8080:8080 app 
```

#### 疎通確認
```cassandraql
$ curl localhost:8080
Hello, World!   
```

### Dockerコンテナの停止
```cassandraql
$ docker container ls |grep app| awk '{print $1}' |xargs  docker stop
```
------------------
## Kubernetes 
イメージは上記の手順で準備されているので、deployのみを行う

[こちらの記事](http://sagantaf.hatenablog.com/entry/2019/07/27/210008)を参考にDocker for macでKuberentesを利用できる、`kubectl`コマンドを利用できる状態にしておいてください

```cassandraql
$ kubectl apply -f manifest/deployment.yaml
$ kubectl get pod
NAME                              READY   STATUS    RESTARTS   AGE
app-deployment-5c746db894-8tz6d   1/1     Running   0          2m57s
app-deployment-5c746db894-dwwsh   1/1     Running   0          2m57s
app-deployment-5c746db894-x8zdg   1/1     Running   0          2m57s

```

#### Podのログの確認
```cassandraql
$ kubectl get pod | awk 'NR==2' | awk '{print $1}' |xargs kubectl logs
2019-10-27 UTC
{"severity":"INFO","timestamp":"2019-10-27T07:49:20.948Z","caller":"gocon-k8s-server/main.go:70","msg":"server start serving"}

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.1.11
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8080

```

### PodのExpose
```cassandraql
$ kubectl expose deployment app --type=NodePort
$ kubectl get service
NAME             TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
app-deployment   NodePort    10.105.62.219   <none>        8080:30441/TCP   5m16s
kubernetes       ClusterIP   10.96.0.1       <none>        443/TCP          2d19h
```

上記の例ではlocalhostの`30441`番portで公開されているのでそのポートに対してリクエストを送ることで、Podに対してリクエストが行える

```cassandraql
$ curl localhost:30441
Hello, World!
```


### リソースの削除
```cassandraql
$ kubectl delete deployment app-deployment
$ kubectl delete service  app-deployment
```
