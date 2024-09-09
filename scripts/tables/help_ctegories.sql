-- 帮助分类
drop table if exists help_categories;
create table if not exists help_categories (
    id int unsigned not null auto_increment,
    name varchar(50) not null default '' comment '名称',
    title varchar(50) not null default '' comment '标题',
    icon varchar(200) not null default '' comment '图标',
    created int unsigned not null default 0 comment '创建时间',
    updated int unsigned not null default 0 comment '最后修改',
    primary key(id)
);

-- 添加初始数据
insert into help_categories (name, title, created, updated) values
('新手帮助', '新手帮助', unix_timestamp(), unix_timestamp()),
('企业事务', '企业事务', unix_timestamp(), unix_timestamp()),
('新手帮助', '新手帮助', unix_timestamp(), unix_timestamp()),
('联系我们', '联系我们', unix_timestamp(), unix_timestamp());

