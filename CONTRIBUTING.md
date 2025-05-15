# 向 Certimate 贡献代码

感谢你抽出时间来改进 Certimate！以下是向 Certimate 主仓库提交 PR（Pull Request）时的操作指南。

- [向 Certimate 贡献代码](#向-certimate-贡献代码)
  - [前提条件](#前提条件)
  - [修改 Go 代码](#修改-go-代码)
  - [修改管理页面 (Admin UI)](#修改管理页面-admin-ui)

## 前提条件

- Go 1.24+ (用于修改 Go 代码)
- Node 20+ (用于修改 UI)

如果还没有这样做，你可以 fork Certimate 的主仓库，并克隆到本地以便进行修改：

```bash
git clone https://github.com/your_username/certimate.git
```

> **重要提示：**
> 建议为每个 Bug 修复或新功能创建一个从 `main` 分支派生的新分支。如果你计划提交多个 PR，请保持不同的改动在独立分支中，以便更容易进行代码审查并最终合并。
> 保持一个 PR 只实现一个功能。

## 修改 Go 代码

假设你已经对 Certimate 的 Go 代码做了一些修改，现在你想要运行它：

1. 进入根目录
2. 运行以下命令启动服务：

   ```bash
   go run main.go serve
   ```

这将启动一个 Web 服务器，默认运行在 `http://localhost:8090`，并使用来自 `ui/dist` 的预构建管理页面。

**在向主仓库提交 PR 之前，建议你：**

- 使用 [gofumpt](https://github.com/mvdan/gofumpt) 对你的代码进行格式化。

- 为你的改动添加单元测试或集成测试（Certimate 使用 Go 的标准 `testing` 包）。你可以通过以下命令运行测试（在项目根目录下）：

  ```bash
  go test ./...
  ```

## 修改管理页面 (Admin UI)

Certimate 的管理页面是一个基于 React 和 Vite 构建的单页应用（SPA）。

要启动 Admin UI：

1. 进入 `ui` 项目目录。

2. 运行 `npm install` 安装依赖。

3. 启动 Vite 开发服务器：

   ```bash
   npm run dev
   ```

你可以通过浏览器访问 `http://localhost:5173` 来查看运行中的管理页面。

由于 Admin UI 只是一个客户端应用，运行时需要 Certimate 的后端服务作为支撑。你可以手动运行 Certimate，或者使用预构建的可执行文件。

所有对 Admin UI 的修改将会自动反映在浏览器中，无需手动刷新页面。

完成修改后，运行以下命令来构建 Admin UI，以便它能被嵌入到 Go 包中：

```bash
npm run build
```

完成所有步骤后，你可以将改动提交 PR 到 Certimate 主仓库。
