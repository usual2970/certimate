# Contributing to Certimate

Thank you for taking the time to improve Certimate! Below is a guide for submitting a PR (Pull Request) to the main Certimate repository.

- [Contributing to Certimate](#contributing-to-certimate)
  - [Prerequisites](#prerequisites)
  - [Making Changes in the Go Code](#making-changes-in-the-go-code)
  - [Making Changes in the Admin UI](#making-changes-in-the-admin-ui)

## Prerequisites

- Go 1.24+ (for Go code changes)
- Node 20+ (for Admin UI changes)

If you haven't done so already, you can fork the Certimate repository and clone your fork to work locally:

```bash
git clone https://github.com/your_username/certimate.git
```

> **Important:**
> It is recommended to create a new branch from `main` for each bug fix or feature. If you plan to submit multiple PRs, ensure the changes are in separate branches for easier review and eventual merge.
> Keep each PR focused on a single feature or fix.

## Making Changes in the Go Code

Once you have made changes to the Go code in Certimate, follow these steps to run the project:

1. Navigate to the root directory.

2. Start the service by running:

   ```bash
   go run main.go serve
   ```

This will start a web server at `http://localhost:8090` using the prebuilt Admin UI located in `ui/dist`.

**Before submitting a PR to the main repository, consider:**

- Format your source code by using [gofumpt](https://github.com/mvdan/gofumpt).

- Adding unit or integration tests for your changes. Certimate uses Go’s standard `testing` package. You can run tests using the following command (while in the root project directory):

  ```bash
  go test ./...
  ```

## Making Changes in the Admin UI

Certimate’s Admin UI is a single-page application (SPA) built using React and Vite.

To start the Admin UI:

1. Navigate to the `ui` project directory.

2. Install the necessary dependencies by running:

   ```bash
   npm install
   ```

3. Start the Vite development server:

   ```bash
   npm run dev
   ```

You can now access the running Admin UI at `http://localhost:5173` in your browser.

Since the Admin UI is a client-side application, you will also need to have the Certimate backend running. You can either manually run Certimate or use a prebuilt executable.

Any changes you make in the Admin UI will be automatically reflected in the browser without requiring a page reload.

After completing your changes, build the Admin UI so it can be embedded into the Go package:

```bash
npm run build
```

Once all steps are completed, you are ready to submit a PR to the main Certimate repository.
