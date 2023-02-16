drop table if exists user_relation;
drop table if exists videos;
drop table if exists `user`;
drop table if exists `like_video`;
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
    key `user_vedio_idx` (`user_id`)
) engine = InnoDB
    default charset = utf8mb4
    collate = utf8mb4_general_ci;

create table if not exists user_relation
(
    `id`              bigint primary key auto_increment comment '用户ID',
    `created_at`      timestamp not null default current_timestamp comment 'User account create time',
    `updated_at`      timestamp not null default current_timestamp on update current_timestamp comment 'User account update time',
    `deleted_at`      timestamp null     default null comment 'User account delete time',
    `followers_count` int comment '粉丝列表人数',
    `follows_count`   int comment '关注列表人数',
    `likes_videos`    longtext comment '喜欢的视频列表',
    `followers`       longtext comment '粉丝列表, JSON表示',
    `follows`         longtext comment '关注列表, JSON表示'

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
)