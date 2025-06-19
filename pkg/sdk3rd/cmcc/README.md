移动云 Go SDK 文档: [https://ecloud.10086.cn/op-help-center/doc/article/53799](https://ecloud.10086.cn/op-help-center/doc/article/53799)

移动云 Go SDK 下载地址: [https://ecloud.10086.cn/api/query/developer/nexus/service/rest/repository/browse/go-sdk/gitlab.ecloud.com/ecloud/](https://ecloud.10086.cn/api/query/developer/nexus/service/rest/repository/browse/go-sdk/gitlab.ecloud.com/ecloud/)

---

将其引入本地目录的原因是:

1. 原始包必须通过移动云私有仓库获取, 为构建带来不便。
2. 原始包存在部分内容错误, 需要自行修改, 如：
   - 存在一些编译错误；
   - 返回错误的时候, 未返回错误信息；
   - 解析响应体错误。
