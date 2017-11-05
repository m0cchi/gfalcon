# identification infrastructure for SaaS Platform
[![CircleCI](https://circleci.com/gh/m0cchi/gfalcon.svg?style=svg)](https://circleci.com/gh/m0cchi/gfalcon)

# TODO
- 認可
- Role の作成
- SSO (SAML?)
- gRPC で API の呼び出し
- Log
- ベンチマーク
- SampleApp(Wiki,ImageStorage)

# Sample App
## https://github.com/m0cchi/gfalcon-signin-service
IdP の Sample Application.
Cookie をサブドメイン間で共有することで SSO を実現する。
サブドメイン群の中に、非SSLのWebServiceが動いている場合は非推奨。