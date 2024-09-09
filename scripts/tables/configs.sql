-- 配置信息
drop table if exists configs;
create table if not exists configs (
    id int unsigned not null auto_increment,
    name varchar(50) not null default '' comment '配置模板名称',
    site_name varchar(50) not null default '' comment '后台名称',
    site_title varchar(50) not null default '' comment '网站标题',
    logo_site varchar(200) not null default '' comment '站点LOGO',
    logo_web varchar(200) not null default '' comment '网站LOGO',
    logo_title varchar(200) not null default '' comment '网页标题LOGO',
    logo_app varchar(200) not null default '' comment '移动端LOGO',
    domain_pc varchar(100) not null default '' comment 'PC主域名',
    domain_app varchar(100) not null default '' comment '综合站APP域名',
    domain_app_phy varchar(100) not null default '' comment '体育APP域名',
    domain_h5 varchar(100) not null default '' comment 'H5域名',
    domain_agent_pc varchar(100) not null default '' comment '代理PC域名',
    domain_agent_app varchar(100) not null default '' comment '代理APP域名',
    agent_qq varchar(20) not null default '' comment '代理QQ',
    agent_skype varchar(20) not null default '' comment '代理Skype',
    agent_sugram varchar(20) not null default '' comment '代理Sugram',
    admin_id int unsigned not null default 0 comment '管理员编号',
    admin_name varchar(20) not null default '' comment '管理员名称',
    created int unsigned not null default 0 comment '添加时间',
    updated int unsigned not null default 0 comment '最后修改',
    state tinyint unsigned not null default 0 comment '状态, 0:未使用; 1:使用中',
    primary key(id)
);

-- 初始化数据库数据
insert into configs (name, created, updated, state) values 
('默认模板', unix_timestamp(), unix_timestamp(), 1);
