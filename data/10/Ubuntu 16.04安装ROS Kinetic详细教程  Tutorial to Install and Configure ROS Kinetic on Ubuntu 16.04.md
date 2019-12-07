[![img](https://upload.jianshu.io/users/upload_avatars/20059030/60ad3b7b-5caf-4ff7-890c-6983c6385619.jpg?imageMogr2/auto-orient/strip|imageView2/1/w/96/h/96/format/webp)](https://www.jianshu.com/u/12461cd06bc9)

[kezunlin](https://www.jianshu.com/u/12461cd06bc9)关注

**本文首发于个人博客https://kezunlin.me/post/e2780b93/，欢迎阅读!**

Tutorial to Install and Configure ROS Kinetic on Ubuntu 16.04.



```undefined
ROS release       ROS version         Ubuntu version
2016.3       ROS Kinetic Kame   Ubuntu 16.04(Xenial)/Ubuntu 15.10
```

# Install Guide

## quick commands



```swift
sudo sh -c '. /etc/lsb-release && echo "deb [arch=amd64] http://mirrors.ustc.edu.cn/ros/ubuntu/ $DISTRIB_CODENAME main" > /etc/apt/sources.list.d/ros-latest.list'

sudo apt-key adv --keyserver hkp://ha.pool.sks-keyservers.net:80 --recv-key 421C365BD9FF1F717815A3895523BAEEB01FA116

sudo apt-get update
sudo apt-get install ros-kinetic-desktop-full

sudo apt-get install python-rosinstall python-rosinstall-generator python-wstool build-essential

sudo rosdep init
rosdep update

echo "source /opt/ros/kinetic/setup.bash" >> ~/.bashrc
source ~/.bashrc

# test
roscore
```

## Notice

### tinghua source



```cpp
http://mirrors.tuna.tsinghua.edu.cn/ubuntu xenial
```

### update source



```csharp
sudo apt-get update
```

if we encouter errors when update source, we need to fix.

e.g remove sougou source to fix errors.



```cpp
grep -r kylin .
./sources.list.d/sogoupinyin.list.save:deb http://archive.ubuntukylin.com:10006/ubuntukylin xenial main
./sources.list.d/sogoupinyin.list:deb http://archive.ubuntukylin.com:10006/ubuntukylin xenial main

rm ./sources.list.d/sogoupinyin.list
```

## Configure ros source

### ros offical(NOT RECOMMEND)



```rust
sudo sh -c 'echo "deb http://packages.ros.org/ros/ubuntu $(lsb_release -sc) main" > /etc/apt/sources.list.d/ros-latest.list'
```

NOT RECOMMEND，when `apt-get update`, error will occur `Hash Sum mismatch`

### ros china(RECOMMEND)



```dart
sudo sh -c '. /etc/lsb-release && echo "deb [arch=amd64] http://mirrors.ustc.edu.cn/ros/ubuntu/ $DISTRIB_CODENAME main" > /etc/apt/sources.list.d/ros-latest.list'
sudo apt-key adv --keyserver hkp://ha.pool.sks-keyservers.net:80 --recv-key 421C365BD9FF1F717815A3895523BAEEB01FA116
```

> ros-latest.list



```ruby
deb [arch=amd64] http://mirrors.ustc.edu.cn/ros/ubuntu/ xenial main
```

## Install ros

ROS, rqt, rviz, robot-generic libraries, 2D/3D simulators, navigation and 2D/3D perception



```csharp
sudo apt-get update
sudo apt-get install ros-kinetic-desktop-full

#sudo apt-get install ros-kinetic-desktop
#sudo apt-get install ros-kinetic-ros-base

#sudo apt-get install ros-kinetic-<PACKAGE>
#sudo apt-get install ros-kinetic-slam-gmapping

#apt-cache search ros-kinetic
```

## Initialize rosdep

Before you can use ROS, you will need to initialize rosdep. rosdep enables you to easily install system dependencies for source you want to compile and is required to run some core components in ROS.



```kotlin
sudo rosdep init
rosdep update
```

will output



```cpp
reading in sources list data from /etc/ros/rosdep/sources.list.d
Hit https://raw.githubusercontent.com/ros/rosdistro/master/rosdep/osx-homebrew.yaml
Hit https://raw.githubusercontent.com/ros/rosdistro/master/rosdep/base.yaml
Hit https://raw.githubusercontent.com/ros/rosdistro/master/rosdep/python.yaml
Hit https://raw.githubusercontent.com/ros/rosdistro/master/rosdep/ruby.yaml
Hit https://raw.githubusercontent.com/ros/rosdistro/master/releases/fuerte.yaml
Query rosdistro index https://raw.githubusercontent.com/ros/rosdistro/master/index.yaml
Add distro "groovy"
Add distro "hydro"
Add distro "indigo"
Add distro "jade"
Add distro "kinetic"
Add distro "lunar"
updated cache in /home/kezunlin/.ros/rosdep/sources.cache
```

## Environment setup



```bash
echo "source /opt/ros/kinetic/setup.bash" >> ~/.bashrc
source ~/.bashrc
```

### check ROS



```bash
env | grep ROS
export | grep ROS
declare -x ROSLISP_PACKAGE_DIRECTORIES=""
declare -x ROS_DISTRO="kinetic"
declare -x ROS_ETC_DIR="/opt/ros/kinetic/etc/ros"
declare -x ROS_MASTER_URI="http://localhost:11311"
declare -x ROS_PACKAGE_PATH="/opt/ros/kinetic/share"
declare -x ROS_ROOT="/opt/ros/kinetic/share/ros"
```

## Dependencies for building packages



```csharp
sudo apt-get install python-rosinstall python-rosinstall-generator python-wstool build-essential
```

## Test install



```undefined
roscore
```

output



```cpp
... logging to /home/kezunlin/.ros/log/b777db6c-ff85-11e8-93c2-80fa5b47928a/roslaunch-ke-17139.log
Checking log directory for disk usage. This may take awhile.
Press Ctrl-C to interrupt
Done checking log file disk usage. Usage is <1GB.

started roslaunch server http://ke:36319/
ros_comm version 1.12.14


SUMMARY
========

PARAMETERS
 * /rosdistro: kinetic
 * /rosversion: 1.12.14

NODES

auto-starting new master
process[master]: started with pid [17162]
ROS_MASTER_URI=http://ke:11311/

setting /run_id to b777db6c-ff85-11e8-93c2-80fa5b47928a
process[rosout-1]: started with pid [17175]
started core service [/rosout]
^C[rosout-1] killing on exit
[master] killing on exit
shutting down processing monitor...
... shutting down processing monitor complete
done
```

# Create Workspace

## Create

Let's create and build a catkin workspace:



```jsx
mkdir -p ~/catkin_ws/src
cd ~/catkin_ws/
catkin_make

ls .
build dist src
```

`tree src` folder



```jsx
src/
└── CMakeLists.txt -> /opt/ros/kinetic/share/catkin/cmake/toplevel.cmake

0 directories, 1 file
```

`tree devel` folder



```css
devel
├── env.sh
├── lib
├── setup.bash
├── setup.sh
├── _setup_util.py
└── setup.zsh

1 directory, 5 files
```

The `catkin_make` command is a convenience tool for working with `catkin workspaces`.

## source devel setup

before `source devel/setup.bash`



```jsx
env | grep ROS
ROS_ROOT=/opt/ros/kinetic/share/ros
ROS_PACKAGE_PATH=/opt/ros/kinetic/share
ROS_MASTER_URI=http://localhost:11311
ROSLISP_PACKAGE_DIRECTORIES=
ROS_DISTRO=kinetic
ROS_ETC_DIR=/opt/ros/kinetic/etc/ros
```

after `source devel/setup.bash`



```jsx
env | grep ROS
ROS_ROOT=/opt/ros/kinetic/share/ros
ROS_PACKAGE_PATH=/home/kezunlin/catkin_ws/src:/opt/ros/kinetic/share
ROS_MASTER_URI=http://localhost:11311
ROSLISP_PACKAGE_DIRECTORIES=/home/kezunlin/catkin_ws/devel/share/common-lisp
ROS_DISTRO=kinetic
ROS_ETC_DIR=/opt/ros/kinetic/etc/ros
```

To make sure your workspace is properly overlayed by the setup script, make sure `ROS_PACKAGE_PATH` environment variable includes the directory you're in.



```bash
echo $ROS_PACKAGE_PATH
/home/kezunlin/catkin_ws/src:/opt/ros/kinetic/share
```

# Reference

- [Official Install Guide](https://links.jianshu.com/go?to=http%3A%2F%2Fwiki.ros.org%2Fkinetic%2FInstallation%2FUbuntu)
- [Ubuntu ROS mirrors](https://links.jianshu.com/go?to=http%3A%2F%2Fwiki.ros.org%2FROS%2FInstallation%2FUbuntuMirrors)
- [Configure ROS Environment](https://links.jianshu.com/go?to=http%3A%2F%2Fwiki.ros.org%2FROS%2FTutorials%2FInstallingandConfiguringROSEnvironment)

# History

- 2018/01/04: created.

# Copyright

- Post author: [kezunlin](https://links.jianshu.com/go?to=https%3A%2F%2Fkezunlin.me)
- Post link: [https://kezunlin.me/post/e2780b93/](https://links.jianshu.com/go?to=https%3A%2F%2Fkezunlin.me%2Fpost%2Fe2780b93%2F)
- Copyright Notice: All articles in this blog are licensed under CC BY-NC-SA 3.0 unless stating additionally.