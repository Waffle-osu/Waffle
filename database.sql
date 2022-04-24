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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friends`
--

LOCK TABLES `friends` WRITE;
/*!40000 ALTER TABLE `friends` DISABLE KEYS */;
INSERT INTO `friends` VALUES (3,2),(2,3);
/*!40000 ALTER TABLE `friends` ENABLE KEYS */;
UNLOCK TABLES;

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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `stats`
--

LOCK TABLES `stats` WRITE;
/*!40000 ALTER TABLE `stats` DISABLE KEYS */;
INSERT INTO `stats` VALUES (2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(2,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(2,3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(3,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(3,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(3,3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(4,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(4,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(4,3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0);
/*!40000 ALTER TABLE `stats` ENABLE KEYS */;
UNLOCK TABLES;

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
  `privileges` int NOT NULL DEFAULT '0',
  `joined_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`,`username`),
  UNIQUE KEY `id_UNIQUE` (`user_id`),
  UNIQUE KEY `username_UNIQUE` (`username`),
  KEY `user_INDEX` (`username`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (2,'Furball','1787d7646304c5d987cf4e64a3973dc7',0,0,'no reason',0,'2022-04-20 22:26:46'),(3,'Eevee','1787d7646304c5d987cf4e64a3973dc7',0,0,'no reason',0,'2022-04-21 20:36:26'),(4,'marshall_racer','1787d7646304c5d987cf4e64a3973dc7',0,0,'no reason',0,'2022-04-22 23:39:09');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-04-24 16:15:54
