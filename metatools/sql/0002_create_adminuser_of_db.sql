DROP USER IF EXISTS gfadmin;
CREATE USER `gfadmin`@`%` IDENTIFIED BY 'gfadmin';
GRANT ALL ON `gfalcon`.* TO `gfadmin`@`%` identified BY 'gfadmin';
GRANT ALL ON `gfalcon`.* TO `gfadmin`@`localhost` identified BY 'gfadmin';
