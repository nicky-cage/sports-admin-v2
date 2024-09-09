alter table admins change username name varchar(20) not null default '' comment '管理员名称';

alter table admins add column mail varchar(100) not null default '' comment '邮箱账号';


update admins set created = unix_timestamp() where created = 0;
update admins set updated = unix_timestamp() where updated = 0;

alter table admins change password password char(32) not null default '' comment '密码';

alter table admins add column secret char(32) not null default '' comment '密盐' after password;

alter table admins drop column secret;

alter table admins change status state tinyint unsigned not null default 0 comment '状态';

alter table admins add column nickname varchar(50) not null default '' comment '昵称';
