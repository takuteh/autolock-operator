# Autolock Kubernetes Operator 仕様

## 1. 概要

Autolock Operator は、Autolock CR（Custom Resource）をもとに Kubernetes リソースを管理する。

Webapp から Autolock CR を更新し、Operator が変更を検知して ConfigMap や各種リソースに反映する。設定変更時には、変更された設定を参照する Pod のみをロールアウトする。

### 基本方針

- Autolock CR をアプリケーション設定の管理元とする。
- Kubernetes リソースの作成・更新は Operator が担当する。
- Webapp は Kubernetes API 経由で Autolock CR を更新する。
- ConfigMap や Secret の変更に応じて、必要な Deployment のみをロールアウトする。
- 認証 DB の認証情報および初期化用設定は、通常のアプリケーション設定とは分離する。
- Operator が管理するリソースは、Reconcile によって望ましい状態を維持する。

## 2. 管理対象リソース

「Reconcile 対象」は、リソースの管理方法を表す。`deployment` は Deployment の Pod テンプレート更新を伴う対象、`configmap` は ConfigMap の更新対象、`-` は個別の設定変更によるロールアウト対象ではないことを示す。

なお、Operator が管理するリソースについても、実際に差分がある場合のみ更新する。

### 2.1 Main

Main は Autolock の主要なロジックを実行し、共通設定と通知用 Secret を参照する。

| リソース | 作成方法 | Reconcile 対象 | 用途 |
|---|---|---|---|
| Deployment `autolock-main` | Operator が自動作成 | - | Main アプリケーションを実行する。共通設定または通知用 Secret の変更時にロールアウトする |
| ConfigMap `autolock-configmap` | Operator が自動作成 | configmap | Main のロジック設定を保持する |
| Secret `autolock-secret` | 手動作成 | deployment | LINE・Slack の通知先情報やトークンを保持する。変更時は Main のみロールアウトする |

### 2.2 MQTT Broker

MQTT Broker は MQTT 通信を提供し、Service を通じて Main などのアプリケーションから接続できるようにする。

| リソース | 作成方法 | Reconcile 対象 | 用途 |
|---|---|---|---|
| Deployment `autolock-mqtt-broker` | Operator が自動作成 | - | MQTT Broker を実行する |
| Service `autolock-mqtt-broker` | Operator が自動作成 | - | MQTT Broker への接続先を提供する |

### 2.3 Webhook

Webhook は LINE の Webhook リクエストを受け取り、内容を整形して Main に通知する。

共通設定を ConfigMap から取得し、設定変更時には必要に応じてロールアウトする。

| リソース | 作成方法 | Reconcile 対象 | 用途 |
|---|---|---|---|
| Deployment `autolock-webhook` | Operator が自動作成 | - | Webhook アプリケーションを実行する |
| Service `autolock-webhook` | Operator が自動作成 | - | Webhook への接続先を提供する |
| Ingress | Operator が自動作成 | - | 外部から Webhook へのアクセスを受け付ける |
| ConfigMap `autolock-configmap` | Operator が自動作成 | deployment | Webhook に共通設定を提供する |

### 2.4 Webapp

Webapp は共通設定を ConfigMap から取得し、設定を変更する場合は Kubernetes API 経由で Autolock CR を更新する。

Webapp には Autolock CR の読み取り・更新に必要な RBAC 権限を付与する。ConfigMap や Deployment の直接変更は Operator が担当する。

| リソース | 作成方法 | Reconcile 対象 | 用途 |
|---|---|---|---|
| Deployment `autolock-webapp` | Operator が自動作成 | - | Webapp を実行する |
| Service `autolock-webapp` | Operator が自動作成 | - | Webapp への接続先を提供する |
| Ingress | Operator が自動作成 | - | 外部から Webapp へのアクセスを受け付ける |
| ConfigMap `autolock-configmap` | Operator が自動作成 | deployment | Webapp に共通設定を提供する |
| ServiceAccount `autolock-webapp` | Operator が自動作成 | - | Webapp が Kubernetes API を利用する際の実行主体 |
| Role `autolock-cr-editor` | Operator が自動作成 | - | Autolock CR の読み取り・更新権限を定義する |
| RoleBinding | Operator が自動作成 | - | Role と ServiceAccount を紐付ける |

### 2.5 認証 DB

認証 DB は StatefulSet として管理し、PersistentVolumeClaim を利用してデータを永続化する。

認証情報用 Secret および初期化用 ConfigMap は手動作成とし、DB の接続・初期化に使用する。これらの変更による自動再起動・再初期化は行わない。

| リソース | 作成方法 | Reconcile 対象 | 用途 |
|---|---|---|---|
| StatefulSet `autolock-auth-db` | Operator が自動作成 | - | 認証 DB を実行する |
| Service `autolock-auth-db` | Operator が自動作成 | - | 認証 DB への接続先を提供する |
| PersistentVolumeClaim | Operator が自動作成 | - | DB データを永続化する |
| 認証情報用 Secret | 手動作成 | - | DB の認証情報を提供する |
| 初期化用 ConfigMap | 手動作成 | - | DB の初期化に使用する |

## 3. ConfigMap の管理

### 3.1 用途

`autolock-configmap` には、Main のロジックに関する共通設定を格納する。

参照するコンポーネントは以下のとおり。

- Main
- Webhook
- Webapp

### 3.2 更新フロー

1. Webapp などから Autolock CR が更新される。
2. Operator が Autolock CR の変更を検知する。
3. Reconcile により ConfigMap の内容を CR の設定値に合わせる。
4. ConfigMap の内容からハッシュ値を計算する。
5. Main、Webhook、Webapp の各 Deployment の Pod テンプレートにハッシュ値を設定する。
6. ハッシュ値が変化した Deployment のみロールアウトする。

ConfigMap の内容に変更がない場合、不要な Deployment の更新や Pod の再起動は行わない。

ConfigMap は、原則として既存リソースを更新する。削除・再作成を必須とする設計にはしない。

## 4. 通知用 Secret の管理

### 4.1 用途

`autolock-secret` には、LINE および Slack による通知に必要な送信先情報やトークンを格納する。

この Secret は Main のみが参照する。

### 4.2 更新フロー

1. Operator が `autolock-secret` の変更イベントを検知する。
2. Secret に対応する Autolock CR の Reconcile を要求する。
3. Operator が Secret の内容からハッシュ値を計算する。
4. Main Deployment の Pod テンプレートにハッシュ値を設定する。
5. ハッシュ値が変化した場合、Main のみをロールアウトする。

Webhook および Webapp は、この Secret の変更によるロールアウト対象としない。

Secret と Autolock CR の対応関係を明示し、CR 名をコード内に固定値として埋め込まない設計とする。対応関係には、Secret の Annotation などのメタデータを利用する。

## 5. 認証 DB 用 Secret と初期化用 ConfigMap

認証 DB の認証情報用 Secret および初期化用 ConfigMap は、DB の接続・初期化を目的とし、Main のロジック設定とは分離する。

- これらの変更を自動 Reconcile のトリガーとはしない。
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
- RoleBinding：Role と ServiceAccount の紐付け

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
| Autolock CR の作成・更新 | CR の設定を各リソースに反映する | 共通設定に変更がある Deployment |
| `autolock-configmap` の変更 | 対応する CR を再調整する | Main / Webhook / Webapp のうち必要なもの |
| `autolock-secret` の変更 | 通知設定を使用する Main を再調整する | Main |
| 認証 DB 用 Secret の変更 | 自動 Reconcile のトリガーとしない | なし |
| DB 初期化用 ConfigMap の変更 | 自動 Reconcile のトリガーとしない | なし |

ConfigMap や Secret の変更を Reconcile の契機にする場合は、対象リソースと Autolock CR の対応関係を解決し、該当する CR のみを再調整する。

### 7.2 ハッシュによるロールアウト

設定のハッシュ値を Deployment の Pod テンプレートの `annotations` に設定する。

- ConfigMap のハッシュ：Main / Webhook / Webapp
- 通知用 Secret のハッシュ：Main

Pod テンプレートのハッシュが変化すると、Kubernetes の Deployment Controller が新しい ReplicaSet を作成し、Pod を順次更新する。

ハッシュが変わらない場合は、不要なロールアウトを発生させない。

### 7.3 Reconcile の責務

Reconcile は、Autolock CR に記載された望ましい状態と Kubernetes 上の実際の状態を一致させる。

- イベントの検知とリソースの更新処理を分離する。
- 変更イベントに応じて、必要なリソースだけを更新する。
- 既存リソースは可能な限り更新し、不要な削除・再作成を避ける。
- 差分がない場合は更新を行わない。
- Operator が管理する Deployment、Service、Ingress、StatefulSet、ServiceAccount、Role、RoleBinding などは、Reconcile によって望ましい状態を維持する。
- DB 用 Secret および初期化用 ConfigMapは、仕様で定めたとおり自動変更検知・再初期化の対象外とする。
