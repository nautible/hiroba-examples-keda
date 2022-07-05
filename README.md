# KEDA動作確認用アプリケーション

KEDAの動作確認用サンプルアプリケーションになります。

なお、本サンプルアプリケーションはAWS上での動作を前提としています。

## 1. AWSリソースの準備

### 1.1 ECR

管理コンソールよりサンプルアプリケーション用のECRリポジトリを作成します。（リポジトリ名：demo-keda）

### 1.2 SQS

管理コンソールよりサンプルアプリケーション用のSQSを作成します。（キュー名：demo-keda-queue）

## 2. KEDAの導入

[nautible-pluginのpod-autoscaler](https://github.com/nautible/nautible-plugin/tree/main/pod-autoscaler)を参考にKEDAをクラスタに導入します。

## 3. サンプルアプリケーションの準備

サンプルアプリケーションのコンテナイメージを準備します。

### 3.1 サンプルアプリケーションのビルド

```bash
cd pod-autoscaler/sample
docker build -t demo-keda:latest .
docker tag demo-keda:latest <レジストリ名>/demo-keda:latest
```

### 3.2 サンプルアプリケーションのプッシュ（ECR利用時の手順）

ECRへログイン

```bash
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin <レジストリ名>
```

ECRへイメージをプッシュ

```bash
docker push <レジストリ名>/demo-keda:latest
```

## 4. サンプルアプリケーションの導入

環境ごとの設定をマニフェストに設定し、反映します。

### 4.1 マニフェストに環境ごとの設定値を定義

イメージの指定及び環境変数（ACCOUNT_ID,REGION）を指定します。

manifests/deployment.yaml

```yaml
...
    spec:
      containers:
      - name: demo-keda
        image: <レジストリ名>/<リポジトリ名>:<バージョン>
        env:
        - name: ACCOUNT_ID
          value: "<AWSアカウントID>"
        - name: REGION
          value: "<リージョン>"
```

### 4.2 マニフェストを反映

```bash
cd manifests
kubectl apply -f .
```

## 4. 確認

AWS管理コンソール上でSQSにテストデータを送信し、しばらく待つとPodが起動してログが出力されます。

SQSに「test1」を送信した場合の例

```text
kubectl get po
NAME                         READY   STATUS    RESTARTS   AGE
demo-keda-57d484c7cb-lqxc7   1/1     Running   0          14s

kubectl logs demo-keda-57d484c7cb-lqxc7 -f
receive message ...
test1
```

また、その後しばらく放置すると、Podが自動的に停止します。（DeploymentのPod要求数が0になる）

## 5. サンプルアプリケーションの削除

```bash
cd pod-autoscaler/sample/manifests
kubectl delete -f .
```

ECRおよびSQSについてはAWS管理コンソールから削除します。
