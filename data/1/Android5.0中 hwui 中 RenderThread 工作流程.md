## 前言

本篇文章是自己的一个学习笔记，记录了 Android 5.0 中 hwui 中的 RenderThread 的简单工作流程。由于是学习笔记，所以其中一些细节不会太详细，我只是将大概的流程走一遍，将其工作流标注出来，下次遇到问题的时候就可以知道去哪里查。

下图是我用 Systrace 抓取的一个应用启动的时候 RenderThread 的第一次 Draw 的 Trace 图，从这里面的顺序来看 RenderThread 的流程。熟悉应用启动流程的话应该知道，只有当第一次 DrawFrame 完成之后，整个应用的界面才会显示在手机上，在这之前，用户看到的是应用的 StartingWindow 的界面。
![img](https://upload-images.jianshu.io/upload_images/20166-9a9a23855a87489c.png?imageMogr2/auto-orient/strip|imageView2/2/w/638/format/webp)

RenderThread Draw first frame

## 从Java层说起

应用程序的每一帧是从接收到 VSYNC 信号开始进行计算和绘制的,这要从 Choreographer 这个类说起了，不过由于篇幅原因，我们直接看一帧的绘制调用关系链即可：
![img](https://upload-images.jianshu.io/upload_images/20166-3a6c09c1b5a07723.png?imageMogr2/auto-orient/strip|imageView2/2/w/937/format/webp)

绘制关系链

Choreographer 的 drawFrame 会调用到 ViewRootImpl 的 performTraversals 方法，而 performTraversals 方法最终会调用到performDraw() 方法， performDraw 又会调用到 draw(boolean fullRedrawNeeded) 方法，这个 draw 方法是 ViewRootImpl 的私有方法，和我们熟知的那个draw并不是同一个方法



```java
            if (mAttachInfo.mHardwareRenderer != null && mAttachInfo.mHardwareRenderer.isEnabled()) {
                mIsAnimating = false;
                boolean invalidateRoot = false;
                if (mHardwareYOffset != yOffset || mHardwareXOffset != xOffset) {
                    mHardwareYOffset = yOffset;
                    mHardwareXOffset = xOffset;
                    mAttachInfo.mHardwareRenderer.invalidateRoot();
                }
                mResizeAlpha = resizeAlpha;

                dirty.setEmpty();

                mBlockResizeBuffer = false;
                mAttachInfo.mHardwareRenderer.draw(mView, mAttachInfo, this);
}
```

如果是走硬件绘制路线的话，则会走这一条先，之后就会调用 mHardwareRenderer 的 draw 方法,这里的 mHardwareRenderer 指的是 ThreadedRenderer ，其 Draw 函数如下：



```java
    @Override
    void draw(View view, AttachInfo attachInfo, HardwareDrawCallbacks callbacks) {
        attachInfo.mIgnoreDirtyState = true;
        long frameTimeNanos = mChoreographer.getFrameTimeNanos();
        attachInfo.mDrawingTime = frameTimeNanos / TimeUtils.NANOS_PER_MS;

        long recordDuration = 0;
        if (mProfilingEnabled) {
            recordDuration = System.nanoTime();
        }

        updateRootDisplayList(view, callbacks);

        if (mProfilingEnabled) {
            recordDuration = System.nanoTime() - recordDuration;
        }

        attachInfo.mIgnoreDirtyState = false;

        // register animating rendernodes which started animating prior to renderer
        // creation, which is typical for animators started prior to first draw
        if (attachInfo.mPendingAnimatingRenderNodes != null) {
            final int count = attachInfo.mPendingAnimatingRenderNodes.size();
            for (int i = 0; i < count; i++) {
                registerAnimatingRenderNode(
                        attachInfo.mPendingAnimatingRenderNodes.get(i));
            }
            attachInfo.mPendingAnimatingRenderNodes.clear();
            // We don't need this anymore as subsequent calls to
            // ViewRootImpl#attachRenderNodeAnimator will go directly to us.
            attachInfo.mPendingAnimatingRenderNodes = null;
        }

        int syncResult = nSyncAndDrawFrame(mNativeProxy, frameTimeNanos,
                recordDuration, view.getResources().getDisplayMetrics().density);
        if ((syncResult & SYNC_INVALIDATE_REQUIRED) != 0) {
            attachInfo.mViewRootImpl.invalidate();
        }
    }
```

这个函数里面的 updateRootDisplayList(view, callbacks) ;即 getDisplayList 操作。接下来就是比较重要的一个操作：



```cpp
        int syncResult = nSyncAndDrawFrame(mNativeProxy, frameTimeNanos,
                recordDuration, view.getResources().getDisplayMetrics().density);
```

可以看出这是一个阻塞操作，等Native层完成后，拿到返回值后才会进行下一步的操作。

## Native层

其Native代码在android_view_ThreadedRenderer.cpp中，对应的实现代码如下：



```cpp
static int android_view_ThreadedRenderer_syncAndDrawFrame(JNIEnv* env, jobject clazz,
        jlong proxyPtr, jlong frameTimeNanos, jlong recordDuration, jfloat density) {
    RenderProxy* proxy = reinterpret_cast<RenderProxy*>(proxyPtr);
    return proxy->syncAndDrawFrame(frameTimeNanos, recordDuration, density);
}
```

RenderProxy的路径位于frameworks/base/libs/hwui/renderthread/RenderProxy.cpp



```cpp
int RenderProxy::syncAndDrawFrame(nsecs_t frameTimeNanos, nsecs_t recordDurationNanos,
        float density) {
    mDrawFrameTask.setDensity(density);
    return mDrawFrameTask.drawFrame(frameTimeNanos, recordDurationNanos);
}
```

其中 mDrawFrameTask 是一个 DrawFrameTask 对象，其路径位于frameworks/base/libs/hwui/renderthread/DrawFrameTask.cpp，其中drawFrame代码：



```cpp
int DrawFrameTask::drawFrame(nsecs_t frameTimeNanos, nsecs_t recordDurationNanos) {
    mSyncResult = kSync_OK;
    mFrameTimeNanos = frameTimeNanos;
    mRecordDurationNanos = recordDurationNanos;
    postAndWait();

    // Reset the single-frame data
    mFrameTimeNanos = 0;
    mRecordDurationNanos = 0;

    return mSyncResult;
}
```

其中 postAndWait() 的实现如下：



```cpp
void DrawFrameTask::postAndWait() {
    AutoMutex _lock(mLock);
    mRenderThread->queue(this);
    mSignal.wait(mLock);
}
```

就是将一个 DrawFrameTask 放入到了 mRenderThread 中,其中 queue 方法实现如下：



```cpp
void RenderThread::queue(RenderTask* task) {
    AutoMutex _lock(mLock);
    mQueue.queue(task);
    if (mNextWakeup && task->mRunAt < mNextWakeup) {
        mNextWakeup = 0;
        mLooper->wake();
    }
}
```

其中 mQueue 是一个 TaskQueue 对象，其



```php
void TaskQueue::queue(RenderTask* task) {
    // Since the RenderTask itself forms the linked list it is not allowed
    // to have the same task queued twice
    LOG_ALWAYS_FATAL_IF(task->mNext || mTail == task, "Task is already in the queue!");
    if (mTail) {
        // Fast path if we can just append
        if (mTail->mRunAt <= task->mRunAt) {
            mTail->mNext = task;
            mTail = task;
        } else {
            // Need to find the proper insertion point
            RenderTask* previous = 0;
            RenderTask* next = mHead;
            while (next && next->mRunAt <= task->mRunAt) {
                previous = next;
                next = next->mNext;
            }
            if (!previous) {
                task->mNext = mHead;
                mHead = task;
            } else {
                previous->mNext = task;
                if (next) {
                    task->mNext = next;
                } else {
                    mTail = task;
                }
            }
        }
    } else {
        mTail = mHead = task;
    }
}
```

接着看 RenderThread 之前的 queue 方法，



```cpp
void Looper::wake() {
    ssize_t nWrite;
    do {
        nWrite = write(mWakeWritePipeFd, "W", 1);
    } while (nWrite == -1 && errno == EINTR);

    if (nWrite != 1) {
        if (errno != EAGAIN) {
            ALOGW("Could not write wake signal, errno=%d", errno);
        }
    }
}
```

wake 函数则更为简单，仅仅向管道的写端写入一个字符“W”，这样管道的读端就会因为有数据可读而从等待状态中醒来。

## HWUI-RenderThread

接下来会到哪里去，我们首先要熟悉一下RenderThread，RenderThread是继承自Thread的，这个Thread是utils/Thread.h,RenderThread的初始化函数



```cpp
RenderThread::RenderThread() : Thread(true), Singleton<RenderThread>()
        , mNextWakeup(LLONG_MAX)
        , mDisplayEventReceiver(0)
        , mVsyncRequested(false)
        , mFrameCallbackTaskPending(false)
        , mFrameCallbackTask(0)
        , mRenderState(NULL)
        , mEglManager(NULL) {
    mFrameCallbackTask = new DispatchFrameCallbacks(this);
    mLooper = new Looper(false);
    run("RenderThread");
}
```

其 run 方法在 Thread 中有说明：



```cpp
    // Start the thread in threadLoop() which needs to be implemented.
    virtual status_t    run(    const char* name = 0,
                                int32_t priority = PRIORITY_DEFAULT,
                                size_t stack = 0);
```

即启动 threadLoop 函数，我们来看 RenderThread 的 threadLoop 函数，这个函数比较重要：



```cpp
bool RenderThread::threadLoop() {
#if defined(HAVE_PTHREADS)
    setpriority(PRIO_PROCESS, 0, PRIORITY_DISPLAY);
#endif
    initThreadLocals();

    int timeoutMillis = -1;
    for (;;) {
        int result = mLooper->pollOnce(timeoutMillis);
        LOG_ALWAYS_FATAL_IF(result == Looper::POLL_ERROR,
                "RenderThread Looper POLL_ERROR!");

        nsecs_t nextWakeup;
        // Process our queue, if we have anything
        while (RenderTask* task = nextTask(&nextWakeup)) {
            task->run();
            // task may have deleted itself, do not reference it again
        }
        if (nextWakeup == LLONG_MAX) {
            timeoutMillis = -1;
        } else {
            nsecs_t timeoutNanos = nextWakeup - systemTime(SYSTEM_TIME_MONOTONIC);
            timeoutMillis = nanoseconds_to_milliseconds(timeoutNanos);
            if (timeoutMillis < 0) {
                timeoutMillis = 0;
            }
        }

        if (mPendingRegistrationFrameCallbacks.size() && !mFrameCallbackTaskPending) {
            drainDisplayEventQueue(true);
            mFrameCallbacks.insert(
                    mPendingRegistrationFrameCallbacks.begin(), mPendingRegistrationFrameCallbacks.end());
            mPendingRegistrationFrameCallbacks.clear();
            requestVsync();
        }
    }

    return false;
}
```

可以看到，一个 for 循环是一个无限循环，而其中 pollOnce 是一个阻塞函数，直到我们上面调用了 mLooper->wake() 之后，会继续往下走，走到 while 循环中：



```cpp
while (RenderTask* task = nextTask(&nextWakeup)) {
            task->run();
            // task may have deleted itself, do not reference it again
        }
```

会将 RenderTask 取出来执行其 run 方法，经过前面的流程我们知道这个 RenderTask 是一个 DrawFrameTask ，其run方法如下：



```cpp
void DrawFrameTask::run() {
    ATRACE_NAME("DrawFrame");

    mContext->profiler().setDensity(mDensity);
    mContext->profiler().startFrame(mRecordDurationNanos);

    bool canUnblockUiThread;
    bool canDrawThisFrame;
    {
        TreeInfo info(TreeInfo::MODE_FULL, mRenderThread->renderState());
        canUnblockUiThread = syncFrameState(info);
        canDrawThisFrame = info.out.canDrawThisFrame;
    }

    // Grab a copy of everything we need
    CanvasContext* context = mContext;

    // From this point on anything in "this" is *UNSAFE TO ACCESS*
    if (canUnblockUiThread) {
        unblockUiThread();
    }

    if (CC_LIKELY(canDrawThisFrame)) {
        context->draw();
    }

    if (!canUnblockUiThread) {
        unblockUiThread();
    }
}
```

## RenderThread.DrawFrame

上面说到了 DrawFrameTask 的 run 方法，这里 run 方法中的执行的方法即我们在最前面那张图中所示的部分（即文章最前面那张图），下面的流程就是那张图中的函数调用，我们结合代码和图，一部分一部分来走整个 DrawFrame 的流程：

### 1. syncFrameState

第一个比较重要的函数是 syncFrameState ，从函数名就可以知道， syncFrameState 的作用就是同步 frame 信息，将 Java 层维护的 frame 信息同步到 RenderThread中。

> Main Thread 和Render Thread 都各自维护了一份应用程序窗口视图信息。各自维护了一份应用程序窗口视图信息的目的，就是为了可以互不干扰，进而实现最大程度的并行。其中，Render Thread维护的应用程序窗口视图信息是来自于 Main Thread 的。因此，当Main Thread 维护的应用程序窗口信息发生了变化时，就需要同步到 Render Thread 去。

所以查看代码就可以知道有两个 RenderNode，一个在 hwui 中，一个在 View 中。简单来说，同步信息就是将 Java 层的 RenderNode 中的信息同步到 hwui 中的 RenderNode 中。 注意syncFrameState的返回值赋给了 canUnblockUiThread ，从名字可以看出这个 canUnblockUiThread 的作用是判断是否唤醒 Main Thread ，也就是说如果返回为 true 的话，会提前唤醒主线程来执行其他的事情，而不用等到 draw 完成后再去唤醒 Main Thread。 这也是 Android 5.0 和 Android 4.x 最大的区别了。


![img](https://upload-images.jianshu.io/upload_images/20166-22aef1af1fcd1745.png?imageMogr2/auto-orient/strip|imageView2/2/w/478/format/webp)

syncFrameState



```rust
bool DrawFrameTask::syncFrameState(TreeInfo& info) {
    mRenderThread->timeLord().vsyncReceived(mFrameTimeNanos);
    mContext->makeCurrent();
    Caches::getInstance().textureCache.resetMarkInUse();

    for (size_t i = 0; i < mLayers.size(); i++) {
        mContext->processLayerUpdate(mLayers[i].get());
    }
    mLayers.clear();
    mContext->prepareTree(info);

    if (info.out.hasAnimations) {
        if (info.out.requiresUiRedraw) {
            mSyncResult |= kSync_UIRedrawRequired;
        }
    }
    // If prepareTextures is false, we ran out of texture cache space
    return info.prepareTextures;
}
```

首先是makeCurrent，这里的mContext是一个CanvasContext对象，其makeCurrent实现如下：



```cpp
void CanvasContext::makeCurrent() {
    // In the meantime this matches the behavior of GLRenderer, so it is not a regression
    mHaveNewSurface |= mEglManager.makeCurrent(mEglSurface);
}
```

mEglManager是一个EglManager对象，其实现为：



```cpp
bool EglManager::makeCurrent(EGLSurface surface) {
    if (isCurrent(surface)) return false;

    if (surface == EGL_NO_SURFACE) {
        // If we are setting EGL_NO_SURFACE we don't care about any of the potential
        // return errors, which would only happen if mEglDisplay had already been
        // destroyed in which case the current context is already NO_CONTEXT
        TIME_LOG("eglMakeCurrent", eglMakeCurrent(mEglDisplay, EGL_NO_SURFACE, EGL_NO_SURFACE, EGL_NO_CONTEXT));
    } else {
        EGLBoolean success;
        TIME_LOG("eglMakeCurrent", success = eglMakeCurrent(mEglDisplay, surface, surface, mEglContext));
        if (!success) {
            LOG_ALWAYS_FATAL("Failed to make current on surface %p, error=%s",
                (void*)surface, egl_error_str());
        }
    }
    mCurrentSurface = surface;
    return true;
}
```

这里会判断mCurrentSurface == surface，如果成立，则不用再初始化操作，如果是另外一个surface。，则会执行eglMakeCurrent，来重新创建上下文。

makeCurrent之后，会调用mContext->prepareTree(info)，其实现如下：



```cpp
void CanvasContext::prepareTree(TreeInfo& info) {
    mRenderThread.removeFrameCallback(this);

    info.damageAccumulator = &mDamageAccumulator;
    info.renderer = mCanvas;
    if (mPrefetechedLayers.size() && info.mode == TreeInfo::MODE_FULL) {
        info.canvasContext = this;
    }
    mAnimationContext->startFrame(info.mode);
    mRootRenderNode->prepareTree(info);
    mAnimationContext->runRemainingAnimations(info);

    if (info.canvasContext) {
        freePrefetechedLayers();
    }

    int runningBehind = 0;
    // TODO: This query is moderately expensive, investigate adding some sort
    // of fast-path based off when we last called eglSwapBuffers() as well as
    // last vsync time. Or something.
    TIME_LOG("nativeWindowQuery", mNativeWindow->query(mNativeWindow.get(),
            NATIVE_WINDOW_CONSUMER_RUNNING_BEHIND, &runningBehind));
    info.out.canDrawThisFrame = !runningBehind;

    if (info.out.hasAnimations || !info.out.canDrawThisFrame) {
        if (!info.out.requiresUiRedraw) {
            // If animationsNeedsRedraw is set don't bother posting for an RT anim
            // as we will just end up fighting the UI thread.
            mRenderThread.postFrameCallback(this);
        }
    }
}
```

其中 mRootRenderNode->prepareTree(info) 又是最重要的。回到Java层，我们知道 ThreadedRenderer 在初始化时，初始化了一个指针



```java
long rootNodePtr = nCreateRootRenderNode();
```

这个RootRenderNode也就是一个根Node，



```undefined
mRootNode = RenderNode.adopt(rootNodePtr);
```

然后会创建一个 mNativeProxy 指针，在 Native 层初始化一个 RenderProxy 对象，将 rootNodePtr 传给 RenderProxy 对象，这样在 RenderProxy 我们就可以得到这个对象的指针了。其中 CanvasContext 也是在 RenderProxy 对象初始化的时候被初始化的，初始化的时候将 rootNodePtr 传给了 CanvasContext 对象。

我们之前提到 ThreadedRenderer 的 draw 方法中首先会调用updateRootDisplayList，即我们熟悉的 getDisplayList 。这个方法中，其实也分为两个步骤，第一个步骤是 updateViewTreeDisplayList，第二个步骤是将根 Node 加入到 DrawOp 中：



```css
canvas.insertReorderBarrier();
canvas.drawRenderNode(view.getDisplayList());
canvas.insertInorderBarrier();
```

其最终实现在



```cpp
status_t DisplayListRenderer::drawRenderNode(RenderNode* renderNode, Rect& dirty, int32_t flags) {
    LOG_ALWAYS_FATAL_IF(!renderNode, "missing rendernode");

    // dirty is an out parameter and should not be recorded,
    // it matters only when replaying the display list
    DrawRenderNodeOp* op = new (alloc()) DrawRenderNodeOp(renderNode, flags, *currentTransform());
    addRenderNodeOp(op);

    return DrawGlInfo::kStatusDone;
}
```

再回到我们之前的 CanvasContext.prepareTree 中提到的 mRootRenderNode->prepareTree(info)，这时候这里的 mRootRenderNode 就是 CanvasContext 初始化是传进来的。

其实现在 RenderNode.cpp 中：



```cpp
void RenderNode::prepareTree(TreeInfo& info) {
    prepareTreeImpl(info);
}

void RenderNode::prepareTreeImpl(TreeInfo& info) {
    TT_START_MARK(getName());
    info.damageAccumulator->pushTransform(this);

    if (info.mode == TreeInfo::MODE_FULL) {
        pushStagingPropertiesChanges(info); //同步当前正在处理的Render Node的Property
    }
    uint32_t animatorDirtyMask = 0;
    if (CC_LIKELY(info.runAnimations)) {
        animatorDirtyMask = mAnimatorManager.animate(info);//执行动画相关的操作
    }
    prepareLayer(info, animatorDirtyMask);
    if (info.mode == TreeInfo::MODE_FULL) {
        pushStagingDisplayListChanges(info);  //同步当前正在处理的Render Node的Display List
    }
    prepareSubTree(info, mDisplayListData); //同步当前正在处理的Render Node的Display List引用的Bitmap，以及当前正在处理的Render Node的子Render Node的Display List等信息
    pushLayerUpdate(info); //检查当前正在处理的Render Node是否设置了Layer。如果设置了的话，就对这些Layer进行处理

    info.damageAccumulator->popTransform();
    TT_END_MARK();
}
```

这里所涉及到的进一步的具体操作大家可以自行去看代码。

## 2. draw

![img](https://upload-images.jianshu.io/upload_images/20166-9748a10652e3ee28.png?imageMogr2/auto-orient/strip|imageView2/2/w/540/format/webp)

Draw

执行完syncFrameState之后，接下来就是执行draw



```php
    if (CC_LIKELY(canDrawThisFrame)) {
        context->draw();
    }
```

CanvasContext的draw函数是一个核心函数，其位置在 frameworks/base/libs/hwui/OpenGLRenderer.cpp ，其实现如下：



```php
void CanvasContext::draw() {
    profiler().markPlaybackStart();

    SkRect dirty;
    mDamageAccumulator.finish(&dirty);

    ......

    status_t status;
    if (!dirty.isEmpty()) {
        status = mCanvas->prepareDirty(dirty.fLeft, dirty.fTop,
                dirty.fRight, dirty.fBottom, mOpaque);
    } else {
        status = mCanvas->prepare(mOpaque);
    }

    Rect outBounds;
    status |= mCanvas->drawRenderNode(mRootRenderNode.get(), outBounds);

    profiler().draw(mCanvas);

    mCanvas->finish();

    profiler().markPlaybackEnd();

    if (status & DrawGlInfo::kStatusDrew) {
        swapBuffers();
    }

    profiler().finishFrame();

    /// M: enable to get overdraw count
    if (CC_UNLIKELY(g_HWUI_debug_overdraw)) {
        if (!mDebugOverdrawLayer) {
            mDebugOverdrawLayer = LayerRenderer::createRenderLayer(mRenderThread.renderState(),
                mCanvas->getWidth(), mCanvas->getHeight());
        } else if (mDebugOverdrawLayer->layer.getWidth() != mCanvas->getWidth() ||
                   mDebugOverdrawLayer->layer.getHeight() != mCanvas->getHeight()) {
            if (!LayerRenderer::resizeLayer(mDebugOverdrawLayer, mCanvas->getWidth(), mCanvas->getHeight())) {
                LayerRenderer::destroyLayer(mDebugOverdrawLayer);
                mDebugOverdrawLayer = NULL;
            }
        }

    ......
}
```

#### 2.1 eglBeginFrame

首先来看eglBeginFrame的实现



```cpp
void EglManager::beginFrame(EGLSurface surface, EGLint* width, EGLint* height) {
    makeCurrent(surface);
    if (width) {
        eglQuerySurface(mEglDisplay, surface, EGL_WIDTH, width);
    }
    if (height) {
        eglQuerySurface(mEglDisplay, surface, EGL_HEIGHT, height);
    }
    eglBeginFrame(mEglDisplay, surface);
}
```

makeCurrent是用来管理上下文，eglBeginFrame主要是校验参数的合法性。

#### 2.2 prepareDirty



```php
    status_t status;
    if (!dirty.isEmpty()) {
        status = mCanvas->prepareDirty(dirty.fLeft, dirty.fTop,
                dirty.fRight, dirty.fBottom, mOpaque);
    } else {
        status = mCanvas->prepare(mOpaque);
    }
```

这里的mCanvas是一个OpenGLRenderer对象，其prepareDirty实现



```cpp
//TODO:增加函数功能描述
status_t OpenGLRenderer::prepareDirty(float left, float top,
        float right, float bottom, bool opaque) {
    setupFrameState(left, top, right, bottom, opaque);

    // Layer renderers will start the frame immediately
    // The framebuffer renderer will first defer the display list
    // for each layer and wait until the first drawing command
    // to start the frame
    if (currentSnapshot()->fbo == 0) {
        syncState();
        updateLayers();
    } else {
        return startFrame();
    }

    return DrawGlInfo::kStatusDone;
}
```

#### 2.3 drawRenderNode



```csharp
Rect outBounds;
status |= mCanvas->drawRenderNode(mRootRenderNode.get(), outBounds);
```

接下来就是调用OpenGLRenderer的drawRenderNode方法进行绘制



```cpp
status_t OpenGLRenderer::drawRenderNode(RenderNode* renderNode, Rect& dirty, int32_t replayFlags) {
    status_t status;
    // All the usual checks and setup operations (quickReject, setupDraw, etc.)
    // will be performed by the display list itself
    if (renderNode && renderNode->isRenderable()) {
        // compute 3d ordering
        renderNode->computeOrdering();
        if (CC_UNLIKELY(mCaches.drawDeferDisabled)) { //判断是否不重排序
            status = startFrame();
            ReplayStateStruct replayStruct(*this, dirty, replayFlags);
            renderNode->replay(replayStruct, 0);
            return status | replayStruct.mDrawGlStatus;
        }

        // 需要重新排序
        bool avoidOverdraw = !mCaches.debugOverdraw && !mCountOverdraw; // shh, don't tell devs!
        DeferredDisplayList deferredList(*currentClipRect(), avoidOverdraw);
        DeferStateStruct deferStruct(deferredList, *this, replayFlags);
        renderNode->defer(deferStruct, 0); //递归进行重排操作

        flushLayers(); // 首先执行设置了 Layer 的子 Render Node 的绘制命令，以便得到一个对应的FBO
        status = startFrame(); //执行一些诸如清理颜色绘冲区等基本操作
        status = deferredList.flush(*this, dirty) | status;
        return status;
    }

    // Even if there is no drawing command(Ex: invisible),
    // it still needs startFrame to clear buffer and start tiling.
    return startFrame();
}
```

这里的 renderNode 是一个 Root Render Node，

可以看到，到了这里虽然只是开始，但是其实已经结束了，这个函数里面最重要的几步:



```kotlin
renderNode->defer(deferStruct, 0); //进行重排序

flushLayers(); 首先执行设置了 Layer 的子 Render Node 的绘制命令，以便得到一个对应的FBO

status = deferredList.flush(*this, dirty) | status;   //对deferredList中的绘制命令进行真正的绘制操作
```

这几个是渲染部分真正的核心部分，其中的代码细节需要自己去研究。老罗在这部分讲的很细，有空可以去看看他的文章[Android应用程序UI硬件加速渲染的Display List渲染过程分析](https://link.jianshu.com?t=http://blog.csdn.net/luoshengyang/article/details/46281499).

#### 2.4 swapBuffers



```php
    if (status & DrawGlInfo::kStatusDrew) {
        swapBuffers();
    }
```

其核心就是调用EGL的 eglSwapBuffers(mEglDisplay, surface), duration)函数。

#### 2.5  FinishFrame



```css
    profiler().finishFrame();
```

主要是记录时间信息。

## 总结

鉴于我比较懒，而且总结能力不如老罗，就直接把他的总结贴过来了。
 RenderThread的总的流程如下：

> 1. 将Main Thread维护的Display List同步到Render Thread维护的Display List去。这个同步过程由Render Thread执行，但是Main Thread会被阻塞住。

> 1. 如果能够完全地将Main Thread维护的Display List同步到Render Thread维护的Display List去，那么Main Thread就会被唤醒，此后Main Thread和Render Thread就互不干扰，各自操作各自内部维护的Display List；否则的话，Main Thread就会继续阻塞，直到Render Thread完成应用程序窗口当前帧的渲染为止。

> 1. Render Thread在渲染应用程序窗口的Root Render Node的Display List之前，首先将那些设置了Layer的子Render Node的Display List渲染在各自的一个FBO上，接下来再一起将这些FBO以及那些没有设置Layer的子Render Node的Display List一起渲染在Frame Buffer之上，也就是渲染在从Surface Flinger请求回来的一个图形缓冲区上。这个图形缓冲区最终会被提交给Surface Flinger合并以及显示在屏幕上。

> 第2步能够完全将Main Thread维护的Display List同步到Render Thread维护的Display List去很关键，它使得Main Thread和Render Thread可以并行执行，这意味着Render Thread在渲染应用程序窗口当前帧的Display List的同时，Main Thread可以去准备应用程序窗口下一帧的Display List，这样就使得应用程序窗口的UI更流畅。

注意最后一段，在 Android 4.x 时代，没有RenderThread的时代，只有 Main Thread ，也就是说 必须要等到 Draw 完成后，才会去准备下一帧的数据。
