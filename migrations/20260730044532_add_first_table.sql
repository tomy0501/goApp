-- +goose Up
-- +goose StatementBegin
CREATE TABLE `facilities` (
  `id` binary(16) NOT NULL DEFAULT (uuid_to_bin(uuid(), 1)),
  `name` varchar(100) NOT NULL,
  `address` varchar(255) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `facility_users` (
  `id` binary(16) NOT NULL DEFAULT (uuid_to_bin(uuid(), 1)),
  `facility_id` binary(16) NOT NULL,
  `user_id` binary(16) NOT NULL,
  `role` enum('owner','manager','front','housekeeping','staff') NOT NULL DEFAULT 'staff',
  `status` enum('active','invited','disabled') NOT NULL DEFAULT 'active',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  UNIQUE KEY `uk_facility_user` (`facility_id`,`user_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_facility_id` (`facility_id`),
  CONSTRAINT `fk_facility_users_facility` FOREIGN KEY (`facility_id`) REFERENCES `facilities` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `guest` (
  `id` binary(16) NOT NULL DEFAULT (uuid_to_bin(uuid(), 1)),
  `customer_name` varchar(50) NOT NULL,
  `customer_age` int DEFAULT NULL,
  `customer_sex` enum('male','female','other','unknown') DEFAULT NULL,
  `customer_from` varchar(200) DEFAULT NULL,
  `created_date` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_date` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `reservation` (
  `id` binary(16) NOT NULL DEFAULT (uuid_to_bin(uuid(), 1)),
  `facility_id` binary(16) NOT NULL,
  `represent_guest_id` binary(16) NOT NULL,
  `created_date` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_date` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `fk_1` (`facility_id`),
  KEY `fk_2` (`represent_guest_id`),
  CONSTRAINT `fk_1` FOREIGN KEY (`facility_id`) REFERENCES `facilities` (`id`),
  CONSTRAINT `fk_2` FOREIGN KEY (`represent_guest_id`) REFERENCES `guest` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `reservation_guest` (
  `id` binary(16) NOT NULL DEFAULT (uuid_to_bin(uuid(), 1)),
  `reservation_id` binary(16) NOT NULL,
  `guest_id` binary(16) NOT NULL,
  `room_id` binary(16) NOT NULL,
  `checkin` datetime DEFAULT NULL,
  `checkout` datetime DEFAULT NULL,
  `created_date` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_date` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `reservation_guest_string` (`reservation_id`,`guest_id`),
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `fk_2` (`guest_id`),
  KEY `fk_3` (`room_id`),
  CONSTRAINT `fk_1` FOREIGN KEY (`reservation_id`) REFERENCES `reservation` (`id`),
  CONSTRAINT `fk_2` FOREIGN KEY (`guest_id`) REFERENCES `guest` (`id`),
  CONSTRAINT `fk_3` FOREIGN KEY (`room_id`) REFERENCES `room` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


CREATE TABLE `room` (
  `id` binary(16) NOT NULL DEFAULT (uuid_to_bin(uuid(), 1)),
  `facility_id` binary(16) NOT NULL,
  `room_name` varchar(50) NOT NULL,
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `fk_1` (`facility_id`),
  CONSTRAINT `fk_1` FOREIGN KEY (`facility_id`) REFERENCES `facilities` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


CREATE TABLE `room_type` (
  `id` binary(16) NOT NULL DEFAULT (uuid_to_bin(uuid(), 1)),
  `room_id` binary(16) NOT NULL,
  `room_type_name` varchar(50) DEFAULT NULL,
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `fk_1` (`room_id`),
  CONSTRAINT `fk_1` FOREIGN KEY (`room_id`) REFERENCES `room` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE `seasonality` (
  `id` binary(16) NOT NULL DEFAULT (uuid_to_bin(uuid(), 1)),
  `facility_id` binary(16) NOT NULL,
  `seasonality_name` varchar(100) NOT NULL,
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `fk_1` (`facility_id`),
  CONSTRAINT `fk_1` FOREIGN KEY (`facility_id`) REFERENCES `facilities` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


CREATE TABLE `seasonality_number_of_people` (
  `id` binary(16) NOT NULL DEFAULT (uuid_to_bin(uuid(), 1)),
  `seasonality_id` binary(16) NOT NULL,
  `number_of_people` int NOT NULL,
  `created_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `seasonality_people_strict` (`seasonality_id`,`number_of_people`),
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  CONSTRAINT `fk_1` FOREIGN KEY (`seasonality_id`) REFERENCES `seasonality` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- +goose StatementEnd