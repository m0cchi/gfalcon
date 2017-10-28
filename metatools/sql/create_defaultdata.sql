SOURCE create_table.sql;

INSERT INTO `services`(`id`) VALUE ('gfalcon');
INSERT INTO `actions`(`service_iid`,`id`) SELECT `iid` as `service_iid`, 'all' as `id` FROM `services` WHERE `services`.`id` = 'gfalcon';
INSERT INTO `teams`(`id`) VALUE ('gfalcon');
INSERT INTO `users`(`team_iid`,`id`) SELECT `iid` as `team_iid`, 'gfadmin' as `id` FROM `teams` WHERE `teams`.`id` = 'gfalcon';

