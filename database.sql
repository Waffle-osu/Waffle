-- MySQL dump 10.13  Distrib 8.0.28, for Win64 (x86_64)
--
-- Host: localhost    Database: waffle
-- ------------------------------------------------------
-- Server version	8.0.28

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `beatmap_comments`
--

DROP TABLE IF EXISTS `beatmap_comments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `beatmap_comments` (
  `comment_id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `beatmap_id` int NOT NULL,
  `beatmapset_id` int NOT NULL,
  `score_id` bigint unsigned NOT NULL,
  `time` bigint NOT NULL,
  `target` tinyint NOT NULL,
  `comment` text NOT NULL,
  `format_string` varchar(16) NOT NULL,
  PRIMARY KEY (`comment_id`),
  KEY `beatmapid_index` (`beatmap_id`,`target`),
  KEY `beatmapset_id_index` (`beatmapset_id`,`target`) /*!80000 INVISIBLE */,
  KEY `score_id_index` (`score_id`,`target`) /*!80000 INVISIBLE */,
  KEY `userid_beatmap_comments_pk_idx` (`user_id`) /*!80000 INVISIBLE */,
  CONSTRAINT `userid_beatmapcomments_pk` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `beatmap_favourites`
--

DROP TABLE IF EXISTS `beatmap_favourites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `beatmap_favourites` (
  `favourite_id` bigint NOT NULL AUTO_INCREMENT,
  `beatmapset_id` int NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`favourite_id`,`user_id`,`beatmapset_id`),
  KEY `userid_index` (`user_id`),
  CONSTRAINT `userid_pk` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `beatmap_ratings`
--

DROP TABLE IF EXISTS `beatmap_ratings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `beatmap_ratings` (
  `beatmapset_id` int NOT NULL,
  `rating_sum` bigint NOT NULL DEFAULT '0',
  `votes` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`beatmapset_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `beatmap_ratings_submissions`
--

DROP TABLE IF EXISTS `beatmap_ratings_submissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `beatmap_ratings_submissions` (
  `submission_id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `beatmapset_id` int NOT NULL,
  PRIMARY KEY (`submission_id`,`user_id`,`beatmapset_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `beatmaps`
--

DROP TABLE IF EXISTS `beatmaps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `beatmaps` (
  `beatmap_id` int NOT NULL AUTO_INCREMENT,
  `beatmapset_id` int NOT NULL,
  `creator_id` bigint NOT NULL,
  `filename` varchar(256) NOT NULL,
  `beatmap_md5` varchar(32) NOT NULL,
  `version` varchar(128) NOT NULL,
  `total_length` int NOT NULL,
  `drain_time` int NOT NULL,
  `count_objects` int NOT NULL,
  `count_normal` int NOT NULL,
  `count_slider` int NOT NULL,
  `count_spinner` int NOT NULL,
  `diff_hp` tinyint NOT NULL,
  `diff_cs` tinyint NOT NULL,
  `diff_od` tinyint NOT NULL,
  `diff_stars` float NOT NULL,
  `playmode` tinyint NOT NULL,
  `ranking_status` tinyint NOT NULL,
  `last_update` datetime NOT NULL,
  `submit_date` datetime NOT NULL,
  `approve_date` datetime NOT NULL,
  `beatmap_source` tinyint NOT NULL COMMENT 'Where the beatmap is from, used for when i may or may not add BSS and might be used for oldsu maps',
  PRIMARY KEY (`beatmap_id`),
  KEY `checksum_index` (`beatmap_md5`),
  KEY `filename_index` (`filename`),
  KEY `setid_index` (`beatmap_id`) /*!80000 INVISIBLE */,
  KEY `source_index` (`beatmap_source`)
) ENGINE=InnoDB AUTO_INCREMENT=3483833 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `beatmapsets`
--

DROP TABLE IF EXISTS `beatmapsets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `beatmapsets` (
  `beatmapset_id` int NOT NULL AUTO_INCREMENT,
  `creator_id` bigint NOT NULL,
  `artist` varchar(256) NOT NULL,
  `title` varchar(256) NOT NULL,
  `creator` varchar(256) NOT NULL,
  `source` varchar(256) NOT NULL,
  `tags` varchar(1024) NOT NULL,
  `has_video` tinyint NOT NULL,
  `has_storyboard` tinyint NOT NULL,
  `bpm` float NOT NULL,
  PRIMARY KEY (`beatmapset_id`),
  KEY `title_index` (`title`) /*!80000 INVISIBLE */,
  KEY `artist_index` (`title`) /*!80000 INVISIBLE */,
  KEY `creator_index` (`creator`),
  FULLTEXT KEY `fulltext_search` (`artist`,`title`,`creator`,`source`,`tags`)
) ENGINE=InnoDB AUTO_INCREMENT=41895 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `failtimes`
--

DROP TABLE IF EXISTS `failtimes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `failtimes` (
  `failtime_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `failtime` int NOT NULL,
  `beatmap_id` int NOT NULL,
  `score_id` bigint unsigned NOT NULL,
  `was_exit` tinyint NOT NULL,
  PRIMARY KEY (`failtime_id`),
  KEY `beatmapid_index` (`beatmap_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `friends`
--

DROP TABLE IF EXISTS `friends`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `friends` (
  `user_1` bigint unsigned NOT NULL,
  `user_2` bigint unsigned NOT NULL,
  PRIMARY KEY (`user_1`,`user_2`),
  KEY `index_user1` (`user_1`) /*!80000 INVISIBLE */,
  KEY `index_user2` (`user_2`),
  CONSTRAINT `user_id2_FK` FOREIGN KEY (`user_2`) REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_id_FK` FOREIGN KEY (`user_1`) REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `irc_log`
--

DROP TABLE IF EXISTS `irc_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `irc_log` (
  `message_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `sender` bigint unsigned NOT NULL,
  `target` varchar(64) NOT NULL,
  `message` text NOT NULL,
  `date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`message_id`)
) ENGINE=InnoDB AUTO_INCREMENT=91 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `scores`
--

DROP TABLE IF EXISTS `scores`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `scores` (
  `score_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `beatmap_id` int NOT NULL,
  `beatmapset_id` int NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `playmode` tinyint NOT NULL,
  `score` int NOT NULL DEFAULT '0',
  `max_combo` int NOT NULL DEFAULT '0',
  `ranking` varchar(2) NOT NULL,
  `hit300` int NOT NULL DEFAULT '0',
  `hit100` int NOT NULL DEFAULT '0',
  `hit50` int NOT NULL DEFAULT '0',
  `hitMiss` int NOT NULL DEFAULT '0',
  `hitGeki` int NOT NULL DEFAULT '0',
  `hitKatu` int NOT NULL DEFAULT '0',
  `enabled_mods` int NOT NULL DEFAULT '0',
  `perfect` tinyint NOT NULL DEFAULT '0',
  `passed` tinyint NOT NULL DEFAULT '0',
  `date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `leaderboard_best` tinyint NOT NULL DEFAULT '0',
  `mapset_best` tinyint NOT NULL DEFAULT '0',
  `score_hash` varchar(64) NOT NULL,
  PRIMARY KEY (`score_id`),
  KEY `userid_index` (`score_id`) /*!80000 INVISIBLE */,
  KEY `leaderboardbest_index` (`leaderboard_best`) /*!80000 INVISIBLE */,
  KEY `mapsetbest_index` (`mapset_best`) /*!80000 INVISIBLE */,
  KEY `scorehash_index` (`score_hash`) /*!80000 INVISIBLE */,
  KEY `beatmapid_index` (`beatmap_id`),
  KEY `userid_fk_idx` (`user_id`),
  CONSTRAINT `userid_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=85 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `screenshots`
--

DROP TABLE IF EXISTS `screenshots`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `screenshots` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `filename` varchar(256) NOT NULL,
  `creation_date` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `stats`
--

DROP TABLE IF EXISTS `stats`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `stats` (
  `user_id` bigint unsigned NOT NULL,
  `mode` tinyint NOT NULL,
  `ranked_score` bigint unsigned NOT NULL DEFAULT '0',
  `total_score` bigint unsigned NOT NULL DEFAULT '0',
  `user_level` double unsigned NOT NULL DEFAULT '0',
  `accuracy` float unsigned NOT NULL DEFAULT '0',
  `playcount` bigint unsigned NOT NULL DEFAULT '0',
  `count_ssh` bigint unsigned NOT NULL DEFAULT '0',
  `count_ss` bigint unsigned NOT NULL DEFAULT '0',
  `count_sh` bigint unsigned NOT NULL DEFAULT '0',
  `count_s` bigint unsigned NOT NULL DEFAULT '0',
  `count_a` bigint unsigned NOT NULL DEFAULT '0',
  `count_b` bigint unsigned NOT NULL DEFAULT '0',
  `count_c` bigint unsigned NOT NULL DEFAULT '0',
  `count_d` bigint unsigned NOT NULL DEFAULT '0',
  `hit300` bigint unsigned NOT NULL DEFAULT '0',
  `hit100` bigint unsigned NOT NULL DEFAULT '0',
  `hit50` bigint unsigned NOT NULL DEFAULT '0',
  `hitMiss` bigint unsigned NOT NULL DEFAULT '0',
  `hitGeki` bigint unsigned NOT NULL DEFAULT '0',
  `hitKatu` bigint unsigned NOT NULL DEFAULT '0',
  `replays_watched` bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`user_id`,`mode`),
  CONSTRAINT `userid` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `updater_items`
--

DROP TABLE IF EXISTS `updater_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `updater_items` (
  `item_id` bigint NOT NULL AUTO_INCREMENT,
  `server_filename` varchar(256) NOT NULL,
  `client_filename` varchar(256) NOT NULL,
  `file_hash` varchar(64) NOT NULL,
  `item_name` varchar(256) NOT NULL DEFAULT '',
  `item_action` varchar(8) DEFAULT 'none',
  PRIMARY KEY (`item_id`,`server_filename`,`client_filename`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `user_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(32) NOT NULL,
  `password` varchar(64) NOT NULL,
  `country` smallint unsigned NOT NULL DEFAULT '0',
  `banned` tinyint NOT NULL DEFAULT '0',
  `banned_reason` varchar(256) NOT NULL DEFAULT 'no reason',
  `privileges` int NOT NULL DEFAULT '1',
  `joined_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`,`username`),
  UNIQUE KEY `id_UNIQUE` (`user_id`),
  UNIQUE KEY `username_UNIQUE` (`username`),
  KEY `user_INDEX` (`username`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-05-26 17:29:57
