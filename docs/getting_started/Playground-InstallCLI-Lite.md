# Install Nexus CLI
[[Prev]](Playground-Lite.md) [[Exit]](../../README.md)  [[Next]](Playground-InstallRuntime-Lite.md)

![InstallCLI](../images/Playground-2-install-cli.png)

### 1. Clone Nexus Repo
```
git clone git@github.com:intel-sandbox/applications.development.framework.nexus.git
cd applications.development.framework.nexus/
export NEXUS_REPO_DIR=${PWD}
```

### 2. Build & Install Nexus CLI
#### For Linux
```
make cli.build.linux
```
For MacOS
```
make cli.build.darwin
```

### 3. Install Nexus CLI
#### For Linux
```
sudo make cli.install.linux
```
For MacOS
```
sudo make cli.install.darwin
```

[[Prev]](Playground-Lite.md) [[Exit]](../../README.md)  [[Next]](Playground-InstallRuntime-Lite.md)
