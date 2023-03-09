
drop table if exists `video`;
drop table if exists `user`;
drop table if exists `like_video`;
drop table if exists `message`;
drop table if exists `comment_video`;
drop table if exists `follow`;
-- 用户表 users --
create table if not exists `user`
(
    `id`         bigint       not null auto_increment comment 'PK',
    `name`       varchar(128) not null default '' comment '用户名',
    `password`   varchar(128) not null default '' comment '密码',
    `created_at` timestamp    not null default current_timestamp comment '创建时间',
    `updated_at` timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',
    `deleted_at` timestamp    null     default null comment '删除时间',
    primary key (`id`),
    unique key (`name`),
    key `name_password_idx` (`name`, `password`)
) engine = InnoDB
  default charset = utf8mb4
  collate = utf8mb4_general_ci comment ='用户表';
create table if not exists video
(
    `id`              bigint primary key auto_increment comment '视频ID',
    `created_at`      timestamp not null default current_timestamp comment '创建时间',
    `updated_at`      timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
    `deleted_at`      timestamp null     default null comment '删除时间',
    `title`           varchar(255) comment '视频标题',
    `label`           varchar(255) comment '视频标签',
    `owner_id`        bigint    not null comment '视频所有者',
    `likes_count`     bigint             default 0 comment '点赞人数',
    `comment_archive` longtext comment '评论信息的最新归档',
    foreign key (owner_id) references `user` (id)
) engine = InnoDB
  default charset = utf8mb4
  collate = utf8mb4_general_ci comment '视频信息表';
create table if not exists `like_video`
(
    `id`         bigint not null auto_increment comment 'PK',
    `user_id`    bigint not null comment '用户id' ,
    `video_id`   bigint not null comment '点赞的视频id',
    `cancel`     tinyint not null comment '是否取消',
    `created_at` timestamp    not null default current_timestamp comment '创建时间',
    `updated_at` timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',
    `deleted_at` timestamp    null     default null comment '删除时间',
    primary key (`id`),
    foreign key `user_id` (`user_id`) references `user` (id),
    foreign key `video_id` (`video_id`) references `video` (id)
) engine = InnoDB
  default charset = utf8mb4
  collate = utf8mb4_general_ci comment '点赞表';

create table if not exists `comment_video`
(
    `id`          bigint not null auto_increment,
    `user_id`     bigint not null comment '用户id',
    `video_id`    bigint not null comment '评论的视频id',
    `comment_text`varchar(128) not null default '' comment '评论内容',
    `created_at`  timestamp    not null default current_timestamp comment '创建时间',
    `updated_at`  timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',
    `deleted_at`  timestamp    null     default null comment '删除时间',
    primary key (`id`),
    foreign key `user_id` (`user_id`) references `user` (id),
    foreign key `video_id` (`video_id`) references `video` (id)
) engine = InnoDB
  default charset = utf8mb4
  collate = utf8mb4_general_ci comment '评论表';

create table if not exists follow
(
    `id`              bigint primary key auto_increment comment '记录ID',
    `followed_id`         bigint comment '被关注者id',
    `follower_id`     bigint comment '关注者id',
    `cancel`          tinyint not null comment '是否取消关注',
    `created_at`  timestamp    not null default current_timestamp comment '创建时间',
    `updated_at`  timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',
    `deleted_at`  timestamp    null     default null comment '删除时间',
    foreign key `followed_idx` (`followed_id`) references `user` (id),
    foreign key `follow_idx` (`follower_id`) references `user` (id)
) engine = InnoDB
  default charset = utf8mb4
  collate = utf8mb4_general_ci comment '用户社交信息表';




create table if not exists `message`
(
    `id`              bigint primary key auto_increment comment '消息ID',
    `created_at`      timestamp not null default current_timestamp comment '创建时间',
    `updated_at`      timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
    `deleted_at`      timestamp null     default null comment '删除时间',
    `to_user_id`      bigint not null  comment '接收者id',
    `from_user_id`    bigint not null  comment '发送者id',
    `content`         varchar(255)     comment '消息内容',
    `create_time`     bigint null      comment '创建时间int64',
    foreign key (to_user_id) references `user` (id),
    foreign key (from_user_id) references `user` (id)
)engine = InnoDB
 default charset = utf8mb4
 collate = utf8mb4_general_ci comment '消息表';;
