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
INSERT INTO `stats` VALUES (1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(1,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(1,3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(2,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(2,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(2,3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(3,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(3,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(3,3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(4,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(4,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),(4,3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0);
/*!40000 ALTER TABLE `stats` ENABLE KEYS */;
UNLOCK TABLES;

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
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `updater_items`
--

LOCK TABLES `updater_items` WRITE;
/*!40000 ALTER TABLE `updater_items` DISABLE KEYS */;
INSERT INTO `updater_items` VALUES (1,'osu!.exe','osu!.exe','3623d9f7c693b786564e2d61b1c43af9','client_debug','none'),(2,'avcodec-51.dll','avcodec-51.dll','b22bf1e4ecd4be3d909dc68ccab74eec','client_ddls','none'),(3,'avformat-52.dll','avformat-52.dll','2e7a800133625f827cf46aa0bb1af800','client_ddls','none'),(4,'avutil-49.dll','avutil-49.dll','c870147dff89c95c81f8fbdfbc6344ac','client_ddls','none'),(5,'bass.dll','bass.dll','bbfc7d855252b0211875769bbf667bcd','client_ddls','none'),(6,'bass_fx.dll','bass_fx.dll','f9ffe0a23a32b79653e31330764c4231','client_ddls','none'),(7,'d3dx9_31.dll','d3dx9_31.dll','797e24743937d67d69f28f2cf5052ee8','client_ddls','none'),(8,'Microsoft.Ink.dll','Microsoft.Ink.dll','a02ee61542caae25f8a44c9428d30247','client_ddls','none'),(9,'Microsoft.Xna.Framework.dll','Microsoft.Xna.Framework.dll','45a786658d3f69717652fed471d03ee0','client_ddls','none'),(10,'osu!common.dll','osu!common.dll','820817b776374a0adcbb7fa2a7ca74f2','client_ddls','none'),(11,'osu.dll.zip','osu.dll','599c14bcfc9c43b88d70d1a9388b33b7','client_ddls','zip'),(12,'OsuP2P.dll','OsuP2P.dll','2342bfd835e2e487d040d8ba62eb1a72','client_ddls','none'),(13,'pthreadGC2.dll','pthreadGC2.dll','ce931021e18f385f519e945a8a10548e','client_ddls','none'),(14,'x3daudio1_1.dll','x3daudio1_1.dll','121b131eaa369d8f58dacc5c39a77d80','client_ddls','none');
/*!40000 ALTER TABLE `updater_items` ENABLE KEYS */;
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
  `privileges` int NOT NULL DEFAULT '1',
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
INSERT INTO `users` VALUES (1,'WaffleBot','No!',0,0,'no reason',31,'2022-04-24 16:40:27'),(2,'Furball','$2a$10$f3Q8GKPnKffV4K3n.9cf1.qpJYoWbKAzD17LHRbGe2x1Nal5EpYFa',0,0,'no reason',31,'2022-04-20 22:26:46'),(3,'Eevee','$2a$10$oTu9vvSE.xaEt2OPrkuIM.JiNGH.U60DBgYh4mIH4JYCRYnrfzLB.',0,0,'no reason',1,'2022-04-21 20:36:26'),(4,'ArNeN','$2a$10$YggUilDeQa7Tl.TwILijTuRLmdS6ndjSxAizqrFrocriL3V118Mem',0,0,'no reason',1,'2022-04-24 18:33:15');
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

-- Dump completed on 2022-05-02 23:52:47
