-- 站点配置 -底部信息
drop table if exists site_bottoms;
create table if not exists site_bottoms (
    id int unsigned not null auto_increment,
    bottom_type tinyint not null default 0 comment '信息类型',
    content_type tinyint not null default 0 comment '内容类型',
    url_type tinyint not null default 0 comment '链接类型',
    image varchar(200) not null default '' comment '图片',
    url varchar(200) not null default '' comment '链接地址',
    title varchar(100) not null default '' comment '名称',
    sort int not null default 0 comment '排序',
    content varchar(300) not null default '' comment '内容',
    created int unsigned not null default 0 comment '创建时间',
    updated int unsigned not null default 0 comment '修改时间',
    primary key(id)
);

-- 初始数据
insert into site_bottoms (bottom_type, created, updated) values 
(1, unix_timestamp(), unix_timestamp()),
(2, unix_timestamp(), unix_timestamp()),
(3, unix_timestamp(), unix_timestamp()),
(4, unix_timestamp(), unix_timestamp()),
(5, unix_timestamp(), unix_timestamp()),
(6, unix_timestamp(), unix_timestamp());

-- 新加字段
alter table site_bottoms add column state tinyint unsigned not null default 0 comment '状态';
