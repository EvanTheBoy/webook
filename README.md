## 环境配置说明

本项目原本是在 Windows 电脑上跑的，但由于后续用到的技术栈包括 redis，docker 和 Kubernetes等，而 Windows 电脑上安装它们实在是太麻烦，因此考虑使用服务器，将这些服务全部部署在服务器上。经过长时间的测试，我目前采用的是 Ubuntu 22.04 的虚拟机，配置是 4 核，22GB 内存，50GB 硬盘。

如果你决定将这个项目下载到本地并把它跑起来，那么接下来的内容就是为你准备的！

## 在 Linux 安装 Go 语言

### 基本流程

安装 Go 语言的方式有很多种，这里使用从阿里镜像源获取之后再手动解压的方式。

阿里镜像源：https://mirrors.aliyun.com/golang/?spm=a2c6h.13651104.d-5243.7.1d351e57KtELYq

由于我的 Windows 电脑上安装的是 go1.20.5，因此在服务器上也安装这个版本：

![image-20240207101224893](F:\webookImages\image-20240207101224893.png)

直接点击下载即可。

接下来在服务器上打开终端，进入 Downloads 目录，输入下面的命令在指定目录解压：

```bash
sudo tar -zxvf go1.20.5.linux-amd64.tar.gz -C /usr/local/
```

接下来输入：

sudo vim /etc/profile

配置环境变量。在文件的末尾加入：

```
export GOROOT=/usr/local/go

export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

保存退出后，输入：

```bash
sudo source /etc/profile
```

使其生效。

### 可能出现的 bug

有可能出现每次进入终端都需要执行 source 命令使配置文件生效，才能真的生效的 bug，可以考虑：

```bash
sudo vim ~/.bashrc
```

然后将 source /etc/profile 加入文件末尾。

## 在 Goland 进行远程部署

### ssh

首先，在 linux 系统安装支持 ssh 相关的软件：

sudo apt-get install openssh-server

然后输入：

ps -ef | grep sshd

查看是否安装成功：

![image.png](F:\webookImages\sshsupport.png)

### 部署

打开 Goland，点击 Tools -> Deployment









