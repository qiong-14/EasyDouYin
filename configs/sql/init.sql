
drop table if exists user_relation;
drop table if exists videos;
drop table if exists `user`;
drop table if exists `like_video`;
drop table if exists `message`;
drop table if exists `comment_video`;
drop table if exists `follows`;
-- 用户表 users --
create table if not exists `user`
(
    `id`         bigint       not null auto_increment comment 'PK',
    `name`       varchar(128) not null default '' comment 'UserName',
    `password`   varchar(128) not null default '' comment 'UserPassword',
    `created_at` timestamp    not null default current_timestamp comment 'User account create time',
    `updated_at` timestamp    not null default current_timestamp on update current_timestamp comment 'User account update time',
    `deleted_at` timestamp    null     default null comment 'User account delete time',
    primary key (`id`),
    key `name_password_idx` (`name`, `password`)
) engine = InnoDB
  default charset = utf8mb4
  collate = utf8mb4_general_ci comment ='User table';

create table if not exists `like_video`
(
    `id`         bigint not null auto_increment,
    `user_id`    bigint not null,
    `video_id`   bigint not null,
    `cancel`     tinyint not null,
    `created_at` timestamp    not null default current_timestamp,
    `updated_at` timestamp    not null default current_timestamp on update current_timestamp,
    `deleted_at` timestamp    null     default null,
    primary key (`id`),
    key `user_video_idx` (`user_id`, `video_id`)
) engine = InnoDB
    default charset = utf8mb4
    collate = utf8mb4_general_ci;

create table if not exists `comment_video`
(
    `id`          bigint not null auto_increment,
    `user_id`     bigint not null,
    `video_id`    bigint not null,
    `comment_text`varchar(128) not null default '' comment 'CommentText',
    `created_at`  timestamp    not null default current_timestamp,
    `updated_at`  timestamp    not null default current_timestamp on update current_timestamp,
    `deleted_at`  timestamp    null     default null,
    primary key (`id`)
) engine = InnoDB
    default charset = utf8mb4
    collate = utf8mb4_general_ci;

create table if not exists user_relation
(
    `user_id`         bigint primary key auto_increment comment 'User id',
    `created_at`      timestamp not null default current_timestamp comment 'User account create time',
    `updated_at`      timestamp not null default current_timestamp on update current_timestamp comment 'User account update time',
    `deleted_at`      timestamp null     default null comment 'User account delete time',
    `name`            varchar(128) not null default '' comment 'User name',
    `is_follow`       boolean comment 'If follow',
    `follower_id`     int comment 'Follower id',
    `follow_count`    int comment 'Follow count',
    `follower_count`  int comment 'Follower count'
)engine = InnoDB;

create table if not exists user_info
(
    `user_id`               bigint primary key auto_increment comment 'User id',
    `created_at`            timestamp not null default current_timestamp comment 'User account create time',
    `updated_at`            timestamp not null default current_timestamp on update current_timestamp comment 'User account update time',
    `deleted_at`            timestamp null     default null comment 'User account delete time',
    `name`                  varchar(128) not null default '' comment 'User name',
    `is_follow`             boolean comment 'If follow',
    `follower_id`           int comment 'Follower id',
    `follow_count`          int comment 'Follow count',
    `follower_count`        int comment 'Follower count',
    `avatar`                varchar(255) comment 'User avatar',
    `background_image`      varchar(255) comment 'User background_image',
    `signature`             varchar(255) comment 'Signature',
    `total_favorited`       int comment 'Total favorited',
    `work_count`            int comment 'Work count',
    `favorite_count`        int comment 'Favorite count'
)engine = InnoDB; 


create table if not exists follows
(
    `id`              bigint primary key auto_increment comment '记录ID',
    `followed_id`         bigint comment '被关注者id',
    `follower_id`     bigint comment '关注者id',
    `cancel`          tinyint not null,
    `created_at`  timestamp    not null default current_timestamp,
    `updated_at`  timestamp    not null default current_timestamp on update current_timestamp,
    `deleted_at`  timestamp    null     default null
) comment '用户社交信息表';

create table if not exists videos
(
    `id`              bigint primary key auto_increment comment '视频ID',
    `created_at`      timestamp not null default current_timestamp comment 'User account create time',
    `updated_at`      timestamp not null default current_timestamp on update current_timestamp comment 'User account update time',
    `deleted_at`      timestamp null     default null comment 'User account delete time',
    `title`           varchar(255) comment '视频标题',
    `label`           varchar(255) comment '视频标签',
    `owner_id`        bigint    not null comment '视频所有者',
    `likes_count`     bigint             default 0 comment '点赞人数',
    `comment_archive` longtext comment '评论信息的最新归档',
    foreign key (owner_id) references `user` (id)
);


create table if not exists `message`
(
    `id`              bigint primary key auto_increment comment '消息ID',
    `created_at`      timestamp not null default current_timestamp comment 'message create time',
    `updated_at`      timestamp not null default current_timestamp on update current_timestamp comment 'message update time',
    `deleted_at`      timestamp null     default null comment 'message delete time',
    `to_user_id`      bigint not null  comment '接收者id',
    `from_user_id`    bigint not null  comment '发送者id',
    `content`         varchar(255)     comment '消息内容',
    `create_time`     bigint null      comment '创建时间int64'
);
