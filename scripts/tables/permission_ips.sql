drop table if exists permission_ips;
create table if not exists permission_ips (
    id int unsigned not null auto_increment,
    permission_type tinyint unsigned not null default 0 comment '授权类型, 0:授权访问',
    ip varchar(30) not null default '' comment 'IP',
    remark varchar(200) not null default '' comment '备注',
    state tinyint not null default 0 comment '状态',
    created int unsigned not null default 0 comment '添加时间', 
    updated int unsigned not null default 0 comment '修改时间', 
    primary key(id)
);

insert into permission_ips (permission_type, ip, remark, created, updated) values 
(0, '127,0.0.1', '测试备注内容', unix_timestamp(), unix_timestamp()),
(0, '127,0.0.1', '测试备注内容', unix_timestamp(), unix_timestamp()),
(0, '127,0.0.1', '测试备注内容', unix_timestamp(), unix_timestamp()),
(0, '127,0.0.1', '测试备注内容', unix_timestamp(), unix_timestamp()),
(0, '127,0.0.1', '测试备注内容', unix_timestamp(), unix_timestamp());
