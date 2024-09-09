-- 参数组表
drop table if exists parameter_groups;
create table if not exists parameter_groups (
    id int unsigned not null auto_increment,
    title varchar(50) not null default '' comment '参数名称',
    name varchar(50) not null default '' comment 'KEY',
    remark varchar(200) not null default '' comment '备注',
    primary key(id)
);

-- 写入部分数据
insert into parameter_groups (title, name, remark) values 
('注册信息校验', 'zhuce', '会员注册信息校验'),
('首存信息校验', 'first_charge', '存款信息校验'),
('首提信息校验', 'first_withdraw', '提款信息校验');

-- 增加参数组编号
alter table parameters add column group_id int unsigned not null default 0 comment '参数组编号';


-- 清空数据表
truncate table parameters;

-- 初始化数据
insert into parameters (title, name, `value`, remark, group_id) values 
('手机验证码', 'phone-number', '2', '1为必填, 2为非必填不可见, 3为必填不可见', '1'),
('银行卡验证', 'cards-number', '1', '1为必填, 2为非必填不可见, 3为必填不可见', '1'),
('邮箱验证码', 'email-number', '3', '1为必填, 2为非必填不可见, 3为必填不可见', '1');
