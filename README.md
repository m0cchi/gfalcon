# identification infrastructure for SaaS Platform
[![CircleCI](https://circleci.com/gh/m0cchi/gfalcon.svg?style=svg)](https://circleci.com/gh/m0cchi/gfalcon)
[![BCH compliance](https://bettercodehub.com/edge/badge/m0cchi/gfalcon?branch=master)](https://bettercodehub.com/)
[![codebeat badge](https://codebeat.co/badges/c50993b8-1cbc-4e6a-ad94-43c31c30d020)](https://codebeat.co/projects/github-com-m0cchi-gfalcon-master)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/d683d6b0287b4d7cb11c6a2893768006)](https://www.codacy.com/app/boom.boom.planet/gfalcon?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=m0cchi/gfalcon&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/m0cchi/gfalcon)](https://goreportcard.com/report/github.com/m0cchi/gfalcon)

# Description
SaaS Platform(e.g. AWS)で利用するための認証基盤。

# Feature
- SaaS Platform で利用することを想定
   - Service 毎に権限の定義できる
   - シングルログアウトができる
   - 認証連携のサポート(TODO)
- Easy for any language(TODO: gRPC で API の呼び出し)

# TODO
- Role の作成
- 認証連携 (SAML?OpenID Connet?Kerberos?)
- gRPC で API の呼び出し
- Log
- ベンチマーク
- SampleApp(Wiki,ImageStorage)
- CLI経由でモデル操作
- ModelInfoの実装(UserやServiceを任意に拡張できる)
- 無効なSessionの削除
- 設定ファイル(セッションの長さや有効期限の設定)

# Sample IdP
## [gfalcon-signin-service](https://github.com/m0cchi/gfalcon-signin-service)
IdP の Sample Application。
認証の利用例となる。
また、SSO の実現も行なっている。
Cookie をサブドメイン間で共有することで SSO を実現している。
サブドメイン群の中に、非SSLのWebServiceが動いている場合は非推奨。
### 参考となるコード: [gfalconにおける認証の仕方](https://github.com/m0cchi/gfalcon-signin-service/blob/master/app/server.go#L58-L71)
### SSO の例: [gfalcon-sso-sample](https://github.com/m0cchi/gfalcon-sso-sample)
![](https://i.gyazo.com/1cde44d51b4356e8cedbc8029b9be131.gif)
### Demo
- Team/UserID/Password: gfalcon/gfadmin/secret
- IdP: https://saas.m0cchi.net/
- SP:  https://note.m0cchi.net/

# Sample SP
## [gfalcon-action-control-sample](https://github.com/m0cchi/gfalcon-action-control-sample)
SP の Sample Application。
認可の利用例となる。
Action(利用制限を行いたい行動) を定義する。
Action を許可したい User は、Action との ActionLink を作成する。
SP では、 ActionLink で認可の制御を実現している。
### 参考となるコード: [gfalconにおける認可の仕方](https://github.com/m0cchi/gfalcon-action-control-sample/blob/master/server.go#L93-L97)
### Demo
- Team/UserID/Password: gfalcon/sahohime/secret
- https://plank.saas.m0cchi.net/

# License
Licensed under the MIT License.
