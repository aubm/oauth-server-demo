DROP TABLE IF EXISTS `access`;

CREATE TABLE `access` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `access_token` varchar(255) NOT NULL DEFAULT '',
  `refresh_token` varchar(255) NOT NULL DEFAULT '',
  `expires_in` int(11) unsigned NOT NULL,
  `scope` varchar(255) NOT NULL DEFAULT '',
  `redirect_uri` varchar(255) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `client_id` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `access_belongs_to_clients` (`client_id`),
  CONSTRAINT `access_belongs_to_clients` FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `clients`;

CREATE TABLE `clients` (
  `id` varchar(255) NOT NULL DEFAULT '',
  `secret` varchar(255) NOT NULL DEFAULT '',
  `redirect_uri` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
