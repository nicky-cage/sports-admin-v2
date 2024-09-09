-- 赞助配置
drop table if exists ad_sponsors;
create table if not exists ad_sponsors (
    id int unsigned not null auto_increment,
    image varchar(200) not null default '' comment '赞助图片',
    icon varchar(200) not null default '' comment '上传图标',
    line_first varchar(100) not null default '' comment '第一行文案',
    line_second varchar(100) not null default '' comment '第二行文案',
    line_key varchar(100) not null default '' comment '按键文案',
    url varchar(200) not null default '' comment '详情链接',
    share_title varchar(100) not null default '' comment '分享标题',
    share_image varchar(200) not null default '' comment '分享图片',
    title varchar(100) not null default '' comment '赞助标题',
    time_start int unsigned not null default 0 comment '开始时间',
    time_end int unsigned not null default 0 comment '结束时间',
    sort int not null default 0 comment '排序',
    state tinyint unsigned not null default 0 comment '状态',
    created int unsigned default 0 comment '创建时间',
    updated int unsigned default 0 comment '修改时间',
    primary key(id)
);

insert into ad_sponsors (line_first, line_second, share_title, title, time_start, time_end, created, updated) values
('上行文案-01', '下行方案-01', '分享标题-01', '准备分享标题-01', unix_timestamp(), unix_timestamp() + 86400 * 1, unix_timestamp(), unix_timestamp()),
('上行文案-02', '下行方案-02', '分享标题-02', '准备分享标题-02', unix_timestamp(), unix_timestamp() + 86400 * 2, unix_timestamp(), unix_timestamp()),
('上行文案-03', '下行方案-03', '分享标题-03', '准备分享标题-03', unix_timestamp(), unix_timestamp() + 86400 * 3, unix_timestamp(), unix_timestamp()),
('上行文案-04', '下行方案-04', '分享标题-04', '准备分享标题-04', unix_timestamp(), unix_timestamp() + 86400 * 4, unix_timestamp(), unix_timestamp()),
('上行文案-05', '下行方案-05', '分享标题-05', '准备分享标题-05', unix_timestamp(), unix_timestamp() + 86400 * 5, unix_timestamp(), unix_timestamp()),
('上行文案-06', '下行方案-06', '分享标题-06', '准备分享标题-06', unix_timestamp(), unix_timestamp() + 86400 * 6, unix_timestamp(), unix_timestamp()),
('上行文案-07', '下行方案-07', '分享标题-07', '准备分享标题-07', unix_timestamp(), unix_timestamp() + 86400 * 7, unix_timestamp(), unix_timestamp());
