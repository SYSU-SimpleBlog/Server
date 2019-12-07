# Flutter定义圆角的BottomSheet

[![img](https://upload.jianshu.io/users/upload_avatars/14802001/af144813-7861-4c84-b743-d901532c51de.jpeg?imageMogr2/auto-orient/strip|imageView2/1/w/96/h/96/format/webp)](https://www.jianshu.com/u/f516cfa3c49c)

[李良逸](https://www.jianshu.com/u/f516cfa3c49c)关注

2019.11.29 14:29:19字数 221阅读 23

## 显示 BottomSheet

我们在一些日常开发的场景中，可能需要使用到BottomSheet来显示菜单，就像这样。



![img](https://upload-images.jianshu.io/upload_images/14802001-5c00c20f98ab3954.png?imageMogr2/auto-orient/strip|imageView2/2/w/353/format/webp)

image.png

而在Flutter中，BottomSheet很好实现，只需要一行代码调用Flutter包中自带的BottomSheet显示方法showModalBottomSheet即可。



![img](https://upload-images.jianshu.io/upload_images/14802001-16894cc2fb74b468.png?imageMogr2/auto-orient/strip|imageView2/2/w/839/format/webp)

image.png

使用方法like this：



```dart
    void _showMenu(context, Menu menu) {
      showModalBottomSheet(
        context: context,
        builder: (context) => Column(
          mainAxisAlignment: MainAxisAlignment.spaceAround,
          children: <Widget>[
            Expanded(
              child: Container(
                child: _getMenuList(menu),
              ),
            ),
          ],
        ),
      );
    }
```

## 圆角 BottomSheet

但有时候，我们的视觉同学会会觉得这个BottomSheet的角不够圆润，我们想把BottomSheet上面左右两个角改为圆角，Android Native 实现效果如下：



![img](https://upload-images.jianshu.io/upload_images/14802001-782c960ea7f44b8b.gif?imageMogr2/auto-orient/strip|imageView2/2/w/160/format/webp)

1575008083736970.gif

而对于初学者来说，在设置圆角的路上，采用了一些网上比较坑的方法，或多或少都有雷区，其实showModalBottomSheet方法中的参数shape足以达到这个效果：



![img](https://upload-images.jianshu.io/upload_images/14802001-715c04a1e4a18e8d.gif?imageMogr2/auto-orient/strip|imageView2/2/w/160/format/webp)

1575008323387025.gif

我们自定义一个Shape，设置左上和右上圆角裁剪：



```dart
        shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.only(
              topLeft: Radius.circular(20.0),
              topRight: Radius.circular(20.0),
            ),
          ),
```

实现方法如下：



```dart
    void _showMenu(context, Menu menu) {
        showModalBottomSheet(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.only(
              topLeft: Radius.circular(20.0),
              topRight: Radius.circular(20.0),
            ),
          ),
          context: context,
          builder: (context) => Column(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            mainAxisSize: MainAxisSize.max,
            children: <Widget>[
              Center(child: Text('Test')),
            ],
          ),
        );
    }
```