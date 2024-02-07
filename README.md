## 环境配置说明

本项目原本是在 Windows 电脑上跑的，但由于后续用到的技术栈包括 redis，docker 和 Kubernetes 等，而 Windows 电脑上安装它们实在是太麻烦，因此考虑使用服务器，将这些服务全部部署在服务器上。经过长时间的测试，我目前采用的是 Ubuntu 22.04 的虚拟机，配置是 4 核，22GB 内存，50GB 硬盘。

如果你决定将这个项目下载到本地并把它跑起来，那么接下来的内容就是为你准备的！

## 在 Linux 安装 Go 语言

### 基本流程

安装 Go 语言的方式有很多种，这里使用从阿里镜像源获取之后再手动解压的方式。

阿里镜像源：https://mirrors.aliyun.com/golang/?spm=a2c6h.13651104.d-5243.7.1d351e57KtELYq

由于我的 Windows 电脑上安装的是 go1.20.5，因此在服务器上也安装这个版本：

![image-20240207101224893](https://github.com/EvanTheBoy/webook/assets/73733942/cd542eff-daf4-492d-947d-ef7b06a31f6f)

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

![sshsupport](https://github.com/EvanTheBoy/webook/assets/73733942/2c1580a5-4cab-4200-adc3-075ff38d3ae0)

### 部署

打开 Goland，打开 Tools -> Deployment，点击左上角的 + 号，选择 SFTP 协议，给服务器命一个名，然后输入远程服务器的 IP 地址，最后保存：

![deployment](https://github.com/EvanTheBoy/webook/assets/73733942/133b4c16-48e8-4733-af39-bcfa2dbd7273)

接着从 Connection 选择到 Mappings：

![mappings](https://github.com/EvanTheBoy/webook/assets/73733942/82b1311f-bbea-476b-8da4-9a766f89af87)

选择 Deployment path，选择 Windows 和 Linux 上的文件对应。然后点击 OK。

接下来右键根目录，然后 Deployment -> Sync with Deployed to ubuntu，把所有的文件全部传到 Linux 对应的目录中去：

![upload](https://github.com/EvanTheBoy/webook/assets/73733942/1277fcab-bf22-4c58-b137-96a6c7264d8e)

## Linux 终端设置代理

打开终端，输入：

```bash
sudo apt install proxychains4
```

然后输入：

```bash
proxychains4 -help
```

查看其配置文件的所在位置。

这里可以下载 mousepad 编辑器，用 mousepad 打开：









