-- 帮助表
drop table if exists helps;
create table if not exists helps (
    id int unsigned not null auto_increment,
    category_id int unsigned not null default 0 comment '分类编号',
    title varchar(100) not null default '' comment '标题',
    sort int not null default 0 comment '排序',
    venue_type tinyint not null default 0 comment '场馆类型, 0:没有关联; 1:存款; 2:取款; 3:转账; 4:奖励; 5:VIP等级; 6:合营计划;',
    content_type tinyint not null default 0 comment '内容类型, 0:文本; 1:链接',
    terminal char(4) not null default '' comment '终端类型, 0:PC; 1:移动;',
    content text comment '内容',
    admin_id int unsigned not null default 0 comment '管理员编号',
    admin_name varchar(20) not null default '' comment '管理员名称',
    created int unsigned not null default 0 comment '创建时间',
    updated int unsigned not null default 0 comment '最后修改',
    img    varchar(255)  comment '图片',
    state tinyint default 1 commint '状态  0关闭 1开启',
    primary key(id)
);

insert into helps (category_id, title, venue_type, content_type, terminal, created, updated) values 
(1, '测试文章标题-01', '1', '0', '0,1', unix_timestamp(), unix_timestamp()),
(1, '测试文章标题-02', '1', '0', '0,1', unix_timestamp(), unix_timestamp()),
(1, '测试文章标题-03', '1', '0', '0,1', unix_timestamp(), unix_timestamp()),
(1, '测试文章标题-04', '1', '0', '0,1', unix_timestamp(), unix_timestamp()),
(1, '测试文章标题-05', '1', '0', '0,1', unix_timestamp(), unix_timestamp()),
(1, '测试文章标题-06', '1', '0', '0,1', unix_timestamp(), unix_timestamp()),
(1, '测试文章标题-07', '1', '0', '0,1', unix_timestamp(), unix_timestamp()),
(1, '测试文章标题-08', '1', '0', '0,1', unix_timestamp(), unix_timestamp());
