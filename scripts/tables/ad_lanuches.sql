-- app启动页广告
drop table if exists ad_lanuches;
create table if not exists ad_lanuches (
    id int unsigned not null auto_increment,
    platform_type tinyint unsigned not null default 0 comment '投放平台, 0:主平台;1:体育;2:电竞;3:真人;4:电游;5:捕鱼;6:彩票;7:棋牌;',
    title varchar(50) not null default '' comment '启动标题',
    url_type tinyint unsigned not null default 0 comment '链接类型, 0:站内链接;1:站外链接;',
    url varchar(200) not null default '' comment '详情链接',
    image_android varchar(200) not null default '' comment '安卓图片',
    image_ios varchar(200) not null default '' comment '苹果图片',
    image_iosx varchar(200) not null default '' comment 'IPHONEX图片',
    state tinyint unsigned not null default 0 comment '状态, 0:停用;1:启用;',
    time_start int unsigned not null default 0 comment '开始时间',
    time_end int unsigned not null default 0 comment '结束时间',
    created int unsigned default 0 comment '创建时间',
    updated int unsigned default 0 comment '修改时间',
    primary key(id)
);


insert into ad_lanuches (platform_type, title, url_type, url, time_start, time_end, created, updated) values 
(1, '测试投放广告-01', 1, '', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp()),
(1, '测试投放广告-02', 1, '', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp()),
(2, '测试投放广告-03', 1, '', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp()),
(5, '测试投放广告-04', 1, '', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp()),
(4, '测试投放广告-05', 1, '', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp()),
(1, '测试投放广告-06', 1, '', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp()),
(2, '测试投放广告-07', 1, '', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp()),
(3, '测试投放广告-08', 1, '', unix_timestamp(), unix_timestamp() + 86400 * 100, unix_timestamp(), unix_timestamp());
