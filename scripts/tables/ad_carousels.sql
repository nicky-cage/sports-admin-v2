-- 轮播图
drop table if exists ad_carousels;
create table if not exists ad_carousels (
    id int unsigned not null auto_increment,
    device_type tinyint not null default 0 comment '轮播设备, 1:WEB;2:H5/App;',
    title varchar(200) not null default '' comment '轮播标题',
    direct_type tinyint unsigned not null default 0 comment '跳转类型',
    url_type tinyint unsigned not null default 0 comment '链接类型',
    url varchar(200) not null default '' comment '链接内容',
    activity_id int unsigned not null default 0 comment '活动编号',
    image varchar(200) not null default '' comment '轮播图片',
    time_start int unsigned not null default 0 comment '开始时间',
    time_end int unsigned not null default 0 comment '结束时间',
    state tinyint unsigned not null default 0 comment '状态',
    sort int not null default 0 comment '排序',
    created int unsigned default 0 comment '创建时间',
    updated int unsigned default 0 comment '修改时间',
    primary key(id)
);


insert into ad_carousels (device_type, title, time_start, time_end, created, updated) values
(1, '测试轮播图-01', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp()),
(1, '测试轮播图-02', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp()),
(2, '测试轮播图-03', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp()),
(1, '测试轮播图-04', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp()),
(2, '测试轮播图-05', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp()),
(1, '测试轮播图-06', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp());
