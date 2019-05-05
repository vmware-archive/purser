# Purser Plugin Setup
_NOTE: This Plugin installation is optional. Install it if you want to use CLI of Purser._

## Linux and macOS

``` bash
# Binary installation
wget -q https://github.com/vmware/purser/blob/master/build/purser-binary-install.sh && sh purser-binary-install.sh
```

Enter your cluster's configuration path when prompted. The plugin binary needs to be in your `PATH` environment variable, so once the download of the binary is finished the script tries to move it to `/usr/local/bin`. This may need your sudo permission.

## Windows/Others

For installation on Windows follow the steps in the [manual installation guide](./docs/manual-installation.md).

## Uninstalling Purser Plugin

### Linux/macOS

``` bash
wget -q https://github.com/vmware/purser/blob/master/build/purser-binary-uninstall.sh && sh purser-binary-uninstall.sh
```
