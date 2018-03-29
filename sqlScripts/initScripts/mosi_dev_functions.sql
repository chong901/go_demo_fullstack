-- MySQL dump 10.13  Distrib 5.7.17, for macos10.12 (x86_64)
--
-- Host: 127.0.0.1    Database: mosi_dev
-- ------------------------------------------------------
-- Server version	5.7.19

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `functions`
--

DROP TABLE IF EXISTS `functions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `functions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `parent` int(11) DEFAULT NULL,
  `uri` varchar(50) DEFAULT NULL,
  `method` varchar(10) DEFAULT NULL,
  `name` varchar(50) DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `order_num` int(11) DEFAULT NULL,
  `is_menu` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_functions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=68 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `functions`
--

LOCK TABLES `functions` WRITE;
/*!40000 ALTER TABLE `functions` DISABLE KEYS */;
INSERT INTO `functions` VALUES (1,0,'','','menu',NULL,NULL,'2017-09-18 16:29:35',0,1),(2,1,'/poM','GET','Open Orders',NULL,NULL,'2017-09-18 16:26:31',1,1),(3,35,'/skuM','GET','Part Management',NULL,NULL,'2017-09-18 16:27:28',1,1),(4,35,'/jobM','GET','Job Management',NULL,NULL,'2017-09-18 16:27:32',2,1),(5,1,'/inventory','GET','Inventory',NULL,NULL,'2017-09-18 16:26:48',2,1),(6,1,'','','Admin',NULL,NULL,'2017-09-18 17:00:42',99,1),(7,6,'/function','GET','Function',NULL,NULL,'2017-09-18 16:30:33',3,1),(8,6,'/role','GET','Role',NULL,NULL,'2017-09-18 16:30:29',2,1),(9,6,'/account','GET','Account',NULL,NULL,'2017-09-18 16:30:26',1,1),(10,0,'/welcome','GET','welcome',NULL,NULL,'2017-09-18 16:29:40',0,1),(11,4,'/skuM/save','GET','edit produce',NULL,NULL,NULL,NULL,NULL),(12,4,'/skuM/save','POST','edit produce',NULL,NULL,NULL,NULL,NULL),(13,4,'/skuM','DELETE','delete produce',NULL,NULL,NULL,NULL,NULL),(14,5,'/inventory','PUT','update inventory',NULL,NULL,'2017-09-18 16:00:42',0,NULL),(15,5,'/inventory','DELETE','delete inventory','2017-09-18 18:33:30',NULL,'2017-09-18 16:00:48',0,NULL),(16,5,'/inventory/history','GET','inventory detail',NULL,NULL,'2017-09-18 16:57:59',0,0),(17,3,'/jobM/save','GET','edit job',NULL,NULL,NULL,NULL,NULL),(18,3,'/jobM/save','POST','edit job',NULL,NULL,NULL,NULL,NULL),(19,3,'/jobM','DELETE','delete job',NULL,NULL,NULL,NULL,NULL),(24,9,'/account/role','PUT','edit account role',NULL,NULL,NULL,NULL,NULL),(25,7,'/function','POST','edit function',NULL,NULL,NULL,NULL,NULL),(26,7,'/function','DELETE','delete function',NULL,NULL,NULL,NULL,NULL),(27,8,'/role','POST','create role',NULL,NULL,NULL,NULL,NULL),(28,8,'/role','PUT','update role','2017-09-19 09:52:56',NULL,NULL,NULL,NULL),(29,2,'/poM/save','GET','edit po',NULL,NULL,'2017-09-18 15:59:52',0,NULL),(30,2,'/poM/save','POST','edit po',NULL,NULL,'2017-09-18 15:59:59',0,NULL),(31,2,'/poM','DELETE','delete po',NULL,NULL,'2017-09-18 16:00:05',0,NULL),(32,9,'/account','POST','add account',NULL,NULL,NULL,NULL,NULL),(33,9,'/account/password','PUT','change password',NULL,NULL,'2017-09-18 16:58:21',0,0),(34,35,'/jobDash','GET','Job Dashboard','2017-09-19 09:53:19',NULL,'2017-09-18 16:27:37',3,1),(35,1,'','','Production',NULL,NULL,'2017-09-18 16:26:53',3,1),(36,35,'/machineM','GET','Machine Management',NULL,NULL,'2017-09-18 16:27:44',4,1),(37,36,'/machineM/save','GET','edit machine',NULL,NULL,NULL,NULL,NULL),(38,36,'/machineM/save','POST','edit machine',NULL,NULL,NULL,NULL,NULL),(39,36,'/machineM','DELETE','delete machine',NULL,NULL,NULL,NULL,NULL),(40,1,'','','Configuration',NULL,'2017-09-15 16:44:52','2017-09-21 11:58:38',5,1),(41,40,'/configuration/skuParameters','GET','Parameter Setting',NULL,'2017-09-15 16:45:59','2017-09-18 16:27:50',0,1),(42,40,'/configuration/unit','GET','Unit Setting',NULL,'2017-09-15 16:46:56','2017-09-18 16:27:54',0,1),(43,40,'/configuration','POST','Add configuration',NULL,'2017-09-15 16:47:41','2017-09-18 16:58:43',0,0),(44,40,'/configuration','DELETE','Delete Configuration',NULL,'2017-09-15 16:47:54','2017-09-18 16:58:46',0,0),(45,1,'','','test','2017-09-18 11:24:04','2017-09-18 11:16:11','2017-09-18 11:22:44',0,NULL),(46,45,'/test1','GET','test1','2017-09-18 11:23:59','2017-09-18 11:17:56','2017-09-18 11:21:56',4,NULL),(47,45,'/test2','GET','test2','2017-09-18 11:23:52','2017-09-18 11:18:09','2017-09-18 11:19:33',2,NULL),(48,45,'/test3','GET','test3','2017-09-18 11:23:56','2017-09-18 11:18:20','2017-09-18 11:21:32',3,NULL),(49,0,'','','API',NULL,'2017-09-18 12:45:06','2017-09-18 12:45:14',0,NULL),(50,49,'/api/v1/role/all','GET','Get All Role',NULL,'2017-09-18 12:45:51','2017-09-18 12:45:51',0,NULL),(51,49,'/api/v1/role','GET','Get Role',NULL,'2017-09-18 12:46:19','2017-09-18 12:46:19',0,NULL),(52,3,'/skuM/uploadFile','POST','Sku Upload File',NULL,'2017-09-18 16:48:49','2017-09-18 16:57:40',0,0),(53,3,'/skuParameter','DELETE','Delete Sku Parameter',NULL,'2017-09-18 16:49:17','2017-09-18 16:57:46',0,0),(54,2,'/poSku','DELETE','Delete Po Sku',NULL,'2017-09-18 16:50:35','2017-09-18 16:57:29',0,0),(55,36,'/machineM/save','PUT','Update Machine',NULL,'2017-09-18 16:51:54','2017-09-18 16:58:33',0,0),(56,8,'/role/name','PUT','Update Role Name',NULL,'2017-09-18 16:52:49','2017-09-18 16:58:13',0,0),(57,8,'/role/functions','PUT','Update Role Functions',NULL,'2017-09-18 16:53:17','2017-09-18 16:58:16',0,0),(58,0,'','','','2017-09-19 10:06:03','2017-09-19 10:05:58','2017-09-19 10:05:58',0,0),(59,0,'','','','2017-09-19 10:06:44','2017-09-19 10:06:40','2017-09-19 10:06:40',0,0),(60,8,'/role','DELETE','Delete Role',NULL,'2017-09-19 11:54:16','2017-09-19 11:54:16',0,0),(61,1,'','','ID Rule',NULL,'2017-09-21 11:53:13','2017-09-21 11:58:29',4,1),(62,61,'/jobIdRuleList','GET','Job ID Rule',NULL,'2017-09-21 11:54:11','2017-09-21 11:54:11',0,1),(63,62,'/jobIdRule','GET','Job ID Edit View',NULL,'2017-09-21 11:54:59','2017-09-21 11:54:59',0,0),(64,62,'/jobIdRule','DELETE','Delete Job ID Rule',NULL,'2017-09-21 11:55:31','2017-09-21 11:55:31',0,0),(65,0,'','','ID Rule',NULL,'2017-09-21 11:55:54','2017-09-21 11:55:54',0,0),(66,65,'/idRule','POST','Add ID Rule',NULL,'2017-09-21 11:56:22','2017-09-21 11:56:22',0,0),(67,65,'/idRule','PUT','Update ID Rule',NULL,'2017-09-21 11:56:45','2017-09-21 11:56:45',0,0);
/*!40000 ALTER TABLE `functions` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-09-21 15:04:17
