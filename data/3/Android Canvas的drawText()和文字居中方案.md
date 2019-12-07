# Android Canvas的drawText()和文字居中方案

[![img](https://upload.jianshu.io/users/upload_avatars/5734256/2c5ca10c-8991-4770-bcef-4f8f0217b13b.jpg?imageMogr2/auto-orient/strip|imageView2/1/w/96/h/96/format/webp)](https://www.jianshu.com/u/c1664cb26499)

[biloba](https://www.jianshu.com/u/c1664cb26499)关注

32017.10.27 17:13:47字数 456阅读 14,016

自定义View是绘制文本有三类方法



```cpp
// 第一类
public void drawText (String text, float x, float y, Paint paint)
public void drawText (String text, int start, int end, float x, float y, Paint paint)
public void drawText (CharSequence text, int start, int end, float x, float y, Paint paint)
public void drawText (char[] text, int index, int count, float x, float y, Paint paint)

// 第二类
public void drawPosText (String text, float[] pos, Paint paint)
public void drawPosText (char[] text, int index, int count, float[] pos, Paint paint)

// 第三类
public void drawTextOnPath (String text, Path path, float hOffset, float vOffset, Paint paint)
public void drawTextOnPath (char[] text, int index, int count, Path path, float hOffset, float vOffset, Paint paint)
```

其中drawText()最常用，drawPosText ()是根据一个个坐标点指定文字位置，drawTextOnPath ()是根据路径绘制。但drawText()的x,y参数是干嘛的呢？

先来测试下



```dart
Paint paint=new Paint();
        paint.setStyle(Paint.Style.FILL);
        paint.setStrokeWidth(12);
        paint.setTextSize(100);

        String text="测试：my text";
        canvas.drawText(text, 200, 400, paint);

        //画两条线标记位置
        paint.setStrokeWidth(4);
        paint.setColor(Color.RED);
        canvas.drawLine(0, 400, 2000, 400, paint);
        paint.setColor(Color.BLUE);
        canvas.drawLine(200, 0, 200, 2000, paint);
```



![img](https://upload-images.jianshu.io/upload_images/5734256-034ce9861f52b1d4.png?imageMogr2/auto-orient/strip|imageView2/2/w/299/format/webp)

左对齐-left

可以看到，x,y并不是指定文字的中点位置，并且x,y与文字对齐方式有关（通过setTextAlign()指定，默认为left）



![img](https://upload-images.jianshu.io/upload_images/5734256-a6c0fc7234bbf914.png?imageMogr2/auto-orient/strip|imageView2/2/w/234/format/webp)

居中对齐-center



![img](https://upload-images.jianshu.io/upload_images/5734256-aa80878af37cf346.png?imageMogr2/auto-orient/strip|imageView2/2/w/241/format/webp)

右对齐-right

（为了使文字完整，上面调整了下x,y的值）

从上面三种情况得出结论，x所对应的竖线：

- 左对齐 — 文字的左边界
- 居中对齐 — 文字的中心位置
- 右对齐 — 文字的左边界

y对应的横线并不是文字的下边界，而是基准线Baseline

看下面这张图



![img](https://upload-images.jianshu.io/upload_images/5734256-512082b482008a3e.png?imageMogr2/auto-orient/strip|imageView2/2/w/1200/format/webp)

红色的Baseline是基准线，紫色的Top是文字的最顶部，也就是在drawText()中指定的x所对应，橙色的Bottom是文字的底部。
拿这些值如何获取呢？



```undefined
Paint.FontMetrics fontMetrics=paint.getFontMetrics();
        fontMetrics.top
        fontMetrics.ascent
        fontMetrics.descent
        fontMetrics.bottom
```

记得要在设置完Paint的文字大小，宽度之类属性后再获取FontMetrics，
baseline对应对应值为0，在它下面的descent和bottom值为正，top和ascent为负。那文字的高度为`bottom - top`



![img](https://upload-images.jianshu.io/upload_images/5734256-64f218ab8594bb3a.png?imageMogr2/auto-orient/strip|imageView2/2/w/839/format/webp)

所以，实际绘制的时候取决于基线上一个点来绘制文字,而这个点有三种分别对应为left,center,right



![img](https://upload-images.jianshu.io/upload_images/5734256-1b709c4384247865.png?imageMogr2/auto-orient/strip|imageView2/2/w/553/format/webp)

好啦，把drawText()中x,y参数讲清楚后实现文字居中就很容易了
直接上代码



```cpp
//矩形背景
        Paint bgRect=new Paint();
        bgRect.setStyle(Paint.Style.FILL);
        bgRect.setColor(Color.YELLOW);
        RectF rectF=new RectF(200, 200, 800, 600);
        canvas.drawRect(rectF, bgRect);

        Paint textPaint=new Paint();
        textPaint.setStyle(Paint.Style.FILL);
        textPaint.setStrokeWidth(8);
        textPaint.setTextSize(50);
        textPaint.setTextAlign(Paint.Align.CENTER);

        String text="测试：my text";
        //计算baseline
        Paint.FontMetrics fontMetrics=textPaint.getFontMetrics();
        float distance=(fontMetrics.bottom - fontMetrics.top)/2 - fontMetrics.bottom;
        float baseline=rectF.centerY()+distance;
        canvas.drawText(text, rectF.centerX(), baseline, textPaint);
```

效果



![img](https://upload-images.jianshu.io/upload_images/5734256-8bf933b5d441a516.png?imageMogr2/auto-orient/strip|imageView2/2/w/216/format/webp)

将对齐方式设置为center，那要让文字居中显示，x值就为矩形中心x值，y值也就是baseline的计算看下图



![img](https://upload-images.jianshu.io/upload_images/5734256-3c9cf9b18ff4e8db.png?imageMogr2/auto-orient/strip|imageView2/2/w/551/format/webp)

```
y = 矩形中心y值 + 矩形中心与基线的距离
```



```undefined
距离 = 文字高度的一半 - 基线到文字底部的距离（也就是bottom）
     = (fontMetrics.bottom - fontMetrics.top)/2 - fontMetrics.bottom
```