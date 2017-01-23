CREATE TABLE `material_library` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `cover` varchar(255) NOT NULL DEFAULT '',
  `name` varchar(255) NOT NULL DEFAULT '',
  `url` varchar(512) NOT NULL DEFAULT '',
  `sha` varchar(255) NOT NULL DEFAULT '',
  `version` varchar(255) NOT NULL DEFAULT '',
  `mate_info` text,
  `hidden_at` int(11) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL,
  `material_type` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;