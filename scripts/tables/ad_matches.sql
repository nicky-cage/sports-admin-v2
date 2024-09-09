-- 体育赛事
drop table if exists ad_matches;
create table if not exists ad_matches (
    id int unsigned not null auto_increment,
    title varchar(200) not null default '' comment '赛事名称',
    match_date char(10) not null default '2020-01-01' comment '比赛时间',
    team_first varchar(50) not null default '' comment '主队名称',
    team_first_icon varchar(200) not null default '' comment '主队队标',
    team_second varchar(50) not null default '' comment '客队名称',
    team_second_icon varchar(200) not null default '' comment '客队队标',
    time_start int unsigned not null default 0 comment '开始时间',
    time_end int unsigned not null default 0 comment '结束时间',
    state tinyint unsigned not null default 0 comment '状态',
    created int unsigned default 0 comment '创建时间',
    updated int unsigned default 0 comment '修改时间',
    primary key(id)
);

insert into ad_matches (title, match_date, team_first, team_second, time_start, time_end, created, updated) values
('体育比赛-01', '2020-10-10', '比赛甲队-01', '比赛丙队-01', unix_timestamp(), unix_timestamp() + 86400 * 1, unix_timestamp(), unix_timestamp()),
('体育比赛-02', '2020-10-10', '比赛甲队-02', '比赛丙队-02', unix_timestamp(), unix_timestamp() + 86400 * 2, unix_timestamp(), unix_timestamp()),
('体育比赛-03', '2020-10-10', '比赛甲队-03', '比赛丙队-03', unix_timestamp(), unix_timestamp() + 86400 * 3, unix_timestamp(), unix_timestamp()),
('体育比赛-04', '2020-10-10', '比赛甲队-04', '比赛丙队-04', unix_timestamp(), unix_timestamp() + 86400 * 4, unix_timestamp(), unix_timestamp()),
('体育比赛-05', '2020-10-10', '比赛甲队-05', '比赛丙队-05', unix_timestamp(), unix_timestamp() + 86400 * 5, unix_timestamp(), unix_timestamp()),
('体育比赛-06', '2020-10-10', '比赛甲队-06', '比赛丙队-06', unix_timestamp(), unix_timestamp() + 86400 * 6, unix_timestamp(), unix_timestamp());
