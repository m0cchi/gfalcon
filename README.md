# identification infrastructure for SaaS Platform
[![CircleCI](https://circleci.com/gh/m0cchi/gfalcon.svg?style=svg)](https://circleci.com/gh/m0cchi/gfalcon)
[![BCH compliance](https://bettercodehub.com/edge/badge/m0cchi/gfalcon?branch=master)](https://bettercodehub.com/)
[![codebeat badge](https://codebeat.co/badges/c50993b8-1cbc-4e6a-ad94-43c31c30d020)](https://codebeat.co/projects/github-com-m0cchi-gfalcon-master)

# TODO
- Role の作成
- SSO (SAML?)
- gRPC で API の呼び出し
- Log
- ベンチマーク
- SampleApp(Wiki,ImageStorage)
- CLI経由でモデル操作
- ModelInfoの実装(UserやServiceを任意に拡張できる)
- 無効なSessionの削除
- 設定ファイル(セッションの長さや有効期限の設定)

# Sample IdP
## https://github.com/m0cchi/gfalcon-signin-service
IdP の Sample Application.
Cookie をサブドメイン間で共有することで SSO を実現する。
サブドメイン群の中に、非SSLのWebServiceが動いている場合は非推奨。

### Sample SSO App: https://github.com/m0cchi/gfalcon-sso-sample
![](https://i.gyazo.com/1cde44d51b4356e8cedbc8029b9be131.gif)
### Demo
- Team/UserID/Password: gfalcon/gfadmin/secret
- IdP: https://saas.m0cchi.net/
- SP:  https://note.saas.m0cchi.net/

# Authorization Sample SP
## https://github.com/m0cchi/gfalcon-action-control-sample
認可を行なっている。
User には、事前に利用可能な Action との ActionLink を作成する。
ActionLink で認可の制御を実現している。

### Demo
- Team/UserID/Password: gfalcon/sahohime/secret
- https://plank.saas.m0cchi.net/
