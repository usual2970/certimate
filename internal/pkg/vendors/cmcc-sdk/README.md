移动云 Go sdk 文档: [https://ecloud.10086.cn/op-help-center/doc/article/53799](https://ecloud.10086.cn/op-help-center/doc/article/53799)

移动云 Go sdk 下载地址: [https://ecloud.10086.cn/api/query/developer/nexus/service/rest/repository/browse/go-sdk/gitlab.ecloud.com/ecloud/](https://ecloud.10086.cn/api/query/developer/nexus/service/rest/repository/browse/go-sdk/gitlab.ecloud.com/ecloud/)

接口文档地址: [https://ecloud.10086.cn/op-help-center/doc/outline/80900](https://ecloud.10086.cn/op-help-center/doc/outline/80900)

将其引入本地目录的原因是:

1. 此包必须通过 cmcc 私有仓库获取, 为构建带来不便, 故直接引入
2. 原始包存在部分内容错误, 需要自行修改, 如
   * 存在一些编译错误
   * 返回错误的时候, 未返回错误信息
   * 解析响应体错误
3. 将一些无用的目录和文件剔除了
