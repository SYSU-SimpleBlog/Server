在Android开发中，系统对Activity、Fragment的生命周期有着非常明显且较于区分的定义，但是在flutter中，由于flutter的生命周期依附在activity或fragment，它的生命周期就不同以往了，下面就展示以下flutter生命周期的理解。

## flutter 生命周期相关函数的调用过程

首先，先上一张图，这张图很简单明了的阐释了一个页面启动所要执行的widget方法流程：





![img](https://upload-images.jianshu.io/upload_images/18813666-b8f11f404d93cb36.png?imageMogr2/auto-orient/strip|imageView2/2/w/485/format/webp)

下面解释一下各个方法的作用：

### initState

在生命周期中只调用一次，此时无法获取widget对象，可以做一些初始化操作。

### didChangeDependencies

当State对象的依赖发生变化时会被调用；例如：在之前build() 中包含了一个InheritedWidget，然后在之后的build() 中InheritedWidget发生了变化，那么此时InheritedWidget的子widget的didChangeDependencies()回调都会被调用。InheritedWidget这个widget可以由父控件向子控件共享数据，案例可以参考 scoped_model开源库。

### didUpdateWidget

widget状态改变的时候调用

### deactivate

类似于Activity的onResume和onStop，两种状态都会调用

### dispose

类似于Android的onDestroy

上面的介绍都比较简单，下面则介绍以下，如何去获取app的生命周期

## flutter app的生命周期

flutter提供了一个枚举类来代表了app各个生命周期的状态：



```php
enum AppLifecycleState {
  /// The application is visible and responding to user input.
  resumed,

  /// The application is in an inactive state and is not receiving user input.
  ///
  /// On iOS, this state corresponds to an app or the Flutter host view running
  /// in the foreground inactive state. Apps transition to this state when in
  /// a phone call, responding to a TouchID request, when entering the app
  /// switcher or the control center, or when the UIViewController hosting the
  /// Flutter app is transitioning.
  ///
  /// On Android, this corresponds to an app or the Flutter host view running
  /// in the foreground inactive state.  Apps transition to this state when
  /// another activity is focused, such as a split-screen app, a phone call,
  /// a picture-in-picture app, a system dialog, or another window.
  ///
  /// Apps in this state should assume that they may be [paused] at any time.
  inactive,

  /// The application is not currently visible to the user, not responding to
  /// user input, and running in the background.
  ///
  /// When the application is in this state, the engine will not call the
  /// [Window.onBeginFrame] and [Window.onDrawFrame] callbacks.
  ///
  /// Android apps in this state should assume that they may enter the
  /// [suspending] state at any time.
  paused,

  /// The application will be suspended momentarily.
  ///
  /// When the application is in this state, the engine will not call the
  /// [Window.onBeginFrame] and [Window.onDrawFrame] callbacks.
  ///
  /// On iOS, this state is currently unused.
  suspending,
}
```

### resumed

应用程序对用户可见的时候输出

### inactive

界面处于不可点击状态，但是可见时候的回调，类似于Android的onpause

### paused

app处于不可见的时候，类似于Android的onStop

### suspending

ios中这个属性无效，android中代表处于后台

### 获取方法



```java
class _MyHomePageState extends State<MyHomePage> with WidgetsBindingObserver {
  AppLifecycleState _lastLifecycleState;

  void dispose() {
    super.dispose();
    WidgetsBinding.instance.removeObserver(this);
  }

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    print(state);
  }

  ...
}
```

注意：第一次进入的时候并不会执行didChangeAppLifecycleState方法。
获取app的生命周期方法很简单，但是注意这并不是当前widget的生命周期，那我们如果获取当前页面的生命周期呢。

## 获取flutter页面的生命周期

当flutter页面跳转切入后台，flutter并没有清楚的给我们展示flutter页面的各个生命周期状态。如果我们想要获取某个widget页面的状态，比如可见不可见那该如何操作呢？

### flutter页面的onDestroy

这个比较简单，重写State的dispose，这个方法即可理解为页面的onDestroy操作。

### flutter页面的onStop、onResume

上面介绍了deactivate类似于activity的onResume、onStop那么我们可以利用这个函数来自己标志一下生命周期。
因为deactivate这个方法第一次是不执行的，因此我们可以定义一个默认值isVisible为true来代表是否可见。



```java
class MyState extends State<StatefulWidget>{

  bool isVisible = true;

  @override
  void deactivate() {
    isVisible = !isVisible;
    if(isVisible){
      //onResume
    }else {
      //onStop
    }
    super.deactivate();
  }

  @override
  Widget build(BuildContext context) {
    // TODO: implement build
    return null;
  }

}
```

这时候我们就可以通过isVisible来判断当前页面是否可见了，以此来做一些操作。
今年金九银十我花一个月的时间收录整理了一套知识体系，如果有想法深入的系统化的去学习的，可以点击[传送门](https://links.jianshu.com/go?to=https%3A%2F%2Fshimo.im%2Fdocs%2FkxRG3DxvRTkWhqc9%2F)，我会把我收录整理的资料都送给大家，帮助大家更快的进阶。



![img](https://upload-images.jianshu.io/upload_images/18813666-21f1bdbabcf5a931.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/240/format/webp)
