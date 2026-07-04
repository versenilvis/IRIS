# Iris commands

This directory contains all modular CLI command specifications and autocompletion definitions supported by **Iris**. Every command is defined as a `core.Spec` registered via package `init()` functions and grouped into category subdirectories.

The **[`core/`](./core)** subdirectory is the primary engine package. It implements the underlying command registry (`core.Registry`), data structures (`Spec`, `Arg`, `Flag`), dynamic generators, and autocompletion matching logic. The top-level **[`all.go`](./all.go)** file anonymously imports all category subpackages to trigger their initialization and register all available commands at startup.

## Overview

Currently, Iris natively supports **567** top-level CLI commands across **14** categories:

- [Cloud, Containers, Kubernetes, DevOps & Databases (`ops/`)](#ops): **118** commands
- [JavaScript, TypeScript, Frontend & Node.js Tools (`js/`)](#js): **82** commands
- [Python Ecosystem & Data Science Tools (`python/`)](#python): **19** commands
- [Rust Ecosystem & Modern CLI Tools (`rust/`)](#rust): **11** commands
- [Go Development & Project Tools (`golang/`)](#golang): **3** commands
- [Java, Kotlin & JVM Build Tools (`jvm/`)](#jvm): **14** commands
- [C/C++ Compilers & Build Systems (`cc/`)](#cc): **16** commands
- [Git Version Control & GitHub Tools (`git/`)](#git): **8** commands
- [System Package Managers (`pkginstaller/`)](#pkginstaller): **12** commands
- [Filesystem, Directory & Archive Utilities (`fs/`)](#fs): **30** commands
- [Editors, Pagers & File Viewers (`view/`)](#view): **27** commands
- [Text Processing, JSON & Stream Manipulation (`text/`)](#text): **28** commands
- [Task Runners & Build Automation (`runner/`)](#runner): **24** commands
- [System Administration, Network & Process Management (`sys/`)](#sys): **175** commands

---

<a id="ops"></a>
## Cloud, Containers, Kubernetes, DevOps & Databases (`ops/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`amplify`** | Environment | [`amplify.go`](./ops/amplify.go) |
| **`ampx`** | CLI for Amplify Gen 2 | [`ampx.go`](./ops/ampx.go) |
| **`ansible`** | Define and run a single Ansible task | [`ansible.go`](./ops/ansible.go) |
| **`ansible-config`** | View ansible configuration | [`ansible_config.go`](./ops/ansible_config.go) |
| **`ansible-doc`** | Displays information on modules installed in Ansible libraries | [`ansible_doc.go`](./ops/ansible_doc.go) |
| **`ansible-galaxy`** | The Galaxy API server URL | [`ansible_galaxy.go`](./ops/ansible_galaxy.go) |
| **`ansible-lint`** | Ansible static code analysis | [`ansible_lint.go`](./ops/ansible_lint.go) |
| **`ansible-playbook`** | Runs Ansible playbooks, executing the defined tasks on the targeted hosts | [`ansible_playbook.go`](./ops/ansible_playbook.go) |
| **`appwrite`** | Appwrite - Open-Source End-to-End Backend Server | [`appwrite.go`](./ops/appwrite.go) |
| **`arch`** | 32-bit intel | [`arch.go`](./ops/arch.go) |
| **`arduino-cli`** | Arduino CLI - build, compile, and upload Arduino sketches | [`arduino_cli.go`](./ops/arduino_cli.go) |
| **`argo`** | If True, Use the HTTP client. Defaults to the ARGO_HTTP1 environment variable | [`argo.go`](./ops/argo.go) |
| **`asdf`** | Plugin name | [`asdf.go`](./ops/asdf.go) |
| **`atlas`** | CLI tool to manage MongoDB Atlas | [`atlas.go`](./ops/atlas.go) |
| **`aws`** | Use a specific profile from your credential file | [`aws.go`](./ops/aws.go) |
| **`aws-vault`** | Add credentials to the secure keystore | [`aws_vault.go`](./ops/aws_vault.go) |
| **`bit`** | Bit documentation: https://bit.dev/docs | [`bit.go`](./ops/bit.go) |
| **`bosh`** | Deployment | [`bosh.go`](./ops/bosh.go) |
| **`capacitor`** | Add a native platform project to your app | [`capacitor.go`](./ops/capacitor.go) |
| **`cdk`** | AWS CDK CLI | [`cdk.go`](./ops/cdk.go) |
| **`cf`** | Cloudfoundry cli | [`cf.go`](./ops/cf.go) |
| **`checkov`** | Branch | [`checkov.go`](./ops/checkov.go) |
| **`circleci`** | CircleCI CLI | [`circleci.go`](./ops/circleci.go) |
| **`cloudflared`** | Specify the hostname of your application | [`cloudflared.go`](./ops/cloudflared.go) |
| **`coda`** | Coda CLI - interact with Coda docs and tables | [`coda.go`](./ops/coda.go) |
| **`command`** | Run an external command | [`command.go`](./ops/command.go) |
| **`copilot`** | Name of the application | [`copilot.go`](./ops/copilot.go) |
| **`cosign`** | Provides utilities for attaching artifacts to other artifacts in a registry | [`cosign.go`](./ops/cosign.go) |
| **`dapr`** | Distributed Application Runtime CLI | [`dapr.go`](./ops/dapr.go) |
| **`datree`** | Help for | [`datree.go`](./ops/datree.go) |
| **`deployctl`** | Command line tool for Deno Deploy | [`deployctl.go`](./ops/deployctl.go) |
| **`direnv`** | Help for direnv | [`direnv.go`](./ops/direnv.go) |
| **`docker`** | container engine | [`docker.go`](./ops/docker.go) |
| **`docker-compose`** | multi-container (legacy) | [`docker.go`](./ops/docker.go) |
| **`doctl`** | The official DigitalOcean command line interface (CLI) | [`doctl.go`](./ops/doctl.go) |
| **`doppler`** | The official CLI for Doppler Secret Operations Platform | [`doppler.go`](./ops/doppler.go) |
| **`eas`** | Log in with your Expo account | [`eas.go`](./ops/eas.go) |
| **`fastly`** | A CLI for interacting with the Fastly platform | [`fastly.go`](./ops/fastly.go) |
| **`firebase`** | ProjectAlias | [`firebase.go`](./ops/firebase.go) |
| **`flyctl`** | Command line tool for Fly.io services | [`flyctl.go`](./ops/flyctl.go) |
| **`fnm`** | Fast and simple Node.js version manager | [`fnm.go`](./ops/fnm.go) |
| **`gcloud`** | Manage Google Cloud Platform resources and developer workflow | [`gcloud.go`](./ops/gcloud.go) |
| **`gh`** | Current branch | [`gh.go`](./ops/gh.go) |
| **`gpg`** | Encryption and signing tool | [`gpg.go`](./ops/gpg.go) |
| **`hasura`** | .env filename to load ENV vars from | [`hasura.go`](./ops/hasura.go) |
| **`helm`** | The Helm package manager for Kubernetes | [`helm.go`](./ops/helm.go) |
| **`helmfile`** | Deploy helm charts | [`helmfile.go`](./ops/helmfile.go) |
| **`hugo`** | The world | [`hugo.go`](./ops/hugo.go) |
| **`k3d`** | Little helper to run k3s in Docker | [`k3d.go`](./ops/k3d.go) |
| **`k6`** | Create an archive | [`k6.go`](./ops/k6.go) |
| **`k9s`** | Kubernetes namespace | [`k9s.go`](./ops/k9s.go) |
| **`kind`** | Cluster | [`kind.go`](./ops/kind.go) |
| **`knex`** | SQL query builder for JavaScript | [`knex.go`](./ops/knex.go) |
| **`kubectl`** | kubernetes cli | [`kubectl.go`](./ops/kubectl.go) |
| **`kubectx`** | Switch between Kubernetes-contexts | [`kubectx.go`](./ops/kubectx.go) |
| **`kubens`** | Switch between Kubernetes-namespaces | [`kubens.go`](./ops/kubens.go) |
| **`limactl`** | Lima: Linux virtual machines, with a focus on running containers | [`limactl.go`](./ops/limactl.go) |
| **`locust`** | Show program | [`locust.go`](./ops/locust.go) |
| **`lpass`** | Command line interface for LastPass | [`lpass.go`](./ops/lpass.go) |
| **`minikube`** | Format to print stdout in | [`minikube.go`](./ops/minikube.go) |
| **`mongocli`** | CLI tool to manage your MongoDB Cloud | [`mongocli.go`](./ops/mongocli.go) |
| **`mongoimport`** | Import data from a JSON, CSV, or TSV file into a MongoDB instance | [`mongoimport.go`](./ops/mongoimport.go) |
| **`mongosh`** | Default Connection String; Equivalent to running mongosh without any commands | [`mongosh.go`](./ops/mongosh.go) |
| **`multipass`** | Displays help on commandline options | [`multipass.go`](./ops/multipass.go) |
| **`mysql`** | Mysql is a terminal-based front-end to MySQL | [`mysql.go`](./ops/mysql.go) |
| **`netlify`** | Print debugging information | [`netlify.go`](./ops/netlify.go) |
| **`newman`** | Newman is a command-line collection runner for Postman | [`newman.go`](./ops/newman.go) |
| **`nginx`** | Nginx (pronounced | [`nginx.go`](./ops/nginx.go) |
| **`ngrok`** | Path to log file, | [`ngrok.go`](./ops/ngrok.go) |
| **`nvm`** | Node version | [`nvm.go`](./ops/nvm.go) |
| **`oci`** | Oracle Cloud Infrastructure CLI | [`oci.go`](./ops/oci.go) |
| **`okteto`** | Context | [`okteto.go`](./ops/okteto.go) |
| **`op`** | Official 1Password CLI | [`op.go`](./ops/op.go) |
| **`opa`** | Open Policy Agent (OPA) | [`opa.go`](./ops/opa.go) |
| **`osqueryi`** | Your OS as a high-performance relational database | [`osqueryi.go`](./ops/osqueryi.go) |
| **`pass`** | Pass - stores, retrieves, generates, and synchronizes passwords securely | [`pass.go`](./ops/pass.go) |
| **`pg_dump`** | Dumps a database as a text file or to other formats | [`pg_dump.go`](./ops/pg_dump.go) |
| **`pgcli`** | Host address of the postgres database | [`pgcli.go`](./ops/pgcli.go) |
| **`pm2`** | Outputs the version number | [`pm2.go`](./ops/pm2.go) |
| **`pod`** | CocoaPods, the Cocoa library package manager | [`pod.go`](./ops/pod.go) |
| **`podman`** | Build an image using instructions from Containerfiles | [`podman.go`](./ops/podman.go) |
| **`pscale`** | The client ID for the PlanetScale CLI application | [`pscale.go`](./ops/pscale.go) |
| **`psql`** | Psql is a terminal-based front-end to PostgreSQL | [`psql.go`](./ops/psql.go) |
| **`pulumi`** | The name of the stack to operate on. Defaults to the current stack | [`pulumi.go`](./ops/pulumi.go) |
| **`qodana`** | Run Qodana as fast as possible, with minimum effort required | [`qodana.go`](./ops/qodana.go) |
| **`railway`** | CLI for managing Railway Apps | [`railway.go`](./ops/railway.go) |
| **`rbenv`** | List all available rbenv commands | [`rbenv.go`](./ops/rbenv.go) |
| **`robot`** | Tag | [`robot.go`](./ops/robot.go) |
| **`rsync`** | remote sync | [`ssh.go`](./ops/ssh.go) |
| **`scp`** | secure copy | [`ssh.go`](./ops/ssh.go) |
| **`serverless`** | AWS profile to use with the command | [`serverless.go`](./ops/serverless.go) |
| **`sfdx`** | Analyze (lint) Aura component code | [`sfdx.go`](./ops/sfdx.go) |
| **`sftp`** | OpenSSH secure file transfer | [`sftp.go`](./ops/sftp.go) |
| **`space`** | Deta Space CLI for mananging Deta Space projects | [`space.go`](./ops/space.go) |
| **`sqlite3`** | A command line interface for SQLite version 3 | [`sqlite3.go`](./ops/sqlite3.go) |
| **`src`** | Interact with Sourcegraph from the command line | [`src.go`](./ops/src.go) |
| **`ssh`** | secure shell | [`ssh.go`](./ops/ssh.go) |
| **`ssh-keygen`** | Generates, manages and converts authentication keys for ssh | [`ssh_keygen.go`](./ops/ssh_keygen.go) |
| **`stripe`** | Stripe CLI - build, test, and manage your Stripe integrations right from your terminal | [`stripe.go`](./ops/stripe.go) |
| **`supabase`** | Supabase CLI | [`supabase.go`](./ops/supabase.go) |
| **`surreal`** | Database authentication password to use when connecting [default: root] | [`surreal.go`](./ops/surreal.go) |
| **`tailscale`** | Connect to Tailscale, logging in if needed | [`tailscale.go`](./ops/tailscale.go) |
| **`terraform`** | Workspace | [`terraform.go`](./ops/terraform.go) |
| **`terragrunt`** | Workspace | [`terragrunt.go`](./ops/terragrunt.go) |
| **`tfenv`** | Version | [`tfenv.go`](./ops/tfenv.go) |
| **`tfsec`** | Terraform workspaces | [`tfsec.go`](./ops/tfsec.go) |
| **`tkn`** | CLI for tekton pipelines | [`tkn.go`](./ops/tkn.go) |
| **`trivy`** | Skip updating built-in policies [$TRIVY_SKIP_POLICY_UPDATE] | [`trivy.go`](./ops/trivy.go) |
| **`tsuru`** | Plan | [`tsuru.go`](./ops/tsuru.go) |
| **`vault`** | Display help | [`vault.go`](./ops/vault.go) |
| **`vela`** | Show the reference doc for component, trait or workflow types | [`vela.go`](./ops/vela.go) |
| **`vercel`** | CLI Interface for Vercel.com | [`vercel.go`](./ops/vercel.go) |
| **`volta`** | Enables verbose diagnostics | [`volta.go`](./ops/volta.go) |
| **`watson`** | A wonderful CLI to track your time | [`watson.go`](./ops/watson.go) |
| **`whois`** | Query a database for information about a domain registrant | [`whois.go`](./ops/whois.go) |
| **`wrangler`** | Path to configuration file [default: wrangler.toml] | [`wrangler.go`](./ops/wrangler.go) |
| **`xc`** | List tasks from an xc-compatible markdown file | [`xc.go`](./ops/xc.go) |
| **`xcodes`** | Manage the Xcode versions installed on your Mac | [`xcodes.go`](./ops/xcodes.go) |

<a id="js"></a>
## JavaScript, TypeScript, Frontend & Node.js Tools (`js/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`asar`** | A simple extensive tar-like archive format with indexing | [`asar.go`](./js/asar.go) |
| **`astro`** | Add an integration | [`astro.go`](./js/astro.go) |
| **`babel`** | A comma-separated list of preset names | [`babel.go`](./js/babel.go) |
| **`blitz`** | Show help for command | [`blitz.go`](./js/blitz.go) |
| **`browser-sync`** | Keep multiple browsers & devices in sync when building websites | [`browser_sync.go`](./js/browser_sync.go) |
| **`build-storybook`** | Storybook build CLI tools | [`build_storybook.go`](./js/build_storybook.go) |
| **`bun`** | bun js runtime | [`bun.go`](./js/bun.go) |
| **`bunx`** | execute package (bun x) | [`bun.go`](./js/bun.go) |
| **`cordova`** | Print out the version of your cordova-cli install | [`cordova.go`](./js/cordova.go) |
| **`create-completion-spec`** | Setup fig folder and create spec with the given name | [`create_completion_spec.go`](./js/create_completion_spec.go) |
| **`create-next-app`** | Output the version number | [`create_next_app.go`](./js/create_next_app.go) |
| **`create-nx-workspace`** | The name of the workspace | [`create_nx_workspace.go`](./js/create_nx_workspace.go) |
| **`create-react-app`** | Creates a new React project | [`create_react_app.go`](./js/create_react_app.go) |
| **`create-react-native-app`** | Creates a new React Native project | [`create_react_native_app.go`](./js/create_react_native_app.go) |
| **`create-redwood-app`** | Name of your Redwood project | [`create_redwood_app.go`](./js/create_redwood_app.go) |
| **`create-remix`** | Display help for command | [`create_remix.go`](./js/create_remix.go) |
| **`create-t3-app`** | The name of the application, as well as the name of the directory to create | [`create_t3_app.go`](./js/create_t3_app.go) |
| **`create-video`** | CLI used to create remotion video project | [`create_video.go`](./js/create_video.go) |
| **`create-vite`** | Create a new project powered by Vite | [`create_vite.go`](./js/create_vite.go) |
| **`create-web3-frontend`** | Quickly create a Next.js project with wagmi and TailwindCSS ready to go | [`create_web3_frontend.go`](./js/create_web3_frontend.go) |
| **`deno`** | A modern JavaScript and TypeScript runtime | [`deno.go`](./js/deno.go) |
| **`dotenv`** | Loads environment variables from .env | [`dotenv.go`](./js/dotenv.go) |
| **`electron`** | Build cross platform desktop apps with JavaScript, HTML and CSS | [`electron.go`](./js/electron.go) |
| **`elm`** | Fig spec for the Elm language cli | [`elm.go`](./js/elm.go) |
| **`elm-format`** | Format your code in the Elm idiomatic way | [`elm_format.go`](./js/elm_format.go) |
| **`elm-json`** | Deal with your elm.json | [`elm_json.go`](./js/elm_json.go) |
| **`elm-review`** | Prints a single JSON object | [`elm_review.go`](./js/elm_review.go) |
| **`esbuild`** | An extremely fast JavaScript bundler | [`esbuild.go`](./js/esbuild.go) |
| **`eslint`** | Pluggable JavaScript linter | [`eslint.go`](./js/eslint.go) |
| **`expo`** | Tools for creating, running, and deploying Universal Expo and React Native apps | [`expo.go`](./js/expo.go) |
| **`expo-cli`** | Tools for creating, running, and deploying Universal Expo and React Native apps | [`expo_cli.go`](./js/expo_cli.go) |
| **`ganache-cli`** | Fast Ethereum RPC client | [`ganache_cli.go`](./js/ganache_cli.go) |
| **`gatsby`** | Set host. Defaults to localhost | [`gatsby.go`](./js/gatsby.go) |
| **`hardhat`** | Ethereum development environment | [`hardhat.go`](./js/hardhat.go) |
| **`ionic`** | Target engine (e.g. browser, cordova) | [`ionic.go`](./js/ionic.go) |
| **`jest`** | A delightful JavaScript Testing Framework with a focus on simplicity | [`jest.go`](./js/jest.go) |
| **`lerna`** | Branch | [`lerna.go`](./js/lerna.go) |
| **`meteor`** | Run the meteor command-line tool | [`meteor.go`](./js/meteor.go) |
| **`ncu`** | Clear the default cache, or the cache file specified by --cacheFile | [`ncu.go`](./js/ncu.go) |
| **`nest`** | Report actions that would be taken without writing out results | [`nest.go`](./js/nest.go) |
| **`next`** | A port number on which to start the application | [`next.go`](./js/next.go) |
| **`ng`** | Project name | [`ng.go`](./js/ng.go) |
| **`node`** | Run the node interpreter | [`node.go`](./js/node.go) |
| **`npm`** | node packages | [`npm.go`](./js/npm.go) |
| **`npx`** | Execute binaries from npm packages | [`npx.go`](./js/npx.go) |
| **`nuxi`** | The directory of the target application | [`nuxi.go`](./js/nuxi.go) |
| **`nuxt`** | Launch the development server | [`nuxt.go`](./js/nuxt.go) |
| **`nx`** | All projects | [`nx.go`](./js/nx.go) |
| **`oxlint`** | All lints (except nursery) | [`oxlint.go`](./js/oxlint.go) |
| **`playwright`** | Display help for command | [`playwright.go`](./js/playwright.go) |
| **`pnpm`** | fast node packages | [`pnpm.go`](./js/pnpm.go) |
| **`pnpx`** | Execute binaries from npm packages | [`pnpx.go`](./js/pnpx.go) |
| **`prettier`** | Run Prettier from the command line | [`prettier.go`](./js/prettier.go) |
| **`quasar`** | Quasar Framework CLI | [`quasar.go`](./js/quasar.go) |
| **`react-native`** | Attempt to fix all diagnosed issues | [`react_native.go`](./js/react_native.go) |
| **`redwood`** | Script | [`redwood.go`](./js/redwood.go) |
| **`remix`** | Represent the directory of the Remix application | [`remix.go`](./js/remix.go) |
| **`remotion`** | Create videos programmatically in React | [`remotion.go`](./js/remotion.go) |
| **`rollup`** | Next-generation ES module bundler | [`rollup.go`](./js/rollup.go) |
| **`rome`** | Rome CLI | [`rome.go`](./js/rome.go) |
| **`rush`** | Projects | [`rush.go`](./js/rush.go) |
| **`sequelize`** | The environment to run the command in | [`sequelize.go`](./js/sequelize.go) |
| **`serve`** | Static file serving and directory listing | [`serve.go`](./js/serve.go) |
| **`shadcn-ui`** | Shadcn UI CLI | [`shadcn_ui.go`](./js/shadcn_ui.go) |
| **`start-storybook`** | Display usage information | [`start_storybook.go`](./js/start_storybook.go) |
| **`stencil`** | CLI to build Stencil projects and generate components | [`stencil.go`](./js/stencil.go) |
| **`swagger-typescript-api`** | Generate api via swagger scheme | [`swagger_typescript_api.go`](./js/swagger_typescript_api.go) |
| **`swc`** | Path to the file | [`swc.go`](./js/swc.go) |
| **`truffle`** | Execute build pipeline (if configuration present) | [`truffle.go`](./js/truffle.go) |
| **`ts-node`** | Run the TypeScript interpreter for Node.JS | [`ts_node.go`](./js/ts_node.go) |
| **`tsc`** | CLI tool for TypeScript compiler | [`tsc.go`](./js/tsc.go) |
| **`tsx`** | Run TypeScript file using tsx | [`tsx.go`](./js/tsx.go) |
| **`turbo`** | Print the version | [`turbo.go`](./js/turbo.go) |
| **`typeorm`** | Show help for command | [`typeorm.go`](./js/typeorm.go) |
| **`vite`** | Native ESM-powered web dev build tool | [`vite.go`](./js/vite.go) |
| **`vr`** | The npm-style script runner for Deno | [`vr.go`](./js/vr.go) |
| **`vsce`** | The Visual Studio Code Extension Manager | [`vsce.go`](./js/vsce.go) |
| **`vue`** | Vue cli tools | [`vue.go`](./js/vue.go) |
| **`watchman`** | A file watching service | [`watchman.go`](./js/watchman.go) |
| **`webpack`** | Run webpack (default command, can be omitted) | [`webpack.go`](./js/webpack.go) |
| **`yalc`** | Work with yarn/npm packages locally like a boss | [`yalc.go`](./js/yalc.go) |
| **`yarn`** | yarn package manager | [`yarn.go`](./js/yarn.go) |

<a id="python"></a>
## Python Ecosystem & Data Science Tools (`python/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`black`** | Version | [`black.go`](./python/black.go) |
| **`conda`** | Name of environment | [`conda.go`](./python/conda.go) |
| **`django-admin`** | Show this help message and exit | [`django_admin.go`](./python/django_admin.go) |
| **`googler`** | Google from the command-line | [`googler.go`](./python/googler.go) |
| **`jupyter`** | Set log level to logging.DEBUG (maximize logging output) | [`jupyter.go`](./python/jupyter.go) |
| **`mamba`** | Mamba is a fast, robust, and cross-platform package manager | [`mamba.go`](./python/mamba.go) |
| **`mypy`** | Mypy is a static type checker for Python | [`mypy.go`](./python/mypy.go) |
| **`pipenv`** | Python package manager | [`pipenv.go`](./python/pipenv.go) |
| **`pipx`** | Installed package | [`pipx.go`](./python/pipx.go) |
| **`poetry`** | python dependency manager | [`python.go`](./python/python.go) |
| **`pre-commit`** | Show help message and exit | [`pre_commit.go`](./python/pre_commit.go) |
| **`pyenv`** | Pyenv | [`pyenv.go`](./python/pyenv.go) |
| **`pytest`** | Control assertion debugging tools. | [`pytest.go`](./python/pytest.go) |
| **`ruff`** | Enable verbose logging | [`ruff.go`](./python/ruff.go) |
| **`sqlfluff`** | A dialect-flexible and configurable SQL linter | [`sqlfluff.go`](./python/sqlfluff.go) |
| **`sqlmesh`** | SQLMesh command line tool | [`sqlmesh.go`](./python/sqlmesh.go) |
| **`streamlit`** | Streamlit | [`streamlit.go`](./python/streamlit.go) |
| **`uv`** | fast python package manager | [`python.go`](./python/python.go) |
| **`youtube-dl`** | Clipboard | [`youtube_dl.go`](./python/youtube_dl.go) |

<a id="rust"></a>
## Rust Ecosystem & Modern CLI Tools (`rust/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`cargo`** | rust toolchain | [`cargo.go`](./rust/cargo.go) |
| **`dprint`** | Prints the help of the given subcommand(s) | [`dprint.go`](./rust/dprint.go) |
| **`pijul`** | Adds a path to the tree | [`pijul.go`](./rust/pijul.go) |
| **`rustc`** | Rust compiler | [`rustc.go`](./rust/rustc.go) |
| **`rustup`** | The Rust toolchain installer | [`rustup.go`](./rust/rustup.go) |
| **`taplo`** | Set color values for the output | [`taplo.go`](./rust/taplo.go) |
| **`tokei`** | Count your code, quickly | [`tokei.go`](./rust/tokei.go) |
| **`trunk`** | Run on all files instead of only changed files | [`trunk.go`](./rust/trunk.go) |
| **`wasm-bindgen`** | Generate bindings between WebAssembly and JavaScript | [`wasm_bindgen.go`](./rust/wasm_bindgen.go) |
| **`wasm-pack`** | Build an npm package | [`wasm_pack.go`](./rust/wasm_pack.go) |
| **`zellij`** | Change where zellij looks for the configuration file | [`zellij.go`](./rust/zellij.go) |

<a id="golang"></a>
## Go Development & Project Tools (`golang/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`go`** | tool for managing Go source code | [`go.go`](./golang/go.go) |
| **`goctl`** | A cli tool to generate go-zero code | [`goctl.go`](./golang/goctl.go) |
| **`goreleaser`** | Deliver Go binaries as fast and easily as possible | [`goreleaser.go`](./golang/goreleaser.go) |

<a id="jvm"></a>
## Java, Kotlin & JVM Build Tools (`jvm/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`clojure`** | An alias to refer to its function or a qualified function | [`clojure.go`](./jvm/clojure.go) |
| **`dart`** | The Dart file containing the main function | [`dart.go`](./jvm/dart.go) |
| **`flutter`** | Available emulators | [`flutter.go`](./jvm/flutter.go) |
| **`fvm`** | Print this usage information | [`fvm.go`](./jvm/fvm.go) |
| **`gradle`** | Log all warnings | [`gradle.go`](./jvm/gradle.go) |
| **`java`** | Java runtime | [`jvm.go`](./jvm/jvm.go) |
| **`javac`** | Java compiler | [`jvm.go`](./jvm/jvm.go) |
| **`jenv`** | Executable file | [`jenv.go`](./jvm/jenv.go) |
| **`jmeter`** | Apache JMeter - 100% Java Load Testing Tool | [`jmeter.go`](./jvm/jmeter.go) |
| **`kdoctor`** | Report a version of KDoctor | [`kdoctor.go`](./jvm/kdoctor.go) |
| **`keytool`** | Show help message | [`keytool.go`](./jvm/keytool.go) |
| **`kotlinc`** | Kotlin compiler | [`jvm.go`](./jvm/jvm.go) |
| **`mvn`** | Maven - a Java based project management and comprehension tool | [`mvn.go`](./jvm/mvn.go) |
| **`spring`** | Initialize a new project using Spring Initializr | [`spring.go`](./jvm/spring.go) |

<a id="cc"></a>
## C/C++ Compilers & Build Systems (`cc/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`bazel`** | Bazel target | [`bazel.go`](./cc/bazel.go) |
| **`c++`** | C++ compiler (alias) | [`cc.go`](./cc/cc.go) |
| **`cc`** | C compiler (alias) | [`cc.go`](./cc/cc.go) |
| **`clang`** | LLVM C compiler | [`cc.go`](./cc/cc.go) |
| **`clang++`** | LLVM C++ compiler | [`cc.go`](./cc/cc.go) |
| **`cmake`** | Command-line interface of the cross-platform buildsystem generator CMake | [`cmake.go`](./cc/cmake.go) |
| **`g++`** | GNU C++ compiler | [`cc.go`](./cc/cc.go) |
| **`gcc`** | GNU C compiler | [`cc.go`](./cc/cc.go) |
| **`premake`** | The premake5.lua file | [`premake.go`](./cc/premake.go) |
| **`swift`** | Show help information | [`swift.go`](./cc/swift.go) |
| **`typst`** | The Typst compiler | [`typst.go`](./cc/typst.go) |
| **`xcode-select`** | Active developer directory for Xcode tools | [`xcode_select.go`](./cc/xcode_select.go) |
| **`xcodebuild`** | Build Xcode projects | [`xcodebuild.go`](./cc/xcodebuild.go) |
| **`xcodeproj`** | Xcodeproj lets you create and modify Xcode projects | [`xcodeproj.go`](./cc/xcodeproj.go) |
| **`xcrun`** | SceneKit CLI utilities | [`xcrun.go`](./cc/xcrun.go) |
| **`zig`** | Enable or disable colored message | [`zig.go`](./cc/zig.go) |

<a id="git"></a>
## Git Version Control & GitHub Tools (`git/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`ghq`** | Clone/sync with a remote repository | [`ghq.go`](./git/ghq.go) |
| **`git`** | version control | [`git.go`](./git/git.go) |
| **`git-cliff`** | Increases the logging verbosity | [`git_cliff.go`](./git/git_cliff.go) |
| **`git-flow`** | Git extensions to provide high-level repository operations for Vincent Driessen | [`git_flow.go`](./git/git_flow.go) |
| **`git-profile`** | Use profile | [`git_profile.go`](./git/git_profile.go) |
| **`git-quick-stats`** | Show help for git-quick-stats | [`git_quick_stats.go`](./git/git_quick_stats.go) |
| **`github`** | Open a git repository in GitHub Desktop | [`github.go`](./git/github.go) |
| **`svn`** | Specify a username ARG | [`svn.go`](./git/svn.go) |

<a id="pkginstaller"></a>
## System Package Managers (`pkginstaller/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`apt`** | Debian/Ubuntu package manager | [`pkgmgr.go`](./pkginstaller/pkgmgr.go) |
| **`apt-get`** | Debian/Ubuntu package manager (low-level) | [`pkgmgr.go`](./pkginstaller/pkgmgr.go) |
| **`brew`** | Homebrew package manager | [`pkgmgr.go`](./pkginstaller/pkgmgr.go) |
| **`dnf`** | Fedora/RHEL package manager | [`pkgmgr.go`](./pkginstaller/pkgmgr.go) |
| **`dpkg`** | Debian package management system | [`dpkg.go`](./pkginstaller/dpkg.go) |
| **`flatpak`** | flatpak package manager | [`pkgmgr.go`](./pkginstaller/pkgmgr.go) |
| **`pacman`** | Arch package manager | [`pkgmgr.go`](./pkginstaller/pkgmgr.go) |
| **`paru`** | AUR helper (feature-rich) | [`pkgmgr.go`](./pkginstaller/pkgmgr.go) |
| **`pkgutil`** | Query and manipulate for macOS Installer packages and receipts | [`pkgutil.go`](./pkginstaller/pkgutil.go) |
| **`snap`** | snap package manager | [`pkgmgr.go`](./pkginstaller/pkgmgr.go) |
| **`yay`** | AUR helper (pacman wrapper) | [`pkgmgr.go`](./pkginstaller/pkgmgr.go) |
| **`yum`** | RHEL/CentOS package manager (legacy) | [`pkgmgr.go`](./pkginstaller/pkgmgr.go) |

<a id="fs"></a>
## Filesystem, Directory & Archive Utilities (`fs/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`broot`** | Show the last modified date of files and directories | [`broot.go`](./fs/broot.go) |
| **`cd`** | change directory | [`cd.go`](./fs/cd.go) |
| **`chmod`** | change file permissions | [`chmod.go`](./fs/chmod.go) |
| **`chown`** | change file owner | [`chown.go`](./fs/chown.go) |
| **`cp`** | copy files and directories | [`cp.go`](./fs/cp.go) |
| **`df`** | Display free disk space | [`df.go`](./fs/df.go) |
| **`dust`** | Like du but more intuitive | [`dust.go`](./fs/dust.go) |
| **`exa`** | A modern replacement for ls | [`exa.go`](./fs/exa.go) |
| **`eza`** | A modern replacement for ls | [`eza.go`](./fs/eza.go) |
| **`find`** | Walk a file hierarchy | [`find.go`](./fs/find.go) |
| **`fold`** | Fold long lines for finite width output device | [`fold.go`](./fs/fold.go) |
| **`ln`** | create links | [`ln.go`](./fs/ln.go) |
| **`ls`** | list directory contents | [`ls.go`](./fs/ls.go) |
| **`lsd`** | An ls command with a lot of pretty colors and some other stuff | [`lsd.go`](./fs/lsd.go) |
| **`mkdir`** | make directories | [`mkdir.go`](./fs/mkdir.go) |
| **`mv`** | move (rename) files | [`mv.go`](./fs/mv.go) |
| **`paper`** | The Paper CLI | [`paper.go`](./fs/paper.go) |
| **`rclone`** | Only list directories | [`rclone.go`](./fs/rclone.go) |
| **`readlink`** | Display file status | [`readlink.go`](./fs/readlink.go) |
| **`rm`** | remove files or directories | [`rm.go`](./fs/rm.go) |
| **`rmdir`** | Remove directories | [`rmdir.go`](./fs/rmdir.go) |
| **`stow`** | Manage farms of symbolic links | [`stow.go`](./fs/stow.go) |
| **`tar`** | Use archive file or device ARCHIVE | [`tar.go`](./fs/tar.go) |
| **`touch`** | create or update file timestamp | [`touch.go`](./fs/touch.go) |
| **`trash`** | Trash, move files/folders to the trash | [`trash.go`](./fs/trash.go) |
| **`tree`** | Display directories as trees (with optional color/HTML output) | [`tree.go`](./fs/tree.go) |
| **`unzip`** | Extract compressed files in a ZIP archive | [`unzip.go`](./fs/unzip.go) |
| **`z`** | jump to directory | [`zoxide.go`](./fs/zoxide.go) |
| **`zi`** | jump to directory interactively | [`zoxide.go`](./fs/zoxide.go) |
| **`zip`** | Package and compress (archive) files into zip file | [`zip.go`](./fs/zip.go) |

<a id="view"></a>
## Editors, Pagers & File Viewers (`view/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`bat`** | A cat(1) clone with syntax highlighting and Git integration | [`bat.go`](./view/bat.go) |
| **`cat`** | concatenate and print | [`cat.go`](./view/cat.go) |
| **`code`** | Read from stdin (e.g. | [`code.go`](./view/code.go) |
| **`cot`** | Command-line utility for CotEditor | [`cot.go`](./view/cot.go) |
| **`du`** | estimate file space usage | [`du.go`](./view/du.go) |
| **`emacs`** | An extensible, customizable, free/libre text editor - and more | [`emacs.go`](./view/emacs.go) |
| **`file`** | determine file type | [`file.go`](./view/file.go) |
| **`glow`** | Render markdown on the CLI, with pizzazz! | [`glow.go`](./view/glow.go) |
| **`head`** | output first lines of file | [`head.go`](./view/head.go) |
| **`idea`** | IntelliJ IDEA CLI | [`idea.go`](./view/idea.go) |
| **`less`** | view file contents (scrollable) | [`less.go`](./view/less.go) |
| **`lvim`** | Hyperextensible Vim-based text editor | [`lvim.go`](./view/lvim.go) |
| **`micro`** | True/false | [`micro.go`](./view/micro.go) |
| **`more`** | Opposite of less | [`more.go`](./view/more.go) |
| **`nano`** | Nano | [`nano.go`](./view/nano.go) |
| **`nvim`** | Hyperextensible Vim-based text editor | [`nvim.go`](./view/nvim.go) |
| **`rich`** | Defined by terminal, appearance may differ | [`rich.go`](./view/rich.go) |
| **`stat`** | display file status | [`stat.go`](./view/stat.go) |
| **`subl`** | Sublime Text | [`subl.go`](./view/subl.go) |
| **`tail`** | output last lines of file | [`tail.go`](./view/tail.go) |
| **`vi`** | Print help message for vi and exit | [`vi.go`](./view/vi.go) |
| **`vim`** | Vi IMproved, a programmer | [`vim.go`](./view/vim.go) |
| **`vimr`** | VimR - Neovim GUI for macOS in Swift | [`vimr.go`](./view/vimr.go) |
| **`wc`** | word, line, character count | [`wc.go`](./view/wc.go) |
| **`xed`** | Xcode text editor invocation tool | [`xed.go`](./view/xed.go) |
| **`xxd`** | Make a hexdump or do the reverse | [`xxd.go`](./view/xxd.go) |
| **`zed`** | A lightning-fast, collaborative code editor written in Rust | [`zed.go`](./view/zed.go) |

<a id="text"></a>
## Text Processing, JSON & Stream Manipulation (`text/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`awk`** | pattern-directed scanning | [`textproc.go`](./text/textproc.go) |
| **`cut`** | extract columns from lines | [`textproc.go`](./text/textproc.go) |
| **`diff`** | Compare files line by line | [`diff.go`](./text/diff.go) |
| **`dos2unix`** | DOS to Unix file format converter | [`dos2unix.go`](./text/dos2unix.go) |
| **`egrep`** | grep with extended regex | [`grep.go`](./text/grep.go) |
| **`fd`** | fast find alternative | [`rg.go`](./text/rg.go) |
| **`find`** | search for files | [`find.go`](./text/find.go) |
| **`gawk`** | GNU awk | [`textproc.go`](./text/textproc.go) |
| **`grep`** | search text in files | [`grep.go`](./text/grep.go) |
| **`iconv`** | Character set conversion | [`iconv.go`](./text/iconv.go) |
| **`jq`** | Output the jq version and exit with zero | [`jq.go`](./text/jq.go) |
| **`pandoc`** | A universal document converter | [`pandoc.go`](./text/pandoc.go) |
| **`rg`** | ripgrep (fast search) | [`rg.go`](./text/rg.go) |
| **`sed`** | stream editor | [`textproc.go`](./text/textproc.go) |
| **`seq`** | Print sequences of numbers. (Defaults to increments of 1) | [`seq.go`](./text/seq.go) |
| **`sha1sum`** | Print or check SHA1 (160-bit) checksums | [`sha1sum.go`](./text/sha1sum.go) |
| **`shasum`** | Print or Check SHA Checksums | [`shasum.go`](./text/shasum.go) |
| **`shred`** | Overwrite a file to hide its contents, and optionally delete it | [`shred.go`](./text/shred.go) |
| **`sort`** | sort lines of text | [`textproc.go`](./text/textproc.go) |
| **`split`** | Use suffix_length letters to form the suffix of the file name | [`split.go`](./text/split.go) |
| **`tee`** | read stdin, write to stdout and files | [`textproc.go`](./text/textproc.go) |
| **`tr`** | translate or delete characters | [`textproc.go`](./text/textproc.go) |
| **`truncate`** | Shrink or extend the size of a file to the specified size | [`truncate.go`](./text/truncate.go) |
| **`typos`** | Source code spelling correction | [`typos.go`](./text/typos.go) |
| **`uniq`** | filter adjacent duplicate lines | [`textproc.go`](./text/textproc.go) |
| **`unix2dos`** | Unix to DOS text file format convertor | [`unix2dos.go`](./text/unix2dos.go) |
| **`vale`** | A syntax-aware linter for prose built with speed and extensibility in mind | [`vale.go`](./text/vale.go) |
| **`xargs`** | build and run commands from stdin | [`textproc.go`](./text/textproc.go) |

<a id="runner"></a>
## Task Runners & Build Automation (`runner/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`ant`** | Apache Ant - Java library and command-line build tool | [`ant.go`](./runner/ant.go) |
| **`composer`** | Composer Command | [`composer.go`](./runner/composer.go) |
| **`dbt`** | CLI for dbt - Data Build Tool | [`dbt.go`](./runner/dbt.go) |
| **`drush`** | Drush is a command line shell and Unix scripting interface for Drupal | [`drush.go`](./runner/drush.go) |
| **`elixir`** | Elixir Language | [`elixir.go`](./runner/elixir.go) |
| **`gem`** | Use HTTP proxy for remote operations | [`gem.go`](./runner/gem.go) |
| **`hexo`** | Draft for | [`hexo.go`](./runner/hexo.go) |
| **`just`** | command runner | [`justfile.go`](./runner/justfile.go) |
| **`laravel`** | The output format (txt, xml, json, or md) | [`laravel.go`](./runner/laravel.go) |
| **`magento`** | Open-source E-commerce | [`magento.go`](./runner/magento.go) |
| **`make`** | build automation | [`makefile.go`](./runner/makefile.go) |
| **`mix`** | Build tool for Elixir | [`mix.go`](./runner/mix.go) |
| **`php`** | Run the PHP interpreter | [`php.go`](./runner/php.go) |
| **`phpunit`** | Generate code coverage report in Clover XML format, | [`phpunit.go`](./runner/phpunit.go) |
| **`phpunit-watcher`** | Automatically rerun PHPUnit tests when source code changes | [`phpunit_watcher.go`](./runner/phpunit_watcher.go) |
| **`rails`** | Create a new rails application | [`rails.go`](./runner/rails.go) |
| **`rake`** | A ruby build program with capabilities similar to make | [`rake.go`](./runner/rake.go) |
| **`rubocop`** | Run only lint cops | [`rubocop.go`](./runner/rubocop.go) |
| **`ruby`** | Interpreted object-oriented scripting language | [`ruby.go`](./runner/ruby.go) |
| **`rvm`** | Show version of rvm | [`rvm.go`](./runner/rvm.go) |
| **`sidekiq`** | Background job framework for Ruby | [`sidekiq.go`](./runner/sidekiq.go) |
| **`symfony`** | Symfony Binary | [`symfony.go`](./runner/symfony.go) |
| **`valet`** | Do not output any message | [`valet.go`](./runner/valet.go) |
| **`vapor`** | Vapor Toolbox (Server-side Swift web framework) | [`vapor.go`](./runner/vapor.go) |

<a id="sys"></a>
## System Administration, Network & Process Management (`sys/`)

| Command | Description | Source File |
| :--- | :--- | :--- |
| **`adb`** | Forward-lock the app | [`adb.go`](./sys/adb.go) |
| **`ag`** | Recursively search for PATTERN in PATH. Like grep or ack, but faster | [`ag.go`](./sys/ag.go) |
| **`airflow`** | Subcommand | [`airflow.go`](./sys/airflow.go) |
| **`aliases`** | Prints help information | [`aliases.go`](./sys/aliases.go) |
| **`asciinema`** | Terminal session recorder | [`asciinema.go`](./sys/asciinema.go) |
| **`asr`** | Can be a disk image, /dev entry, or volume mountpoint | [`asr.go`](./sys/asr.go) |
| **`atuin`** | Magical shell history | [`atuin.go`](./sys/atuin.go) |
| **`basename`** | Return filename portion of pathname | [`basename.go`](./sys/basename.go) |
| **`bc`** | An arbitrary precision calculator language | [`bc.go`](./sys/bc.go) |
| **`btop`** | Beautifuler htop (interactive process viewer) | [`btop.go`](./sys/btop.go) |
| **`bundle`** | Gem | [`bundle.go`](./sys/bundle.go) |
| **`cal`** | Displays a calendar and the date of Easter | [`cal.go`](./sys/cal.go) |
| **`cci`** | CumulusCI command line interface | [`cci.go`](./sys/cci.go) |
| **`cdk8s`** | CDK for K8s | [`cdk8s.go`](./sys/cdk8s.go) |
| **`chezmoi`** | Attribute modifier | [`chezmoi.go`](./sys/chezmoi.go) |
| **`chsh`** | Change your login shell | [`chsh.go`](./sys/chsh.go) |
| **`codesign`** | Create and manipulate code signatures | [`codesign.go`](./sys/codesign.go) |
| **`croc`** | Send file(s), or folder | [`croc.go`](./sys/croc.go) |
| **`crontab`** | Maintain crontab file for individual users | [`crontab.go`](./sys/crontab.go) |
| **`curl`** | transfer data via URL | [`network.go`](./sys/network.go) |
| **`date`** | Display or set date and time | [`date.go`](./sys/date.go) |
| **`dateseq`** | Print help and exit | [`dateseq.go`](./sys/dateseq.go) |
| **`dcli`** | Display help for command | [`dcli.go`](./sys/dcli.go) |
| **`dd`** | The same as | [`dd.go`](./sys/dd.go) |
| **`ddev`** | DDEV-Local local development environment | [`ddev.go`](./sys/ddev.go) |
| **`defaults`** | Global domain | [`defaults.go`](./sys/defaults.go) |
| **`degit`** | Straightforward project scaffolding | [`degit.go`](./sys/degit.go) |
| **`deta`** | Runtime | [`deta.go`](./sys/deta.go) |
| **`dig`** | DNS lookup | [`network.go`](./sys/network.go) |
| **`dirname`** | Return directory portion of pathname | [`dirname.go`](./sys/dirname.go) |
| **`do-release-upgrade`** | Upgrade Ubuntu to latest release | [`do_release_upgrade.go`](./sys/do_release_upgrade.go) |
| **`dog`** | Human-readable host names, nameservers, types, or classes | [`dog.go`](./sys/dog.go) |
| **`dotnet`** | The dotnet cli | [`dotnet.go`](./sys/dotnet.go) |
| **`dscacheutil`** | Utility for managing the Directory Service cache | [`dscacheutil.go`](./sys/dscacheutil.go) |
| **`dscl`** | Prompt for password | [`dscl.go`](./sys/dscl.go) |
| **`dtm`** | Plugin | [`dtm.go`](./sys/dtm.go) |
| **`echo`** | Environment Variable | [`echo.go`](./sys/echo.go) |
| **`eleventy`** | Eleventy is a simpler static site generator | [`eleventy.go`](./sys/eleventy.go) |
| **`env`** | print environment | [`env.go`](./sys/env.go) |
| **`exec`** | Replace the current shell with a program | [`exec.go`](./sys/exec.go) |
| **`export`** | set environment variable | [`env.go`](./sys/env.go) |
| **`fastlane`** | Helps you with your initial fastlane setup | [`fastlane.go`](./sys/fastlane.go) |
| **`fdisk`** | Manipulate disk partition table | [`fdisk.go`](./sys/fdisk.go) |
| **`ffmpeg`** | Play, record, convert, and stream audio and video | [`ffmpeg.go`](./sys/ffmpeg.go) |
| **`firefox`** | Free open-source web browser developer by Mozilla | [`firefox.go`](./sys/firefox.go) |
| **`fisher`** | [Prompt] - 🌊 The ultimate Fish prompt | [`fisher.go`](./sys/fisher.go) |
| **`fmt`** | Simple text formatter | [`fmt.go`](./sys/fmt.go) |
| **`forc`** | Fuel Orchestrator | [`forc.go`](./sys/forc.go) |
| **`forge`** | A command line interface for managing Atlassian-hosted apps | [`forge.go`](./sys/forge.go) |
| **`fzf`** | A general-purpose command-line fuzzy finder | [`fzf.go`](./sys/fzf.go) |
| **`fzf-tmux`** | Opens a fuzzy finder in a tmux pane | [`fzf_tmux.go`](./sys/fzf_tmux.go) |
| **`gltfjsx`** | GLTF to JSX converter | [`gltfjsx.go`](./sys/gltfjsx.go) |
| **`goto`** | Goto | [`goto.go`](./sys/goto.go) |
| **`gum`** | Background Color | [`gum.go`](./sys/gum.go) |
| **`herd`** | Display this application version | [`herd.go`](./sys/herd.go) |
| **`hop`** | Interact with Hop in your terminal | [`hop.go`](./sys/hop.go) |
| **`hostname`** | Set or print name of current host system | [`hostname.go`](./sys/hostname.go) |
| **`htop`** | Improved top (interactive process viewer) | [`htop.go`](./sys/htop.go) |
| **`http`** | HTTPie: command-line HTTP client for the API era | [`http.go`](./sys/http.go) |
| **`hyper`** | Hyper is an Electron-based terminal | [`hyper.go`](./sys/hyper.go) |
| **`hyperfine`** | A command-line benchmarking tool | [`hyperfine.go`](./sys/hyperfine.go) |
| **`ibus`** | Set or get engine | [`ibus.go`](./sys/ibus.go) |
| **`id`** | Display the full name of the user | [`id.go`](./sys/id.go) |
| **`ifconfig`** | configure network interface | [`network.go`](./sys/network.go) |
| **`ignite-cli`** | Output usage information | [`ignite_cli.go`](./sys/ignite_cli.go) |
| **`install`** | Use suffix as the backup suffix if -b is given | [`install.go`](./sys/install.go) |
| **`ip`** | show/manage network | [`network.go`](./sys/network.go) |
| **`join`** | The join utility performs an | [`join.go`](./sys/join.go) |
| **`julia`** | The Julia Programming Language | [`julia.go`](./sys/julia.go) |
| **`kafkactl`** | Command-line interface for Apache Kafka | [`kafkactl.go`](./sys/kafkactl.go) |
| **`kamal`** | Skip image build and push | [`kamal.go`](./sys/kamal.go) |
| **`kill`** | send signal to process | [`ps.go`](./sys/ps.go) |
| **`killall`** | kill by process name | [`ps.go`](./sys/ps.go) |
| **`kitty`** | A cat like utility to display images in the terminal | [`kitty.go`](./sys/kitty.go) |
| **`klist`** | Credential cache to list | [`klist.go`](./sys/klist.go) |
| **`kool`** | Script | [`kool.go`](./sys/kool.go) |
| **`launchctl`** | Service or domain target | [`launchctl.go`](./sys/launchctl.go) |
| **`leaf`** | Create and interact with your leaf projects | [`leaf.go`](./sys/leaf.go) |
| **`lima`** | Lima is an alias for | [`lima.go`](./sys/lima.go) |
| **`login`** | Begin session on the system | [`login.go`](./sys/login.go) |
| **`lsblk`** | List block devices | [`lsblk.go`](./sys/lsblk.go) |
| **`lsof`** | List open files | [`lsof.go`](./sys/lsof.go) |
| **`man`** | Format and display manual pages | [`man.go`](./sys/man.go) |
| **`meroxa`** | The Meroxa CLI | [`meroxa.go`](./sys/meroxa.go) |
| **`mkdocs`** | Project documentation with Markdown | [`mkdocs.go`](./sys/mkdocs.go) |
| **`mkfifo`** | Make FIFOs (first-in, first-out) | [`mkfifo.go`](./sys/mkfifo.go) |
| **`mkinitcpio`** | Create an initial ramdisk environment | [`mkinitcpio.go`](./sys/mkinitcpio.go) |
| **`mknod`** | Create device special file | [`mknod.go`](./sys/mknod.go) |
| **`mosh`** | Address of remote machine to log into | [`mosh.go`](./sys/mosh.go) |
| **`mount`** | Mount disks and manage subtrees | [`mount.go`](./sys/mount.go) |
| **`nc`** | netcat - TCP/UDP tool | [`network.go`](./sys/network.go) |
| **`ncal`** | Displays a calendar and the date of Easter | [`ncal.go`](./sys/ncal.go) |
| **`neofetch`** | The most complete system information CLI tool | [`neofetch.go`](./sys/neofetch.go) |
| **`netstat`** | network statistics | [`network.go`](./sys/network.go) |
| **`networkQuality`** | Measure the different aspects of network quality | [`networkquality.go`](./sys/networkquality.go) |
| **`networksetup`** | Configuration tool for network settings in macOS | [`networksetup.go`](./sys/networksetup.go) |
| **`nextflow`** | Session ID | [`nextflow.go`](./sys/nextflow.go) |
| **`nhost`** | Nhost | [`nhost.go`](./sys/nhost.go) |
| **`nmap`** | Network exploration tool and security / port scanner | [`nmap.go`](./sys/nmap.go) |
| **`nrm`** | Use the right package manage - remove | [`nrm.go`](./sys/nrm.go) |
| **`ns`** | Forces rebuilding the native application | [`ns.go`](./sys/ns.go) |
| **`nslookup`** | query DNS | [`network.go`](./sys/network.go) |
| **`nylas`** | A command line interface for Nylas | [`nylas.go`](./sys/nylas.go) |
| **`oh-my-posh`** | The config file to use | [`oh_my_posh.go`](./sys/oh_my_posh.go) |
| **`okta`** | The Okta CLI is the easiest way to get started with Okta! | [`okta.go`](./sys/okta.go) |
| **`ollama`** | A command-line tool for managing and deploying machine learning models | [`ollama.go`](./sys/ollama.go) |
| **`omz`** | Oh My Zsh | [`omz.go`](./sys/omz.go) |
| **`pac`** | 7 | [`pac.go`](./sys/pac.go) |
| **`passwd`** | Modify a user | [`passwd.go`](./sys/passwd.go) |
| **`pathchk`** | Check pathnames for POSIX portability | [`pathchk.go`](./sys/pathchk.go) |
| **`pdfunite`** | Combine multiple pdfs | [`pdfunite.go`](./sys/pdfunite.go) |
| **`pgrep`** | find process by pattern | [`ps.go`](./sys/ps.go) |
| **`ping`** | test network connectivity | [`network.go`](./sys/network.go) |
| **`pkg-config`** | Return metainformation about installed libraries | [`pkg_config.go`](./sys/pkg_config.go) |
| **`pkill`** | kill by pattern | [`ps.go`](./sys/ps.go) |
| **`pmset`** | Display sleep timer (value in minutes, or 0 to disable) | [`pmset.go`](./sys/pmset.go) |
| **`pocketbase`** | PocketBase CLI | [`pocketbase.go`](./sys/pocketbase.go) |
| **`printenv`** | print environment variables | [`env.go`](./sys/env.go) |
| **`prisma`** | Display this help message | [`prisma.go`](./sys/prisma.go) |
| **`pro`** | Manage Ubuntu Pro services from Canonical | [`pro.go`](./sys/pro.go) |
| **`pry`** | Interactive Ruby | [`pry.go`](./sys/pry.go) |
| **`ps`** | report processes | [`ps.go`](./sys/ps.go) |
| **`publish`** | Set up a new website in the current folder | [`publish.go`](./sys/publish.go) |
| **`pwd`** | Return working directory name | [`pwd.go`](./sys/pwd.go) |
| **`rancher`** | Output format: | [`rancher.go`](./sys/rancher.go) |
| **`repeat`** | Interpret the result as a number and repeat the commands this many times | [`repeat.go`](./sys/repeat.go) |
| **`rscript`** | Scripting Front-End for R | [`rscript.go`](./sys/rscript.go) |
| **`sam`** | Host of locally emulated Lambda container | [`sam.go`](./sys/sam.go) |
| **`sanity`** | Displays help information about Sanity | [`sanity.go`](./sys/sanity.go) |
| **`screen`** | Screen manager with VT100/ANSI terminal emulation | [`screen.go`](./sys/screen.go) |
| **`shell-config`** | Display help for command | [`shell_config.go`](./sys/shell_config.go) |
| **`shortcuts`** | Run a shortcut | [`shortcuts.go`](./sys/shortcuts.go) |
| **`simctl`** | Add photos, live photos, videos, or contacts to the library of a device | [`simctl.go`](./sys/simctl.go) |
| **`source`** | Source files in shell | [`source.go`](./sys/source.go) |
| **`speedtest-cli`** | Command line interface for testing internet bandwidth using speedtest.net | [`speedtest_cli.go`](./sys/speedtest_cli.go) |
| **`spotify`** | CLI to use Spotify from the terminal | [`spotify.go`](./sys/spotify.go) |
| **`ss`** | socket statistics | [`network.go`](./sys/network.go) |
| **`st2`** | Show this help and exit | [`st2.go`](./sys/st2.go) |
| **`stack`** | The Haskell Tool Stack | [`stack.go`](./sys/stack.go) |
| **`starkli`** | Starkli, a ⚡ blazing ⚡ fast ⚡ CLI tool for Starknet powered by 🦀 starknet-rs 🦀 | [`starkli.go`](./sys/starkli.go) |
| **`su`** | (no letter) The same as -l | [`su.go`](./sys/su.go) |
| **`sudo`** | Execute a command as the superuser or another user | [`sudo.go`](./sys/sudo.go) |
| **`sysctl`** | Variable name | [`sysctl.go`](./sys/sysctl.go) |
| **`systemctl`** | Control the systemd system and service manager | [`systemctl.go`](./sys/systemctl.go) |
| **`tac`** | Concatenate and print files in reverse | [`tac.go`](./sys/tac.go) |
| **`tailcall`** | TailCall CLI for managing and optimizing GraphQL configurations | [`tailcall.go`](./sys/tailcall.go) |
| **`tailwindcss`** | Display usage information | [`tailwindcss.go`](./sys/tailwindcss.go) |
| **`time`** | Time how long a command takes! | [`time.go`](./sys/time.go) |
| **`tldr`** | Tldr page | [`tldr.go`](./sys/tldr.go) |
| **`tmux`** | Format output | [`tmux.go`](./sys/tmux.go) |
| **`tmuxinator`** | Project | [`tmuxinator.go`](./sys/tmuxinator.go) |
| **`top`** | Display Linux tasks | [`top.go`](./sys/top.go) |
| **`traceroute`** | Print the route packets take to network host | [`traceroute.go`](./sys/traceroute.go) |
| **`trap`** | Prints all defined signal handlers | [`trap.go`](./sys/trap.go) |
| **`trex`** | trex script | [`trex.go`](./sys/trex.go) |
| **`tsh`** | Remote host login | [`tsh.go`](./sys/tsh.go) |
| **`tuist`** | Build the project in the current directory | [`tuist.go`](./sys/tuist.go) |
| **`twilio`** | Level of logging messages | [`twilio.go`](./sys/twilio.go) |
| **`uname`** | Print operating system name | [`uname.go`](./sys/uname.go) |
| **`unset`** | unset variable | [`env.go`](./sys/env.go) |
| **`visudo`** | Checking existing sudoers file for syntax errors | [`visudo.go`](./sys/visudo.go) |
| **`vultr-cli`** | Bare Metal ID | [`vultr_cli.go`](./sys/vultr_cli.go) |
| **`wezterm`** | Wez | [`wezterm.go`](./sys/wezterm.go) |
| **`wget`** | non-interactive downloader | [`network.go`](./sys/network.go) |
| **`where`** | For each name, indicate how it should be interpreted | [`where.go`](./sys/where.go) |
| **`whereis`** | Locate the binary, source, and manual page files for a command | [`whereis.go`](./sys/whereis.go) |
| **`which`** | Executable file | [`which.go`](./sys/which.go) |
| **`who`** | Display who is logged in | [`who.go`](./sys/who.go) |
| **`wing`** | Runs a Wing executable in the Wing Console | [`wing.go`](./sys/wing.go) |
| **`wp`** | Path to the WordPress files | [`wp.go`](./sys/wp.go) |
| **`wrk`** | Wrk - a HTTP benchmarking tool | [`wrk.go`](./sys/wrk.go) |
| **`wscat`** | Communicate over websocket | [`wscat.go`](./sys/wscat.go) |
| **`yank`** | Yank terminal output to clipboard | [`yank.go`](./sys/yank.go) |
| **`ykman`** | Configure your YubiKey via the command line | [`ykman.go`](./sys/ykman.go) |
| **`zapier`** | Change the way structured data is presented. If | [`zapier.go`](./sys/zapier.go) |

