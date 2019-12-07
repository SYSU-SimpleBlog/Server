> 本文将简单了解下ViewPager的onPageScrolled方法

### onPageScrolled方法



```java
@Override
public void onPageScrolled(int position, float positionOffset, int positionOffsetPixels) {
        //position 当前所在页面
        //positionOffset 当前所在页面偏移百分比
        //positionOffsetPixels 当前所在页面偏移量
            }
```

系统是如何计算当前所在页面(position)，以及如何计算偏移百分比(positionOffset)



![img](https://upload-images.jianshu.io/upload_images/2007513-258e54bcdfd2c33f.gif?imageMogr2/auto-orient/strip|imageView2/2/w/744/format/webp)

单向滑动图.gif

!!!!!!留心红色短线 系统根据手机左边缘所处位置计算值
!!!!!!图很丑--!
0%（绿色页面完全在屏幕中）position: 0
~
25%(绿色页面的25%处已移出) position: 0
~
50%(绿色页面的50%处已移出) position: 0
~
75%(绿色页面的75%处已移出) position: 0
~
0%(绿色页面完全移出、黄色页面完全在屏幕中) position: 1
~
25%(黄色页面的25%处已移出) position: 1
~
50%(黄色页面的50%处已移出) position: 1
~
75%(黄色页面的75%处已移出) position: 1
~
0%(黄色页面完全移出、红色页面完全在屏幕中) position: 2



![img](https://upload-images.jianshu.io/upload_images/2007513-45dbd6540c59ea01.gif?imageMogr2/auto-orient/strip|imageView2/2/w/744/format/webp)

折返滑动图.gif

0%（黄色页面完全在屏幕中）position: 1
~
25%(黄色页面的25%处已移出) position: 1
~
50%(黄色页面的50%处已移出) position: 1
~
75%(黄色页面的75%处已移出) position: 1
~
0%(黄色页面完全移出、红色页面完全在屏幕中) position: 2
~
75%(黄色页面的75%处进入屏幕) position: 1
~
25%(黄色页面的25%处进入屏幕) position: 1
~
0%(黄色页面完全进入屏幕) position: 1
