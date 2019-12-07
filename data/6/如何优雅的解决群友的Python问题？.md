# 如何优雅的解决群友的Python问题？

[![img](https://upload.jianshu.io/users/upload_avatars/3629157/51226190-fced-4a0e-8246-65a70e65fcf0.png?imageMogr2/auto-orient/strip|imageView2/1/w/96/h/96/format/webp)](https://www.jianshu.com/u/9104ebf5e177)

[罗罗攀](https://www.jianshu.com/u/9104ebf5e177)关注

0.8632019.11.29 10:33:23字数 229阅读 62



![img](https://upload-images.jianshu.io/upload_images/3629157-0ca8e5025c28a4b4.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/1200/format/webp)

##### 01 问题描述

这个问题来源于自己Python交流群中的一个问题，如下图所示，需要计算每列中各值的出现次数，然后组成一个新的表。



![img](https://upload-images.jianshu.io/upload_images/3629157-b3b4c4ac08936b43.png?imageMogr2/auto-orient/strip|imageView2/2/w/1061/format/webp)

##### 02 解决思路

计算每列各值的出现次数，我们可以使用groupby方法，当然最简单的还是使用value_counts方法。

- 首先读取数据
- 接着使用一个循环语句，依次计算每列的值计算
- （由于每列的值计数返回的是series数据，而且我们也需要在结果表中的一列加上列名），构建每列值计数的dataframe。
- 最后将这些dataframe合并即可。

##### 03 解决代码



```kotlin
import pandas as pd

data = pd.read_excel('例子.xlsx',sheetname='Sheet1',index_col='index')

frames = []
for i in data.columns:
    s = data[i].value_counts().sort_values()
    d = pd.DataFrame({'列名':i,'变量名':s.index,'次数':s.values})
    frames.append(d)
    
result = pd.concat(frames)
result
```



![img](https://upload-images.jianshu.io/upload_images/3629157-09c2d795995ba79e.png?imageMogr2/auto-orient/strip|imageView2/2/w/320/format/webp)

这样，就可以通过不到10行的代码就可以优雅的解决群友的问题啦，不得不说Python以及pandas的强大了。