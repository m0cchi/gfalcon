INSERT INTO `services`(`id`) VALUE ('gfalcon');
INSERT INTO `actions`(`service_iid`,`id`) SELECT `iid` as `service_iid`, 'all' as `id` FROM `services` WHERE `services`.`id` = 'gfalcon';
INSERT INTO `teams`(`id`) VALUE ('gfalcon');
INSERT INTO `users`(`team_iid`,`id`) SELECT `iid` as `team_iid`, 'gfadmin' as `id` FROM `teams` WHERE `teams`.`id` = 'gfalcon';

INSERT INTO `action_links`(`action_iid`, `user_iid`) select `actions`.`iid` as `action_iid` ,`users`.`iid` as `user_iid` from (select `iid` from `users`) as `users`, (select `iid` from `actions`) as `actions` where `users`.`iid` = 1 and `actions`.`iid` = 1;
