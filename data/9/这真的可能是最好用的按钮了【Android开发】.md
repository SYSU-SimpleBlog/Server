#### 背景

按钮应该是我们的App里面最普遍的组件之一了，特别常用。

通常我们写一个按钮的套路很简单也很固定。大概分为以下几个步骤：

1. 在xml布局里面按照设计稿的尺寸位置写一个Textview
2. 按照设计稿规定的颜色和圆角在drawable目录下创建一个shape文件
3. 将这个shape文件作为Textview的背景

这样一个很标准的按钮就诞生了，然后就可以继续愉快的开发了。这本来没有什么问题，也比让UI妹纸切图高级了很多，但是随着开发的进行你会发现，UI妹纸的想法很多，不同的界面有各种不同圆角和不同背景颜色的按钮，这个时候你需要把上面的三个步骤再进行N次。

最后你会发现你的drawable目录下有各种各样的按钮背景资源，难以管理。特别是假如有的按钮要求有点击效果时需要使用selector，这个时候可能就会产生三个文件用来满足需求，所以总得来说很繁琐。

#### 想法

基于以上原因以及按钮使用的普遍程度，感觉很有必要写一个使用简单且能满足日常各种需求的按钮。我们先梳理下按钮需要达到的效果：

1. 使用简单(即可以利用属性对按钮进行各种设置)
2. 可以支持设置按钮文字、按钮文字颜色、按钮文字大小
3. 可以支持统一设置圆角大小，也可以单独设置按钮每个圆角的大小
4. 可以支持设置按钮背景颜色
5. 可以支持selector
6. 可以支持圆形按钮
7. 可以支持渐变色按钮，可以支持各个方向设置渐变色
8. 可以支持设置带边框的按钮，可以设置边框的颜色和宽度
9. 可以支持设置按钮是否可以点击
10. 可以设置带图标的按钮，支持自定义按钮大小，和自动缩放，图标支持设置在文字上下左右四个方向，支持自定义文字距离图标的距离

#### 引入



```bash
implementation 'top.androidman:superbutton:1.0.1'
```

#### 实现效果





![img](https://upload-images.jianshu.io/upload_images/15217452-ff2d30968ced94d4.png?imageMogr2/auto-orient/strip|imageView2/2/w/688/format/webp)

[image](https://links.jianshu.com/go?to=http%3A%2F%2Fwww.apkbus.com%2Fdata%2Fattachment%2Falbum%2F201907%2F18%2F110132u15ijxvc1nmjpzs6.png)



0x1 基本使用代码解释

##### 效果





![img](https://upload-images.jianshu.io/upload_images/15217452-21b3668458201685.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/362/format/webp)

[image](https://links.jianshu.com/go?to=http%3A%2F%2Fwww.apkbus.com%2Fdata%2Fattachment%2Falbum%2F201907%2F18%2F105921pd2gqggf08b2gh20.jpg)



##### 代码



```xml
<top.androidman.SuperButton
        android:layout_width="match_parent"
        android:layout_height="60dp"
        android:layout_margin="20dp"
        app:color_normal="@color/color_accent"
        app:corner="10dp"
        app:text="@string/poetry_1"
        app:textColor="@color/color_white"
        app:textSize="22sp" />
```

##### 属性解释

- 按钮文字
  app:text="床前明月光"
- 按钮文字颜色
  app:textColor="@color/color_white"
- 按钮文字大小
  app:textSize="22sp"
- 按钮背景颜色
  app:color_normal="@color/color_accent"

#### 0x2 单独设置每个圆角

##### 效果





![img](https://upload-images.jianshu.io/upload_images/15217452-02decbef180954f8.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/364/format/webp)

[image](https://links.jianshu.com/go?to=http%3A%2F%2Fwww.apkbus.com%2Fdata%2Fattachment%2Falbum%2F201907%2F18%2F105925fle7oto7nnx6tolt.jpg)



##### 代码



```xml
<top.androidman.SuperButton
    android:layout_width="match_parent"
    android:layout_height="60dp"
    android:layout_margin="20dp"
    app:color_normal="@color/color_6596ff"
    app:corner="40dp"
    app:corner_left_bottom="0dp"
    app:text="单独设置左下角为0dp"
    app:textColor="@color/color_white"
    app:textSize="22sp" />
```

##### 属性解释

- 单独设置左下角角度
  app:corner_left_bottom="0dp"
- 单独设置左上角角度
  app:corner_left_top="5dp"
- 单独设置右上角角度
  app:corner_right_top="22dp"
- 单独设置右下角角度
  app:corner_right_bottom="0dp"
- 按钮四个角的圆角角度
  app:corner="10dp"

注意：单独设置角度会覆盖corner属性

#### 0x3 Selector

##### 效果





![img](https://upload-images.jianshu.io/upload_images/15217452-5529143edc7b79c1.gif?imageMogr2/auto-orient/strip|imageView2/2/w/384/format/webp)

[image](https://links.jianshu.com/go?to=http%3A%2F%2Fwww.apkbus.com%2Fdata%2Fattachment%2Falbum%2F201907%2F18%2F105931wpd8dldvfpy7gpsp.gif)



##### 代码



```xml
<top.androidman.SuperButton
    android:layout_width="match_parent"
    android:layout_height="60dp"
    android:layout_margin="20dp"
    android:onClick="test"
    app:color_normal="@color/color_accent"
    app:color_pressed="@color/color_ffB419"
    app:corner="10dp"
    app:text="点击会变色哦"
    app:textColor="@color/color_white"
    app:textSize="22sp" />
```

##### 属性解释

- 当按钮点下时会显示设置的颜色效果
  app:color_normal="@color/color_accent"

#### 0x4 圆形按钮

##### 效果





![img](https://upload-images.jianshu.io/upload_images/15217452-67e04fb359535f22.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/368/format/webp)

[image](https://links.jianshu.com/go?to=http%3A%2F%2Fwww.apkbus.com%2Fdata%2Fattachment%2Falbum%2F201907%2F18%2F105934s6zpw6vgx6m6fuu2.jpg)



##### 代码



```xml
 <top.androidman.SuperButton
    android:layout_width="160dp"
    android:layout_height="160dp"
    android:layout_margin="20dp"
    android:onClick="test"
    app:color_normal="@color/color_accent"
    app:drawable_middle="@mipmap/icon_like"
    app:drawable_middle_height="50dp"
    app:drawable_middle_width="50dp"
    app:shape="CIRCLE" />
```

##### 属性解释

- 按钮上只有图标没有文字
  app:drawable_middle="@mipmap/icon_like"
- 按钮上图标高度
  app:drawable_middle_height="50dp"
- 按钮上图标宽度
  app:drawable_middle_width="50dp"

注意：当设置此属性时，文字的相关属性将会失效

#### 0x5 渐变背景的按钮

##### 效果





![img](https://upload-images.jianshu.io/upload_images/15217452-c48d5831d1696c24.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/376/format/webp)

[image](https://links.jianshu.com/go?to=http%3A%2F%2Fwww.apkbus.com%2Fdata%2Fattachment%2Falbum%2F201907%2F18%2F110356gz55ec65we86rc6y.jpg)



##### 代码



```xml
<top.androidman.superbutton.SuperButton
    android:layout_width="match_parent"
    android:layout_height="60dp"
    android:layout_margin="20dp"
    app:color_direction="RIGHT_LEFT"
    app:color_end="@color/color_14c7de"
    app:color_start="@color/color_red"
    app:corner="20dp"
    app:text="从右到左渐变"
    app:textColor="@color/color_white"
    app:textSize="22sp" />
```

##### 属性解释

- 开始颜色
  app:color_start="@color/color_14c7de"
- 结束颜色
  app:color_end="@color/color_red"
- 颜色渐变方向
  app:color_direction="RIGHT_LEFT"
- TOP_BOTTOM 从上到下
  BOTTOM_TOP 从下到上
  LEFT_RIGHT 从左到右
  RIGHT_LEFT 从右到左
  TR_BL 从右上到左下
  BL_TR 从左下到右上
  BR_TL 从右下到左上
  TL_BR 从左上到右下

注意：当设置颜色渐变时，color_normal，color_pressed设置将失效

#### 0x6 有边框按钮

##### 效果





![img](https://upload-images.jianshu.io/upload_images/15217452-389a3575c2d4ce6b.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/372/format/webp)

[image](https://links.jianshu.com/go?to=http%3A%2F%2Fwww.apkbus.com%2Fdata%2Fattachment%2Falbum%2F201907%2F18%2F110359kz8x6i2rg8gf5f76.jpg)



##### 代码



```xml
<top.androidman.SuperButton
    android:layout_width="match_parent"
    android:layout_height="60dp"
    android:layout_margin="20dp"
    app:border_color="@color/color_red"
    app:border_width="2dp"
    app:color_normal="@color/color_accent"
    app:corner="10dp"
    app:text="@string/poetry_1"
    app:textColor="@color/color_white"
    app:textSize="22sp" />
```

##### 属性解释

- 边框宽度
  app:border_width="2dp"
- 边框颜色
  app:border_color="@color/color_red"

#### 0x7 按钮不可点击

##### 效果





![img](https://upload-images.jianshu.io/upload_images/15217452-e8485ba55078b11a.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/364/format/webp)

[image](https://links.jianshu.com/go?to=http%3A%2F%2Fwww.apkbus.com%2Fdata%2Fattachment%2Falbum%2F201907%2F18%2F110402ktdjfc7f70w70m7t.jpg)



##### 代码



```xml
<top.androidman.SuperButton
    android:layout_width="match_parent"
    android:layout_height="60dp"
    android:layout_margin="20dp"
    android:onClick="test"
    app:button_click_able="false"
    app:color_normal="@color/color_accent"
    app:corner="10dp"
    app:singleLine="false"
    app:text="app:button_click_able=false\n也可以代码设置"
    app:textColor="@color/color_white"
    app:textSize="22sp" />
```

##### 属性解释

- 按钮是否可以点击，默认为true可以点击
  app:button_click_able="false"
- 按钮文字是否单行，默认为true单行
  app:singleLine="false"

注意：按钮提供有方法对点击事件进行控制，请参考高级使用里面相关方法

#### 0x8 带图标按钮

##### 效果



![img](https://upload-images.jianshu.io/upload_images/15217452-870faa0384e932a5.png?imageMogr2/auto-orient/strip|imageView2/2/w/376/format/webp)

image.png

##### 代码



```xml
<top.androidman.superbutton.SuperButton
    android:layout_width="match_parent"
    android:layout_height="60dp"
    android:layout_margin="20dp"
    app:color_normal="@color/color_red"
    app:corner="20dp"
    app:drawable_padding="20dp"
    app:drawable_right="@mipmap/icon_like"
    app:text="图标在右边"
    app:textColor="@color/color_white"
    app:textSize="22sp" />
```

#### 属性解释

- 图标在文字右边
  app:drawable_right="@mipmap/icon_like"
- 图标在文字左边
  app:drawable_left="@mipmap/icon_like"
- 图标在文字上边
  app:drawable_top="@mipmap/icon_like"
- 图标在文字下边
  app:drawable_bottom="@mipmap/icon_like"
- 图标距文字距离
  app:drawable_padding="20dp"
- 根据文字大小缩放图标，默认为true,当为false时显示原图标大小
  app:drawable_auto="true"

##### 按钮支持的所有属性



```xml
<?xml version="1.0" encoding="utf-8"?>
<resources>
    <declare-styleable name="SuperButton">
        <!--默认配置-->
        <attr name="text" format="reference|string" />
        <!--按钮文字颜色-->
        <attr name="textColor" format="reference|color" />
        <!--按钮文字大小-->
        <attr name="textSize" format="dimension" />
        <!--文字是否单行，默认单行-->
        <attr name="singleLine" format="boolean" />

        <!--默认背景颜色-->
        <attr name="color_normal" format="reference|color" />
        <!--按压时的背景颜色-->
        <attr name="color_pressed" format="reference|color" />

        <!--图片在文字左边-->
        <attr name="drawable_left" format="reference|color" />
        <!--图片在文字右边-->
        <attr name="drawable_right" format="reference|color" />
        <!--图片在文字上边-->
        <attr name="drawable_top" format="reference|color" />
        <!--图片在文字下边-->
        <attr name="drawable_bottom" format="reference|color" />
        <!--图片距文字的距离-->
        <attr name="drawable_padding" format="dimension" />
        <!--图标根据文字大小自动缩放图标-->
        <attr name="drawable_auto" format="boolean" />

        <!--只有图片的情况,此时会忽略文字，即便设置-->
        <attr name="drawable_middle" format="reference|color" />
        <!--图片在中间时宽-->
        <attr name="drawable_middle_width" format="dimension" />
        <!--图片在中间时高-->
        <attr name="drawable_middle_height" format="dimension" />

        <!--形状-->
        <attr name="shape" format="enum">
            <!--圆形-->
            <enum name="CIRCLE" value="1" />
            <!--矩形-->
            <enum name="RECT" value="2" />
        </attr>

        <!--按钮背景是渐变色时设置-->
        <!--开始颜色-->
        <attr name="color_start" format="color" />
        <!--结束颜色-->
        <attr name="color_end" format="color" />
        <!--颜色方向-->
        <attr name="color_direction" format="enum">
            <!--从上到下-->
            <enum name="TOP_BOTTOM" value="0x1" />
            <!--从右上到左下-->
            <enum name="TR_BL" value="0x2" />
            <!--从右到左-->
            <enum name="RIGHT_LEFT" value="0x3" />
            <!--从右下到左上-->
            <enum name="BR_TL" value="0x4" />
            <!--从下到上-->
            <enum name="BOTTOM_TOP" value="0x5" />
            <!--从左下到右上-->
            <enum name="BL_TR" value="0x6" />
            <!--从左到右-->
            <enum name="LEFT_RIGHT" value="0x7" />
            <!--从左上到右下-->
            <enum name="TL_BR" value="0x8" />
        </attr>

        <!--按钮圆角，如果单独设置会覆盖此设置-->
        <attr name="corner" format="dimension" />
        <!--左上角圆角半径-->
        <attr name="corner_left_top" format="dimension" />
        <!--左下角圆角半径-->
        <attr name="corner_left_bottom" format="dimension" />
        <!--右上角圆角半径-->
        <attr name="corner_right_top" format="dimension" />
        <!--右下角圆角半径-->
        <attr name="corner_right_bottom" format="dimension" />

        <!--边框宽度-->
        <attr name="border_width" format="dimension" />
        <!--边框颜色-->
        <attr name="border_color" format="color" />
        <!--按钮是否可以点击-->
        <attr name="button_click_able" format="boolean" />

    </declare-styleable>
</resources>
```

#### 高级应用

1.想修改按钮相关调用如下方法：



```css
 /**
     * 修改文字
     */
    superButton.setText(CharSequence text);

    /**
     * 修改文字颜色
     */
    superButton.setTextColor(@ColorInt int textColor);

    /**
     * 修改按钮背景颜色
     */
    superButton.setColorNormal(@ColorInt int colorNormal);
```

当某些状态下需要代码控制，将按钮置灰然后不可点击，可以直接调用如下方法：



```css
 /**
     * 调用此方法后按钮颜色被改变，按钮无法点击
     */
    superButton.setUnableColor(@ColorInt int color);
```

如果只是想设置按钮不可点击，不需要改变按钮颜色，可以使用如下方法



```css
/**
     * 设置按钮是否可以点击
     */
    superButton.setButtonClickable(boolean buttonClickable);
```

# 最后

**给大家分享一份移动架构大纲，包含了移动架构师需要掌握的所有的技术体系，大家可以对比一下自己不足或者欠缺的地方有方向的去学习提升；**

**需要高清架构图以及图中视频资料和文章项目源码的可以加入我的技术交流群：825106898私聊群主小姐姐免费获取**



![img](https://upload-images.jianshu.io/upload_images/15217452-3e4ac16c885c487a.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/720/format/webp)
