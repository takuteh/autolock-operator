# Autolock Kubernetes Operator 仕様

## 1. 概要

Autolock Operator は、Autolock CR（Custom Resource）をもとに Kubernetes リソースを管理する。

Webapp から Autolock CR を更新し、Operator が変更を検知して ConfigMap や各種リソースに反映する。設定変更時には、変更された設定を参照する Pod のみをロールアウトする。

### 基本方針

- Autolock CR をアプリケーション設定の管理元とする。
- Kubernetes リソースの作成・更新は Operator が担当する。
- Webapp は Kubernetes API 経由で Autolock CR を更新する。
- ConfigMap や通知用 Secret の変更に応じて、必要な Deployment をロールアウトする。
- 認証 DB の認証情報および初期化用設定は、通常のアプリケーション設定とは分離する。
- Operator が管理するリソースは、Reconcile によって望ましい状態を維持する。
- リソースに差分がない場合は不要な更新を行わず、不要な Pod の再起動を防止する。
- 手動管理のリソースは、Operator による自動作成・更新・削除の対象としない。ただし、変更イベントを監視するリソースについては、仕様で定めた処理を行う。

## 2. 管理対象リソース

本仕様における「変更検知時の再調整対象」は、対象リソースの変更を検知した際に、Operator が再調整するリソースを表す。

「ロールアウト対象」は、設定変更に伴い Pod の再作成を必要とする Deployment を表す。Ingress や Service の変更は、通常、それ自体の更新によって Pod のロールアウトを発生させない。

Operator が管理するリソースについても、実際に差分がある場合のみ更新する。

### 2.1 Main

Main は Autolock の主要なロジックを実行し、共通設定と通知用 Secret を参照する。  
メインロジックについてはこちらを参照→[詳細](https://github.com/takuteh/autolock)
| リソース | 名称 | 作成・管理方法 | 変更検知時の再調整対象 | 用途 |
|---|---|---|---|---|
| Deployment | `autolock-main` | Operator が管理 | なし | Main アプリケーションを実行する。共通設定または通知用 Secret の変更時にロールアウトする |
| ConfigMap | `autolock-config` | Operator が管理 | Main / Webhook / Webapp の各 Deployment | Main のロジック設定を保持する |
| Secret | `autolock-secret` | 手動作成 | Main Deployment | LINE・Slack の通知先情報やトークンを保持する。変更時は Main のみロールアウトする |

### 2.2 MQTT Broker

MQTT Broker は MQTT 通信を提供し、Service を通じて Main などのアプリケーションから接続できるようにする。

| リソース | 名称 | 作成・管理方法 | 変更検知時の再調整対象 | 用途 |
|---|---|---|---|---|
| Deployment | `autolock-mqtt-broker` | Operator が管理 | なし | MQTT Broker を実行する |
| Service | `autolock-mqtt-broker` | Operator が管理 | なし | MQTT Broker への接続先を提供する |

### 2.3 Webhook

Webhook は LINE の Webhook リクエストを受け取り、内容を整形して Main に通知する。

共通設定を ConfigMap から取得し、設定変更時には必要に応じてロールアウトする。

| リソース | 名称 | 作成・管理方法 | 変更検知時の再調整対象 | 用途 |
|---|---|---|---|---|
| Deployment | `autolock-webhook` | Operator が管理 | なし | Webhook アプリケーションを実行する |
| Service | `autolock-webhook` | Operator が管理 | なし | Webhook への接続先を提供する |
| Ingress | `autolock-webhook` | Operator が管理 | なし | 外部から Webhook へのアクセスを受け付ける |
| ConfigMap | `autolock-config` | Operator が管理 | Main / Webhook / Webapp の各 Deployment | Webhook に共通設定を提供する |

### 2.4 Webapp

Webapp は共通設定を ConfigMap から取得し、設定を変更する場合は Kubernetes API 経由で Autolock CR を更新する。

Webapp には Autolock CR の読み取り・更新に必要な RBAC 権限を付与する。ConfigMap や Deployment の直接変更は Operator が担当する。

| リソース | 名称 | 作成・管理方法 | 変更検知時の再調整対象 | 用途 |
|---|---|---|---|---|
| Deployment | `autolock-webapp` | Operator が管理 | なし | Webapp を実行する |
| Service | `autolock-webapp` | Operator が管理 | なし | Webapp への接続先を提供する |
| Ingress | `autolock-webapp` | Operator が管理 | なし | 外部から Webapp へのアクセスを受け付ける |
| ConfigMap | `autolock-config` | Operator が管理 | Main / Webhook / Webapp の各 Deployment | Webapp に共通設定を提供する |
| ServiceAccount | `autolock-webapp` | Operator が管理 | なし | Webapp が Kubernetes API を利用する際の実行主体 |
| Role | `autolock-cr-editor` | Operator が管理 | なし | Autolock CR の読み取り・更新権限を定義する |
| RoleBinding | `autolock-webapp-to-cr-editor` | Operator が管理 | なし | Role と ServiceAccount を紐付ける |

### 2.5 認証 DB

認証 DB は StatefulSet として管理し、PersistentVolumeClaim を利用してデータを永続化する。

認証情報用 Secret および初期化用 ConfigMap は手動作成とし、DB の接続・初期化に使用する。これらの変更による自動再起動・再初期化は行わない。

| リソース | 名称 | 作成・管理方法 | 変更検知時の再調整対象 | 用途 |
|---|---|---|---|---|
| StatefulSet | `autolock-auth-db` | Operator が管理 | なし | 認証 DB を実行する |
| Service | `autolock-auth-db` | Operator が管理 | なし | 認証 DB への接続先を提供する |
| PersistentVolumeClaim | StatefulSet の命名規則に従う | StatefulSet の `volumeClaimTemplates` により作成 | なし | DB データを永続化する |
| Secret | `autolock-auth-db-secret` | 手動作成 | なし | DB の認証情報を提供する |
| ConfigMap | `autolock-auth-db-init` | 手動作成 | なし | DB の初期化に使用する |

認証 DB 用 Secret と初期化用 ConfigMap は、通常のアプリケーション設定とは独立して管理する。変更時の対応については「5. 認証 DB 用 Secret と初期化用 ConfigMap」を参照する。

## 3. ConfigMap の管理

### 3.1 用途

`autolock-config` には、Main のロジックに関する共通設定を格納する。

設定の管理元は Autolock CR とし、Operator が CR の内容に基づいて ConfigMap を作成・更新する。原則として手動で変更しない。

参照するコンポーネントは以下のとおり。

- Main
- Webhook
- Webapp

### 3.2 更新フロー

1. Webapp などから Autolock CR が更新される。
2. Operator が Autolock CR の変更を検知する。
3. Reconcile により ConfigMap の内容を CR の設定値に合わせる。
4. CR の設定値から、対象コンポーネントが参照する設定内容のハッシュ値を計算する。
5. Main、Webhook、Webapp の各 Deployment の Pod テンプレートにハッシュ値を設定する。
6. Pod テンプレートのハッシュ値が変化した Deployment のみロールアウトする。

ConfigMap の内容に変更がない場合、不要な Deployment の更新や Pod の再起動は行わない。

### 3.3 更新時の注意事項

- ConfigMap は Operator が管理するため、手動変更は原則として行わない。
- ConfigMap が手動で変更された場合は、Operator が CR に基づく望ましい状態へ戻す。
- 設定値の変更がない場合は、不要なロールアウトを発生させない。
- ConfigMap の内容と Pod テンプレートに設定するハッシュ値が一致するように管理する。

## 4. 通知用 Secret の管理

### 4.1 用途

`autolock-secret` には、LINE および Slack による通知に必要な送信先情報やトークンを格納する。

この Secret は Main のみが参照する。

Secret 自体は手動作成・管理とし、Operator は内容を変更せず、変更検知とハッシュ値の計算に使用する。

### 4.2 更新フロー

1. Operator が `autolock-secret` の変更イベントを検知する。
2. Secret に対応する Autolock CR を特定し、Reconcile を要求する。
3. Operator が Secret の内容からハッシュ値を計算する。
4. Main Deployment の Pod テンプレートにハッシュ値を設定する。
5. ハッシュ値が変化した場合、Main のみをロールアウトする。

Webhook および Webapp は、この Secret の変更によるロールアウト対象としない。

### 4.3 Secret と Autolock CR の対応付け

Secret と Autolock CR の対応関係を明示し、CR 名をコード内に固定値として埋め込まない設計とする。

対応付けには Secret の Annotation を利用する。Annotation のキーは、例えば以下とする。

`autolock.takuteh.uk/autolock-cr`

Annotation の値には対応する Autolock CR の名前を設定する。

Operator は Secret の Namespace と Annotation の値をもとに対応する CR を特定する。対応する CR が存在しない場合や Annotation が設定されていない場合は、誤った CR を更新せず、適切にログを出力する。

## 5. 認証 DB 用 Secret と初期化用 ConfigMap

認証 DB の認証情報用 Secret および初期化用 ConfigMap は、DB の接続・初期化を目的とし、Main のロジック設定とは分離する。

- これらの変更を自動 Reconcile のトリガーとしない。
- DB 初期化用 ConfigMap は初期化時に使用する。
- Operator は DB の StatefulSet、Service、永続ストレージなどを管理する。
- 認証情報や初期化設定の変更に伴う DB の再初期化・再起動は自動では行わない。
- 変更が必要な場合は、DB の状態やデータへの影響を確認したうえで、別途運用手順に従って対応する。

## 6. Webapp からの Autolock CR 更新

Webapp は Kubernetes API を通じて Autolock CR を更新する。

### 6.1 RBAC 構成

Operator は以下のリソースを作成・管理する。

- ServiceAccount：`autolock-webapp`
- Role：`autolock-cr-editor`
- RoleBinding：`autolock-webapp-to-cr-editor`

### 6.2 権限方針

Webapp には Autolock CR の読み取り・更新に必要な最小限の権限を付与する。

- CR の取得・一覧取得・監視が必要な場合は、必要な読み取り権限を付与する。
- CR の更新には `update`、`patch` など必要な権限を付与する。
- ConfigMap や Deployment を直接変更する権限は原則として付与しない。
- Operator が設定の反映と Pod のロールアウトを担当する。

ServiceAccount、Role、RoleBinding は Operator が望ましい状態を管理し、差分がある場合に修正する。

## 7. Reconcile 方針

### 7.1 Reconcile のトリガー

| イベント | Reconcile の目的 | ロールアウト対象 |
|---|---|---|
| Autolock CR の作成・更新 | CR の設定を各リソースに反映する | 共通設定のハッシュ値が変化した Deployment |
| `autolock-config` の変更 | 対応する CR を再調整し、ConfigMap を望ましい状態に戻す | 必要に応じて Main / Webhook / Webapp |
| `autolock-secret` の変更 | 通知設定を使用する Main を再調整する | Main |
| 認証 DB 用 Secret の変更 | 自動 Reconcile のトリガーとしない | なし |
| DB 初期化用 ConfigMap の変更 | 自動 Reconcile のトリガーとしない | なし |

ConfigMap や通知用 Secret の変更を Reconcile の契機にする場合は、対象リソースと Autolock CR の対応関係を解決し、該当する CR のみを再調整する。

なお、Operator が管理する ConfigMap については、CR を管理元とする。ConfigMap の手動変更を検知した場合も、最終的には CR の設定値に一致させる。

### 7.2 ハッシュによるロールアウト

設定のハッシュ値を Deployment の Pod テンプレートの `annotations` に設定する。

- ConfigMap のハッシュ：Main / Webhook / Webapp
- 通知用 Secret のハッシュ：Main

Pod テンプレートのハッシュ値が変化すると、Kubernetes の Deployment Controller が新しい ReplicaSet を作成し、Pod を順次更新する。

ハッシュ値が変わらない場合は、不要なロールアウトを発生させない。

Ingress、Service、RBAC などの変更は、Pod テンプレートの変更とは分けて管理する。これらのリソースは、必要な差分がある場合にそれぞれ更新する。

### 7.3 Reconcile の責務

Reconcile は、Autolock CR に記載された望ましい状態と Kubernetes 上の実際の状態を一致させる。

- イベントの検知とリソースの更新処理を分離する。
- 変更イベントに応じて、必要なリソースだけを更新する。
- 既存リソースは可能な限り更新し、不要な削除・再作成を避ける。
- 差分がない場合は更新を行わない。
- Operator が管理する Deployment、Service、Ingress、StatefulSet、ServiceAccount、Role、RoleBinding などは、Reconcile によって望ましい状態を維持する。
- DB 用 Secret および初期化用 ConfigMap は、仕様で定めたとおり自動変更検知・再初期化の対象外とする。
- Secret などの変更イベントから CR を特定する際は、Annotation などの明示的な対応関係を利用する。
- 設定のハッシュ値は対象コンポーネントが参照する内容に基づいて計算し、無関係な設定変更によるロールアウトを避ける。
- リソースの削除や再作成は、永続データや外部接続への影響を考慮し、必要な場合に限定する。

# インストール手順

以下では、`v1.0.0` を使用する例を示します。使用するバージョンに応じて URL 内のタグを変更してください。

## 1. Operator のインストール

GitHub Releases で公開しているインストールマニフェストを適用します。

```bash
kubectl apply -f https://github.com/takuteh/autolock-operator/releases/download/v1.0.0/install.yaml
```

## 2. Namespace の作成

Autolock を配置する Namespace を作成します。

```bash
kubectl create namespace autolock
```

すでに存在する場合は、この手順を省略してください。

## 3. Secret の準備

認証 DB 用 Secret と Autolock 用 Secret のテンプレートをダウンロードします。

```bash
curl -fLO https://raw.githubusercontent.com/takuteh/autolock-operator/main/config/secrets/authdb-secret.yaml
curl -fLO https://raw.githubusercontent.com/takuteh/autolock-operator/main/config/secrets/autolock-secret.yaml
```

ダウンロードしたファイルを編集し、実際の環境に合わせて認証情報を設定してください。編集後、Secret を作成します。

```bash
kubectl apply -f authdb-secret.yaml
kubectl apply -f autolock-secret.yaml
```

Secret は `autolock` Namespace に作成してください。

**注意:** 実際のパスワードやトークンを公開リポジトリにコミットしないでください。

## 4. Autolock の Custom Resource の準備

サンプル CR をダウンロードします。

```bash
curl -fLO https://raw.githubusercontent.com/takuteh/autolock-operator/main/config/samples/autolock_v1alpha1_autolock.yaml
```

CR の設定値を環境に合わせて編集します。
編集後、CR を適用します。

```bash
kubectl apply -f autolock_v1alpha1_autolock.yaml
```

Operator は CR の内容に基づいて、Autolock に必要な Kubernetes リソースを作成・管理します。

## 5. 起動状態の確認

各リソースの状態を確認します。

```bash
kubectl get all -n autolock
kubectl get ingress -n autolock
kubectl get autolock -n autolock
```

Pod が正常に起動し、Ready 状態になっていることを確認してください。

起動しない場合は、CR の設定値、Secret の内容、Operator のログを確認してください。



Pod が正常に起動し、Ready 状態になっていることを確認してください。立ち上がらない場合はsecretかCRが正しく適用されていない可能性があります。
