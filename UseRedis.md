# Redis在项目中的引用

## 视频流接口(不做)

1. 当大量用户存在时，可能在同一秒内就有lastTime相同的Feed请求发过来。当前，不需要对用户画像做个性化推荐，可以用Redis缓存一些视频的链接，减小MinIO生成的链接数量以降低其资源占用，单机版可以用本地内存即可, 分布式可以考虑Redis

## 用户鉴权

> 当用户量极大且分布式时，可以用Redis缓存用户Token，这个需求可以做
>
> **分配DB:** 0
>
> Token过期应当自动删除, 考虑用string
>
> key命名：`token:user:<userid>`
>
> 值：Token字符串

## 用户信息

> 项目中，用户信息的查询是一个很频繁的操作，在返回Feed流或返回用户点赞过的视频的时候，需要对user表进行查询，用户信息最好缓存到Redis中
>
> **分配DB:** 1
>
> 考虑用string
>
> key命名：`info:user:<userid>`
>
> 值：User结构体的JSON表示
>
> 淘汰策略（WatchDog）：访问的时候，更新key的时间（expire）；长时间没有访问，就淘汰。
>
> 为防止缓存雪崩，每次设置的过期时间加一个随机偏移。

## 视频信息结构体

> 同上，分配DB：2
>
> key命名：`info:video:<videoId>`

## 用户点赞过的视频，视频的点赞数(用户点赞操作)

>前者可以用Redis的Zset(有序集合)数据结构(点赞时间排序)，将用户点赞过的视频写入到数据库中
>
>注意到该操作在SQL执行的时候没有加Limit限制，如果用户点赞过的视频很多，则容易对数据库造成负载。
>
>> 建议：增加limit限制，仅保留最后点赞的100条视频
>
>分配DB：3
>
>key命名: `fav:user:<userId>`

> 视频的点赞数：用INCR和DECR
>
> 分配DB：3
>
> key命名: `fav:video:<videoId>`

> 点赞数据应该适时和MySQL同步。

## 用户关注列表和用户粉丝列表

> 这个操作如果用纯RDBMS实现，则会比较耗时，可以用zset存储
>
> 分配DB：4
>
> 前者：key命名：`follows:<userId>`
>
> 后者：key命名：`fans:<userId>`