#### 目录 

**生产环境Redis中的热点key如何发现并优化？**

**有哪些常见的核心的MySQL应用开发规范?**

**高可用架构MHA有什么样的不足和风险点吗？**

**为什么pt-osc操作表以后中文注释显示???，如何避免？**

**MySQL 5.6升级5.7都有什么注意事项**

**在用阿里云、腾讯云等公有云时，你是如何评估新建主机/数据库对象的配置级别？**

**ALTER TABLE 出现duplicate primary xxx报错的原因及处理？**

**InnoDB在什么情况下会触发检查点(checkpoint)?**

**update t set a=29 and b in (1,2,3,4);这样写有什么问题吗？**

------





![img](https://upload-images.jianshu.io/upload_images/20009699-e06917d8b10cc07d.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/900/format/webp)

**一、生产环境Redis中的热点key如何发现并优化？**

1、用户消费的数据远大于生产的数据（热卖商品、热点新闻、热点评论、明星直播）

2、请求分片集中，超过单Server的性能极限。

**热点key可能造成如下问题：**

1、流量集中，达到物理网卡上限

2、请求过多，缓存分片服务被打垮

3、DB击穿，引起业务雪崩

**如何发现热点key：**

1、通过tcpdump抓包，可以分析抓取到网络包分析key的频率

2、redis客户端抓取，例如请求key的时候记录日志，分析日志得到key的访问频率

3、redis的monitor可以记录redis的所有操作，记录并分析monitor日志得到key的访问频率（注意：monitor可能会造成性能问题，慎重使用）

4、Redis 4.0提供了—hot-keys配合maxmemory-policy可以统计热点key

5、第三方开源项目，如facebook开源项目redis-faina，原理同3

**解决方案如下：**

1、对于”get”类型的热点key，通常可以为redis添加slave，通过slave承担读压力来缓解

2、服务端本地缓存，服务端先请求本地缓存，缓解redis压力

3、多级缓存方案，通过多级请求，层层过滤解决热点key问题

4、proxy方案，有些方案会探测分片热点key，缓存在proxy上缓解redis压力

5、同解决big方案类似，将一个key通过hash分解为多个key，value值一样，将这些key分散到集群多个分片中，需要访问时先根据hash算出对应的key，然后访问的具体的分片



**二、有哪些常见的核心的MySQL应用开发规范?**

这里重点介绍下Schema设计规范，其他规范请参考末尾《知数堂开发规范》

**（一）schema设计原则：**

1、尽量小的原则，不浪费

2、为了高并发，禁止使用外键

3、每个表必须有主键

4、字符集和库级保持一致，不单独定义字段字符集

**（二）字段规范：**

1、每个表建议不超过30-50个字段

2、优先选择utf8mb4字符集，它的兼容性最好，而且还支持emoji字符。如果对存储容量比较敏感的，可以改成latin1字符集

3、严禁在数据库中明文存储用户密码、身份证、信用卡号（信用卡PIN码）等核心机密数据，务必先行加密

4、存储整型数据时，默认加上UNSIGNED，扩大存储范围

5、建议用INT UNSIGNED存储IPV4地址，查询时再利用INET_ATON()、INET_NTOA()函数转换

6、如果遇到BLOB、TEXT字段，则尽量拆出去，再用主键做关联

7、在够用的前提下，选择尽可能小的字段，用于节省磁盘和内存空间

8、涉及精确金额相关用途时，建议扩大N倍后，全部转成整型存储（例如把分扩大百倍），避免浮点数加减出现不准确问题

**（三）常用数据类型参考：**

1、字符类型建议采用varchar数据类型（InnoDB建议用varchar替代char）

2、金额货币科学计数建议采用decimal数据类型，如果运算在数据库中完成可以考虑使用bigint存储，单位：分

3、自增长标识建议采用int或bigint数据类型，如果该表有大量的删除及再写入就使用bigint,反之int就够用

4、时间类型建议采用为datetime/timestamp数据类型

5、禁止使用text、longtext等的数据类型

6、字段值如果为非负数，就加上unsigned定语，提升可用范围

**（四）SQL规范**

1、在MySQL中SQL语句一般不区分大小写，全部小写

2、sql语句在使用join, 子查询一定先要进行explain确定执行计划

3、为每个业务收集sql list.

**知数堂开发规范：**[https://github.com/zhishutech...](https://links.jianshu.com/go?to=https%3A%2F%2Flink.zhihu.com%2F%3Ftarget%3Dhttps%3A%2F%2Fgithub.com%2Fzhishutech%2Fmysql-sql-standard%2Fblob%2Fmaster%2FSUMMARY.md)



**三、高可用架构MHA有什么样的不足和风险点吗？**

MHA作为传统复制下的高可用霸主，在今天的GTID环境下，开始慢慢走向没落，更多的人开始开始选择replication-manager或者orchestrator等高可用解决方案

**不足及风险点：**

1、failover依赖于外部脚本，比如VIP切换需要自己编写脚本实现

2、MHA启动后只检测主库是否正常，并不检查从库状态及主从延迟

3、需要基于SSH免认证配置，存在一定的安全隐患

4、没有提供从服务器的读负载均衡功能

5、从节点出现宕机等异常并没有能力处理，即没有从库故障转移能力

6、在高可用切换期间，某些场景下可能出现数据丢失的情况，并不保证数据0丢失

7、无法控制RTO恢复时间

**具体的数据丢失场景移步吴老师公开课《把MHA拉下神坛》**

[https://ke.qq.com/course/430673?tuin=2ce85033](https://links.jianshu.com/go?to=https%3A%2F%2Flink.zhihu.com%2F%3Ftarget%3Dhttps%3A%2F%2Fke.qq.com%2Fcourse%2F430673%3Ftuin%3D2ce85033)



**四、为什么pt-osc操作表以后中文注释显示???，如何避免？**

一般来说，生产环境使用的表都会使用中文注释表信息以及字段信息，但是如果使用pt-osc且未指定字符类型的情况下进行在线变更后，中文注释都会变成"???"，虽然不影响正常使用，但是对于认为阅读起来会造成困扰，某些平台会依据注释生成数据字典，因此正确的姿势是在使用pt-osc工具时通过--charset=utf8指定utf8字符集

**示例：**

**pt-online-schema-change -h 127.0.0.1 -u xxx -p xxx --alter="add index idx_id(id)" --chunk-size=5000 \**

**--print --no-version-check --execute D=xucl,t=test --charset=utf8**



**五、MySQL 5.6升级5.7都有什么注意事项**

**（一）MySQL升级的方式一般来说有两种**

1、通过inplace方式原地升级，升级系统表

2、通过新建实例，高版本作为低版本的从库进行滚动升级

MySQL5.7版本做了非常多的改变，升级5.6到5.7时需要考虑兼容性，避免升级到5.7之后因为种种参数设置不正确导致业务受影响，建议首先逐一查看release note

**(二)需要注意的参数及问题：**

1、sql_mode：MySQL 5.7采用严格模式，例如ONLY_FULL_GROUP_BY等

2、innodb_status_output_locks：MySQL 5.7支持将死锁信息打印到error log（其实这个参数MySQL 5.6就已支持）

3、innodb_page_cleaners：MySQL 5.7将脏页刷新线程从master线程独立出来了，对应参数为innodb_page_cleaners

4、innodb_strict_mode：控制CREATE TABLE, ALTER TABLE, CREATE INDEX, 和 OPTIMIZE TABLE的语法问题

5、show_compatibility_56=ON：控制show变量及状态信息输出，如果未开启show status 命令无法获取Slave_xxx 的状态

6、log_timestamps：控制error log/slow_log/genera log日志的显示时间，该参数可以设置为:UTC 和 SYSTEM，但是默认使用 UTC

7、disable_partition_engine_check：在表多的情况下可能导致启动非常慢

8、range_optimizer_max_mem_size：范围查询优化参数，这个参数限制范围查询优化使用的内存，默认8M

9、MySQL 5.7新增优化器选项derived_merge=on，可能导致SQL全表扫描，而在MySQL 5.6下可能表现为auto key

10、innodb_undo_directory && innodb_undo_logs：MySQL 5.7支持将undo从ibdata1独立出来（只支持实例初始化，不支持在线变更）

11、主从复制问题：MySQL5.7到小于5.6.22的复制存在bug(bug 74683)

12、SQL兼容性问题：SQL在MySQL 5.7和MySQL 5.6环境下结果可能不一致，因此建议获取线上SQL，在同样数据的环境下，在两个实例运行获取到的结果计算hash，比较hash值做兼容性判断

**（三）友情提醒**

1、升级前一定要做好备份！！！

2、升级正式环境前提前在测试环境进行仔细测试，确认无误以后再升级正式环境

3、做好相应的回退方案



**六、在用阿里云、腾讯云等公有云时，你是如何评估新建主机/数据库对象的配置级别？**

这里以云下业务迁移云上为例来探讨

1、首先熟悉现有业务的基本架构，比如一主多从、sharding架构等，并且知道相应的业务分布

2、获取现有业务的监控获取到的峰值QPS、TPS、IOPS、CPU使用率、磁盘使用量、内存使用量、最大连接数等关键指标

3、获取现有数据库的关键参数指标，如innodb_buffer_pool_size等

4、公有云每个规格都提供了相应的参数指标，如：核数、内存、IOPS、最大连接数等指标

5、根据第2、3、4步选择相应规格的RDS，原则为RDS规格参数大于现有环境状态指标，其中IOPS需要进行换算(云上的IOPS一般按4k算，而自建的一般按16k算)

6、上云前最好先购买实例进行测试，包括使用sysbench进行标准压测、业务兼容性测试、业务压测等来判断实例规格是否满足性能要求，，建议云上实例性能预留比如20-30%浮动空间

7、特别提醒，云上实例通常会把binlog以及SQL运行产生的临时表、临时文件也计入磁盘空间，此外云上数据表的碎片率可能会比自建实例大很多（曾经遇到本地5G的表云上占用120G），因此要特别注意磁盘空间要预留充足

8、最后说明一点，迁移云上最好选择数据库版本同自建版本。还有，尽量不要使用云上数据



**七、ALTER TABLE 出现duplicate primary xxx报错的原因及处理？**

好多同学都曾经问过这个问题，还有同学说这是bug，实际上这并不是bug

**(一)原因分析**

1、Online DDL操作时MySQL会将DML操作缓存起来存入到变更日志

2、等到DDL执行完成后再应用变更日志中的DML操作

3、在Oline DDL执行期间，并行的DML可能会没先检查唯一性直接插入一条相同主键的数据，这时并不会导致DDL报错，而是在DDL执行完成再次应用变更日志时才报错，最终导致DDL报错执行失败

**（二）问题说明**

其实这是Online DDL的正常情况，官方文档说明如下：

When running an in-place online DDL operation, the thread that runs the ALTER TABLE statement applies an online log of DML operations that were run concurrently on the same table from other connection threads. When the DML operations are applied, it is possible to encounter a duplicate key entry error (ERROR 1062 (23000): Duplicate entry), even if the duplicate entry is only temporary and would be reverted by a later entry in the online log. This is similar to the idea of a foreign key constraint check in InnoDB in which constraints must hold during a transaction

ref：[https://dev.mysql.com/doc/ref...](https://links.jianshu.com/go?to=https%3A%2F%2Flink.zhihu.com%2F%3Ftarget%3Dhttps%3A%2F%2Fdev.mysql.com%2Fdoc%2Frefman%2F8.0%2Fen%2Finnodb-online-ddl-limitations.html)

**（三）建议**

1、推荐使用pt-osc、gh-ost等第三方工具进行DDL操作

2、建议在业务低谷期进行操作



**八、InnoDB在什么情况下会触发检查点(checkpoint)?**

**（一）MySQL的checkpoint分类**

1、sharp checkpoint（激烈检查点，要求尽快将所有脏页都刷到磁盘上，对I/O资源的占有优先级高）

2、fuzzy checkpoint（模糊检查点，会根据系统负载及脏页数量适当平衡，不要求立即将所有脏页写入磁盘，这事默认的方式）

**（二）触发时机**

1、数据库正常关闭时，即innodb_fast_shutdown=0时需要执行sharp checkpoint

2、redo log发生切换时或者redo log快满的时候进行fuzzy checkpoint

3、master thread每隔1秒或10秒定期进行fuzzy checkpoint

4、innodb保证有足够多的空闲page，如果发现不足，需要移除lru list末尾的page，如果这些page是脏页，那么也需要fuzzy checkpoint

5、innodb buffer pool中脏页比超过innodb_max_dirty_pages_pct时也会触发fuzzy checkpoint

**（三）checkpoint相关参数及状态**

1、innodb_fast_shutdown

2、innodb_io_capacity/innodb_io_capacity_max

3、innodb_lru_scan_depth

4、innodb_max_dirty_pages_pct/innodb_max_dirty_pages_pct_lwm

5、Innodb_buffer_pool_pages_dirty/Innodb_buffer_pool_pages_total

6、Innodb_buffer_pool_wait_free



**九、update t set a=29 and b in (1,2,3,4);这样写有什么问题吗？**

**（一）SQL分析**

乍一看这个SQL貌似没有什么问题，本意是将t表中b的值属于1/2/3/4的数据的a列修改为29，但实际上该SQL是将t表数据的a列改成了条件29 & b in (1,2,3,4)的真假判断值

即：update t set a = (29 and b in (1,2,3,4));

修改后的SQL应该为update t set a = 29 where b in (1,2,3,4);

**（二）注意事项**

1、生产环境中进行批量数据修改时应该开启事务，修改确认后再进行提交操作

2、进行DML操作时，建议还是通过SQL审核工具审核后执行

3、建议打开sql_safe_updates选项，避免没有WHERE条件的更新、删除操作

------

有任何问题都可以加微信讨论，欢迎沟通~~ 互相进步！

**微信：lvqingshan_**
