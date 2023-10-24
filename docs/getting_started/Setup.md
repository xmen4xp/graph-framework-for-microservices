# Playground Lite

[[Exit]](../../README.md)  

Welcome to the Nexus Setup.

This document will walk you through the setup needed to get your developer env ready to work with nexus.


## Compiling using Linux OS behind INTEL VPN

### setup the apt repository to use proxy (Add/update the apt conf file)
```
> sudo vi /etc/apt/apt.conf
Acquire::http::Proxy "http://proxy-dmz.intel.com:911";
Acquire::https::Proxy "http://proxy-dmz.intel.com:912";
```

### add the following lines to your shell startup script (.bashrc or .zshrc)
```
export http_proxy=http://proxy-dmz.intel.com:911
export https_proxy=http://proxy-dmz.intel.com:912
export ftp_proxy=http://proxy-dmz.intel.com:911
export socks_proxy=http://proxy-dmz.intel.com:1080
export no_proxy=intel.com,.intel.com,localhost,127.0.0.1,10.0.0.0/8,192.168.0.0/16,172.16.0.0/12
```

### Install the proxy connect package 
```
> sudo apt install connect-proxy
```

### Update your ssh to use proxy when pulling images from github, add the following to ~/.ssh/config
```
Host github.com
     ProxyCommand connect -a none -S proxy-dmz.intel.com:1080 %h %p
```

### Update the docker daemon to use the proxy for pulling docker images (add to /etc/docker/daemon.json)
```
{
 "proxies": {
     "http-proxy": "http://proxy-dmz.intel.com:911",
     "https-proxy": "http://proxy-dmz.intel.com:912",
     "no-proxy": "*.test.example.com,.example.org,intel.com,.intel.com,localhost,127.0.0.1,10.0.0.0/8,192.168.0.0/16,172.16.0.0/12"
 }
}
```

### Restart the docker deamon for the chanages to take effect
```
> sudo systemctl restart docker.service
```

[[Exit]](../../README.md) 
