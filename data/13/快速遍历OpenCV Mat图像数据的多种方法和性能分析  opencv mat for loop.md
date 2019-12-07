[![img](https://upload.jianshu.io/users/upload_avatars/20059030/60ad3b7b-5caf-4ff7-890c-6983c6385619.jpg?imageMogr2/auto-orient/strip|imageView2/1/w/96/h/96/format/webp)](https://www.jianshu.com/u/12461cd06bc9)

[kezunlin](https://www.jianshu.com/u/12461cd06bc9)关注

0.3452019.11.19 08:48:23字数 343阅读 8

**本文首发于个人博客https://kezunlin.me/post/61d55ab4/，欢迎阅读!**

opencv mat for loop

# Series

- [Part 1: compile opencv on ubuntu 16.04](https://links.jianshu.com/go?to=https%3A%2F%2Fkezunlin.me%2Fpost%2F15f5c3e8%2F)
- [Part 2: compile opencv with CUDA support on windows 10](https://links.jianshu.com/go?to=https%3A%2F%2Fkezunlin.me%2Fpost%2F6580691f%2F)
- **Part 3: opencv mat for loop**
- [Part 4: speed up opencv image processing with openmp](https://links.jianshu.com/go?to=https%3A%2F%2Fkezunlin.me%2Fpost%2F7a6ba82e%2F)

# Guide

## Mat

- for gray image, use type `<uchar>`
- for RGB color image，use type `<Vec3b>`

gray format storage





![img](https://upload-images.jianshu.io/upload_images/20059030-335851b487857f6b.png?imageMogr2/auto-orient/strip|imageView2/2/w/446/format/webp)

gray

color format storage: BGR





![img](https://upload-images.jianshu.io/upload_images/20059030-6b0e59c74fa978d6.png?imageMogr2/auto-orient/strip|imageView2/2/w/704/format/webp)

BGR

> we can use method `isContinuous()` to judge whether the memory buffer is continuous or not.

## color space reduction



```cpp
uchar color_space_reduction(uchar pixel)
{
    /*
    0-9 ===>0
    10-19===>10
    20-29===>20
    ...
    240-249===>24
    250-255===>25

    map from 256*256*256===>26*26*26
    */

    int divideWith = 10;
    uchar new_pixel = (pixel / divideWith)*divideWith;
    return new_pixel;
}
```

## color table



```cpp
void get_color_table()
{
    // cache color value in table[256]
    int divideWith = 10;
    uchar table[256];
    for (int i = 0; i < 256; ++i)
        table[i] = divideWith* (i / divideWith);
}
```

# C++

## ptr []



```cpp
// C ptr []: faster but not safe
Mat& ScanImageAndReduce_Cptr(Mat& I, const uchar* const table)
{
    // accept only char type matrices
    CV_Assert(I.depth() != sizeof(uchar));
    int channels = I.channels();
    int nRows = I.rows;
    int nCols = I.cols* channels;
    if (I.isContinuous())
    {
        nCols *= nRows;
        nRows = 1;
    }
    int i, j;
    uchar* p;
    for (i = 0; i < nRows; ++i)
    {
        p = I.ptr<uchar>(i);
        for (j = 0; j < nCols; ++j)
        {
            p[j] = table[p[j]];
        }
    }
    return I;
}
```

## ptr ++



```cpp
// C ptr ++: faster but not safe
Mat& ScanImageAndReduce_Cptr2(Mat& I, const uchar* const table)
{
    // accept only char type matrices
    CV_Assert(I.depth() != sizeof(uchar));
    int channels = I.channels();
    int nRows = I.rows;
    int nCols = I.cols* channels;
    if (I.isContinuous())
    {
        nCols *= nRows;
        nRows = 1;
    }
    uchar* start = I.ptr<uchar>(0); // same as I.ptr<uchar>(0,0)
    uchar* end = start + nRows * nCols;
    for (uchar* p=start; p < end; ++p)
    {
        *p = table[*p];
    }
    return I;
}
```

## at<uchar>(i,j)



```cpp
// at<uchar>(i,j): random access, slow
Mat& ScanImageAndReduce_atRandomAccess(Mat& I, const uchar* const table)
{
   // accept only char type matrices
   CV_Assert(I.depth() != sizeof(uchar));
   const int channels = I.channels();
   switch (channels)
   {
   case 1:
   {
       for (int i = 0; i < I.rows; ++i)
           for (int j = 0; j < I.cols; ++j)
               I.at<uchar>(i, j) = table[I.at<uchar>(i, j)];
       break;
   }
   case 3:
   {
       Mat_<Vec3b> _I = I;

       for (int i = 0; i < I.rows; ++i)
           for (int j = 0; j < I.cols; ++j)
           {
               _I(i, j)[0] = table[_I(i, j)[0]];
               _I(i, j)[1] = table[_I(i, j)[1]];
               _I(i, j)[2] = table[_I(i, j)[2]];
           }
       I = _I;
       break;
   }
   }
   return I;
}
```

## Iterator



```cpp
// MatIterator_<uchar>: safe but slow
Mat& ScanImageAndReduce_Iterator(Mat& I, const uchar* const table)
{
   // accept only char type matrices
   CV_Assert(I.depth() != sizeof(uchar));
   const int channels = I.channels();
   switch (channels)
   {
   case 1:
   {
       MatIterator_<uchar> it, end;
       for (it = I.begin<uchar>(), end = I.end<uchar>(); it != end; ++it)
           *it = table[*it];
       break;
   }
   case 3:
   {
       MatIterator_<Vec3b> it, end;
       for (it = I.begin<Vec3b>(), end = I.end<Vec3b>(); it != end; ++it)
       {
           (*it)[0] = table[(*it)[0]];
           (*it)[1] = table[(*it)[1]];
           (*it)[2] = table[(*it)[2]];
       }
   }
   }
   return I;
}
```

## opencv LUT



```cpp
// LUT
Mat& ScanImageAndReduce_LUT(Mat& I, const uchar* const table)
{
   Mat lookUpTable(1, 256, CV_8U);
   uchar* p = lookUpTable.data;
   for (int i = 0; i < 256; ++i)
       p[i] = table[i];

   cv::LUT(I, lookUpTable, I);
   return I;
}
```

## forEach

> `forEach` method of the `Mat` class that utilizes all the cores on your machine to apply any function at every pixel.



```cpp
// Parallel execution with function object.
struct ForEachOperator
{
    uchar m_table[256];
    ForEachOperator(const uchar* const table)
    {
        for (size_t i = 0; i < 256; i++)
        {
            m_table[i] = table[i];
        }
    }

    void operator ()(uchar& p, const int * position) const
    {
        // Perform a simple operation
        p = m_table[p];
    }
};

// forEach use multiple processors, very fast
Mat& ScanImageAndReduce_forEach(Mat& I, const uchar* const table)
{
    I.forEach<uchar>(ForEachOperator(table));
    return I;
}
```

## forEach with lambda



```cpp
// forEach lambda use multiple processors, very fast (lambda slower than ForEachOperator)
Mat& ScanImageAndReduce_forEach_with_lambda(Mat& I, const uchar* const table)
{
    I.forEach<uchar>
    (
        [=](uchar &p, const int * position) -> void
        {
            p = table[p];
        }
    );
    return I;
}
```

## time cost

### no foreach



```bash
[1 Cptr   ] times=5000, total_cost=988 ms, avg_cost=0.1976 ms
[1 Cptr2  ] times=5000, total_cost=1704 ms, avg_cost=0.3408 ms
[2 atRandom] times=5000, total_cost=9611 ms, avg_cost=1.9222 ms
[3 Iterator] times=5000, total_cost=20195 ms, avg_cost=4.039 ms
[4 LUT    ] times=5000, total_cost=899 ms, avg_cost=0.1798 ms

[1 Cptr   ] times=10000, total_cost=2425 ms, avg_cost=0.2425 ms
[1 Cptr2  ] times=10000, total_cost=3391 ms, avg_cost=0.3391 ms
[2 atRandom] times=10000, total_cost=20024 ms, avg_cost=2.0024 ms
[3 Iterator] times=10000, total_cost=39980 ms, avg_cost=3.998 ms
[4 LUT    ] times=10000, total_cost=103 ms, avg_cost=0.0103 ms
```

### foreach



```bash
[5 forEach     ] times=200000, total_cost=199 ms, avg_cost=0.000995 ms
[5 forEach lambda] times=200000, total_cost=521 ms, avg_cost=0.002605 ms

[5 forEach     ] times=20000, total_cost=17 ms, avg_cost=0.00085 ms
[5 forEach lambda] times=20000, total_cost=23 ms, avg_cost=0.00115 ms
```

### results

Loop Type | Time Cost (us)
:----: |
ptr [] | 242
ptr ++ | 339
at<uchar> | 2002
iterator | 3998
LUT | 10
forEach | 0.85
forEach lambda | 1.15

`forEach` is 10x times faster than `LUT`, 240~340x times faster than `ptr []` and `ptr ++`, and 2000~4000x times faster than `at` and `iterator`.

## code

[code here](https://links.jianshu.com/go?to=https%3A%2F%2Fgist.github.com%2Fkezunlin%2F8a8f1be7c0e101ce3f0e16e529288afc)

# Python

## pure python



```python
# import the necessary packages
import matplotlib.pyplot as plt
import cv2
print(cv2.__version__)

%matplotlib inline
```



```css
3.4.2
```



```python
# load the original image, convert it to grayscale, and display
# it inline
image = cv2.imread("cat.jpg")
image = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
print(image.shape)
#plt.imshow(image, cmap="gray")
```



```undefined
(360, 480)
```



```python
%load_ext cython
```



```dart
The cython extension is already loaded. To reload it, use:
  %reload_ext cython
```



```cython
%%cython -a
 
def threshold_python(T, image):
    # grab the image dimensions
    h = image.shape[0]
    w = image.shape[1]
    
    # loop over the image, pixel by pixel
    for y in range(0, h):
        for x in range(0, w):
            # threshold the pixel
            image[y, x] = 255 if image[y, x] >= T else 0
            
    # return the thresholded image
    return image
```



```python
%timeit threshold_python(5, image)
```



```rust
263 ms ± 20.2 ms per loop (mean ± std. dev. of 7 runs, 1 loop each)
```

## cython



```cython
%%cython -a
 
import cython
 
@cython.boundscheck(False)
cpdef unsigned char[:, :] threshold_cython(int T, unsigned char [:, :] image):
    # set the variable extension types
    cdef int x, y, w, h
    
    # grab the image dimensions
    h = image.shape[0]
    w = image.shape[1]
    
    # loop over the image
    for y in range(0, h):
        for x in range(0, w):
            # threshold the pixel
            image[y, x] = 255 if image[y, x] >= T else 0
    
    # return the thresholded image
    return image
```

## numba



```python
%timeit threshold_cython(5, image)
```



```cpp
150 µs ± 7.14 µs per loop (mean ± std. dev. of 7 runs, 10000 loops each)
```



```python
from numba import njit

@njit
def threshold_njit(T, image):
    # grab the image dimensions
    h = image.shape[0]
    w = image.shape[1]
    
    # loop over the image, pixel by pixel
    for y in range(0, h):
        for x in range(0, w):
            # threshold the pixel
            image[y, x] = 255 if image[y, x] >= T else 0
            
    # return the thresholded image
    return image
```



```python
%timeit threshold_njit(5, image)
```



```cpp
43.5 µs ± 142 ns per loop (mean ± std. dev. of 7 runs, 10000 loops each)
```

## numpy



```python
def threshold_numpy(T, image):
    image[image > T] = 255
    return image
```



```python
%timeit threshold_numpy(5, image)
```



```cpp
111 µs ± 334 ns per loop (mean ± std. dev. of 7 runs, 10000 loops each)
```

## conclusions



```python
image = cv2.imread("cat.jpg")
image = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
print(image.shape)

%timeit threshold_python(5, image)
%timeit threshold_cython(5, image)
%timeit threshold_njit(5, image)
%timeit threshold_numpy(5, image)
```



```rust
(360, 480)
251 ms ± 6.5 ms per loop (mean ± std. dev. of 7 runs, 1 loop each)
143 µs ± 1.19 µs per loop (mean ± std. dev. of 7 runs, 10000 loops each)
43.8 µs ± 284 ns per loop (mean ± std. dev. of 7 runs, 10000 loops each)
113 µs ± 957 ns per loop (mean ± std. dev. of 7 runs, 10000 loops each)
```



```python
image = cv2.imread("big.jpg")
image = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
print(image.shape)

%timeit threshold_python(5, image)
%timeit threshold_cython(5, image)
%timeit threshold_njit(5, image)
%timeit threshold_numpy(5, image)
```



```css
(2880, 5120)
21.8 s ± 460 ms per loop (mean ± std. dev. of 7 runs, 1 loop each)
12.3 ms ± 231 µs per loop (mean ± std. dev. of 7 runs, 100 loops each)
3.91 ms ± 66.1 µs per loop (mean ± std. dev. of 7 runs, 100 loops each)
10.3 ms ± 179 µs per loop (mean ± std. dev. of 7 runs, 100 loops each)
```

60,480

- python: 251 ms
- cython: 143 us
- numba: 43 us
- numpy: 113 us

2880, 5120

- python: 21 s
- cython: 12 ms
- numba: 4 ms
- numpy: 10 ms

# Reference

- [Part1: OpenCV访问Mat图像中每个像素的值 4种对比](https://links.jianshu.com/go?to=https%3A%2F%2Fblog.csdn.net%2Fxiaowei_cqu%2Farticle%2Fdetails%2F7771760)
- [Part2: OpenCV访问Mat图像中每个像素的值 13种对比](https://links.jianshu.com/go?to=https%3A%2F%2Fblog.csdn.net%2Fxiaowei_cqu%2Farticle%2Fdetails%2F19839019)
- [parallel-pixel-access-in-opencv-using-foreach](https://links.jianshu.com/go?to=https%3A%2F%2Fwww.learnopencv.com%2Fparallel-pixel-access-in-opencv-using-foreach%2F)
- [fast-optimized-for-pixel-loops-with-opencv-and-python](https://links.jianshu.com/go?to=https%3A%2F%2Fwww.pyimagesearch.com%2F2017%2F08%2F28%2Ffast-optimized-for-pixel-loops-with-opencv-and-python%2F)
- [python performance tips](https://links.jianshu.com/go?to=https%3A%2F%2Fwiki.python.org%2Fmoin%2FPythonSpeed%2FPerformanceTips)

# History

- 20180823: created.

# Copyright

- Post author: [kezunlin](https://links.jianshu.com/go?to=https%3A%2F%2Fkezunlin.me)
- Post link: [https://kezunlin.me/post/61d55ab4/](https://links.jianshu.com/go?to=https%3A%2F%2Fkezunlin.me%2Fpost%2F61d55ab4%2F)
- Copyright Notice: All articles in this blog are licensed under CC BY-NC-SA 3.0 unless stating additionally.





2人点赞



[kezunlin.me](https://www.jianshu.com/nb/40683949)





"❤️随心❤️"