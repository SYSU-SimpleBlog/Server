# (Swift)设置UIcollectionView的section底色

[![img](https://cdn2.jianshu.io/assets/default_avatar/2-9636b13945b9ccf345bc98d0d81074eb.jpg)](https://www.jianshu.com/u/1f837584607b)

[kingjiajie](https://www.jianshu.com/u/1f837584607b)关注

0.0082019.11.24 18:16:20字数 571阅读 12

# 前言

具体代码demo如下：

GitHub_OC版本:[Demo具体代码](https://links.jianshu.com/go?to=https%3A%2F%2Fgithub.com%2Fkingjiajie%2FJJCollectionViewRoundFlowLayout)

GitHub_Swift版本:[Demo具体代码](https://links.jianshu.com/go?to=https%3A%2F%2Fgithub.com%2Fkingjiajie%2FJJCollectionViewRoundFlowLayout_Swift)

码云_OC:[Demo具体代码](https://links.jianshu.com/go?to=https%3A%2F%2Fgitee.com%2Fkingjiajie%2FJJCollectionViewRoundFlowLayout)

>   简单设计collectionview 底色和根据section不同设置不同颜色，支持collection横竖样式、自定义偏移量、投影。
>   由于APP设计样式的多样性，很多时候我们需要用到一些特别的样式，例如投影、圆角、某个空间下增加底色和投影等组合，这些看似很简单的样式，其实也需要花不少时间进行样式的布局和调整等。
>   例如本人遇到需要在collectionView，根据section不同设置不同的底色，需要动态设置是否包含headerview，还需要设置投影等等，所以设计了这个可配置且动态更新的 collection 背景颜色 控件。可基本满足各种要求。

# 设计思路

> 1、继承UICollectionViewFlowLayout，重写prepareLayout方法，在方法内计算每个section的大小，并根据用户设置的sectiooninset，进行frame补充。
> 2、继承UICollectionViewLayoutAttributes，增加底色、投影等参数。
> 3、在prepareLayout计算每个section的UICollectionViewLayoutAttributes并设置底色参数，并进行保存，
> 4、在layoutAttributesForElementsInRect进行rect判断获取attr。
> 5、在applyLayoutAttributes内机进行样式设置。

# 效果图：



![img](https://upload-images.jianshu.io/upload_images/7386003-41327157c45ff4be.gif?imageMogr2/auto-orient/strip|imageView2/2/w/300/format/webp)

image

# 支持类型：

> 1、collectionView section底色。
> 2、是否包含headerview。
> 3、是否包含footerview。
> 4、支持borderWidth、borderColor。
> 5、支持shadow投影。
> 6、支持collectionView，Vertical，Horizontal。
> 7、支持根据不同section分别设置不同底色显示。

# 核心代码



```swift
/// 计算默认不包含headerview和footerview的背景大小

/// @paramframeframe description
/// @paramsectionInsetsectionInset description
//MARK: 默认section无偏移大小
extension JJCollectionViewRoundFlowLayout_Swift{
    func calculateDefaultFrameWithSectionFrame(_ frame:CGRect ,sectionInset:UIEdgeInsets) -> CGRect{
        var sectionFrame = frame;
        sectionFrame.origin.x -= sectionInset.left;
        sectionFrame.origin.y -= sectionInset.top;
        if (self.scrollDirection == UICollectionView.ScrollDirection.horizontal) {
            sectionFrame.size.width += sectionInset.left + sectionInset.right;
            //减去系统adjustInset的top
            if #available(iOS 11, *) {
                sectionFrame.size.height = self.collectionView!.frame.size.height - self.collectionView!.adjustedContentInset.top;
            } else {
                sectionFrame.size.height = self.collectionView!.frame.size.height - abs(self.collectionView!.contentOffset.y)/*适配iOS11以下*/;
            }
        }else{
            sectionFrame.size.width = self.collectionView!.frame.size.width;
            sectionFrame.size.height += sectionInset.top + sectionInset.bottom;
        }
        return sectionFrame;
    }
}

override func layoutAttributesForElements(in rect: CGRect) -> [UICollectionViewLayoutAttributes]? {
    var attrs = super.layoutAttributesForElements(in: rect) ?? []
    for attr in self.decorationViewAttrs {
        attrs.append(attr)
    }
    return attrs
}

  override public func prepare() 代码有点多，就不贴出来了。下面有demo。
```

# 如何使用：

> pod 'JJCollectionViewRoundFlowLayout_Swift'



```swift
//可选设置
open var isCalculateHeader : Bool = false    // 是否计算header
open var isCalculateFooter : Bool = false    // 是否计算footer
```



```swift
/// 设置底色偏移量(该设置只设置底色，与collectionview原sectioninsets区分）
/// - Parameter collectionView: collectionView description
/// - Parameter collectionViewLayout: collectionViewLayout description
/// - Parameter section: section description
func collectionView(_ collectionView : UICollectionView , layout collectionViewLayout:UICollectionViewLayout,borderEdgeInsertsForSectionAtIndex section : Int) -> UIEdgeInsets;

/// 设置底色相关
/// - Parameter collectionView: collectionView description
/// - Parameter collectionViewLayout: collectionViewLayout description
/// - Parameter section: section description
func collectionView(_ collectionView : UICollectionView, layout collectionViewLayout : UICollectionViewLayout , configModelForSectionAtIndex section : Int ) -> JJCollectionViewRoundConfigModel_Swift;
```

在collectionview页面代码上加入代理（JJCollectionViewDelegateRoundFlowLayout）

# 并实现如下两个方法：



```swift
#pragma mark - JJCollectionViewDelegateRoundFlowLayout

func collectionView(_ collectionView: UICollectionView, layout collectionViewLayout: UICollectionViewLayout, borderEdgeInsertsForSectionAtIndex section: Int) -> UIEdgeInsets {
    return UIEdgeInsets.init(top: 5, left: 12, bottom: 5, right: 12)
}

func collectionView(_ collectionView: UICollectionView, layout collectionViewLayout: UICollectionViewLayout, configModelForSectionAtIndex section: Int) -> JJCollectionViewRoundConfigModel_Swift {
    let model = JJCollectionViewRoundConfigModel_Swift.init();
    
    model.backgroundColor = UIColor.init(red: 233/255.0, green:233/255.0 ,blue:233/255.0,alpha:1.0)
    model.cornerRadius = 10;
    return model;
}
```

# 效果如下：



![img](https://upload-images.jianshu.io/upload_images/7386003-46d60dd6c7c93d4d.png?imageMogr2/auto-orient/strip|imageView2/2/w/300/format/webp)

image

具体代码demo如下：

GitHub_OC版本:[Demo具体代码](https://links.jianshu.com/go?to=https%3A%2F%2Fgithub.com%2Fkingjiajie%2FJJCollectionViewRoundFlowLayout)

GitHub_Swift版本:[Demo具体代码](https://links.jianshu.com/go?to=https%3A%2F%2Fgithub.com%2Fkingjiajie%2FJJCollectionViewRoundFlowLayout_Swift)

码云_OC:[Demo具体代码](https://links.jianshu.com/go?to=https%3A%2F%2Fgitee.com%2Fkingjiajie%2FJJCollectionViewRoundFlowLayout) 大家有空可star。

后续可能会单独更新swift版本，敬请期待。

如有问题，可以直接提issues，或者发送邮件到[kingjiajie@sina.com](https://links.jianshu.com/go?to=mailto%3Akingjiajie%40sina.com)，或者直接回复。谢谢。