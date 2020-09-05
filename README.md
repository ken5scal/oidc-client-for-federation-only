# oidc-client-for-federation-only
federation only oidc client

CloudFrontでホストされるSPAの前段にAzure ADへの認証機構を作成しようとしたものの、
Lambda関数がサポートする[ランタイムにGoがなく、また、Viewer Requestに対して設定するLambda関数の最大圧縮サイズに１MB制限](https://docs.aws.amazon.com/ja_jp/AmazonCloudFront/latest/DeveloperGuide/lambda-requirements-limits.html)
があるため、断念。
供養。

## Ref
* [Microsoft identity platform and OpenID Connect protocol](https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-protocols-oidc)