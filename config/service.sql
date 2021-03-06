-- MySQL dump 10.13  Distrib 5.7.26, for Win64 (x86_64)
--
-- Host: localhost    Database: service
-- ------------------------------------------------------
-- Server version	5.7.26

/*!40101 SET @OLD_CHARACTER_SET_CLIENT = @@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS = @@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION = @@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE = @@TIME_ZONE */;
/*!40103 SET TIME_ZONE = '+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS = @@UNIQUE_CHECKS, UNIQUE_CHECKS = 0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS = 0 */;
/*!40101 SET @OLD_SQL_MODE = @@SQL_MODE, SQL_MODE = 'NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES = @@SQL_NOTES, SQL_NOTES = 0 */;

--
-- Table structure for table `advisor`
--

DROP TABLE IF EXISTS `advisor`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `advisor`
(
    `id`                int(11)                              NOT NULL AUTO_INCREMENT,
    `phone`             varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `password`          varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `name`              varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `coin`              int(2) unsigned zerofill             DEFAULT '00',
    `total_order_num`   int(1) unsigned zerofill             DEFAULT '0',
    `status`            int(1) unsigned zerofill             DEFAULT '0',
    `rank`              float unsigned zerofill              DEFAULT '000000000000',
    `rank_num`          int(1) unsigned zerofill             DEFAULT '0',
    `work_experience`   int(1) unsigned zerofill             DEFAULT '0',
    `bio`               varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `about`             varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `on_time`           float                                DEFAULT NULL,
    `total_comment_num` int(11)                              DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 30004
  DEFAULT CHARSET = utf8
  COLLATE = utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `advisor`
--

LOCK TABLES `advisor` WRITE;
/*!40000 ALTER TABLE `advisor`
    DISABLE KEYS */;
INSERT INTO `advisor`
VALUES (30001, '17633333333', '$2a$10$nxlweDGYboxQKGrpiXyVWeSxbYFqdFh/MCtZ2heGv97X5ZDGhdP8S', 'Boooob', 20630, 7, 1,
        0000000002.2, 0, 4, 'BioBioBioBiobio', 'this is huiofda fdashuif fhsduiafas', 0.714286, 5),
       (30002, '17633331234', '$2a$10$WUix5JYkjXXQa4JNsoPE4OCeJQ4JZJ0hdR/6F8eQJSqRrmjSBsnI2', '鏂扮殑椤鹃棶Luke', 00, 0, 1,
        000000000000, 0, 2, 'BioBioBioBiobio', 'this is huiofda fdashuif fhsduiafas', 0, 0),
       (30003, '17633339999', '$2a$10$ss5TJR5.efKhuexdB5XWxeFhvMb1HdQzA2o8sSiBJnmeWSXzSfXJu', NULL, 800, 0, 0,
        000000000000, 0, 0, NULL, NULL, 0, 0);
/*!40000 ALTER TABLE `advisor`
    ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bill`
--

DROP TABLE IF EXISTS `bill`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bill`
(
    `id`         int(2) NOT NULL AUTO_INCREMENT,
    `user_id`    int(2) NOT NULL,
    `advisor_id` int(2) NOT NULL,
    `order_id`   int(2) NOT NULL,
    `amount`     int(2) NOT NULL,
    `type`       int(2) NOT NULL,
    `time`       int(2) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = MyISAM
  AUTO_INCREMENT = 18
  DEFAULT CHARSET = utf8
  COLLATE = utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bill`
--

LOCK TABLES `bill` WRITE;
/*!40000 ALTER TABLE `bill`
    DISABLE KEYS */;
INSERT INTO `bill`
VALUES (1, 3, 0, 1, 5000, 1, 1654755050),
       (2, 3, 0, 2, 150, 1, 1654755054),
       (3, 3, 0, 3, 330, 1, 1654755060),
       (4, 3, 0, 3, 165, 3, 1654755073),
       (5, 0, 30001, 1, 5000, 5, 1654755546),
       (6, 3, 0, 3, 165, 4, 1654755722),
       (7, 3, 0, 2, 150, 2, 1654759670),
       (8, 3, 0, 3, 330, 2, 1654759670),
       (9, 3, 0, 4, 5000, 1, 1654767137),
       (10, 0, 30001, 4, 5000, 5, 1654767143),
       (11, 3, 0, 5, 330, 1, 1654767260),
       (12, 3, 0, 5, 165, 3, 1654767278),
       (13, 0, 30001, 5, 495, 5, 1654767282),
       (14, 3, 0, 6, 150, 1, 1654767374),
       (15, 0, 30001, 6, 150, 5, 1654767380),
       (16, 3, 0, 7, 150, 1, 1654767413),
       (17, 0, 30001, 7, 150, 5, 1654767422);
/*!40000 ALTER TABLE `bill`
    ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `collection`
--

DROP TABLE IF EXISTS `collection`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `collection`
(
    `id`         int(2) NOT NULL AUTO_INCREMENT,
    `user_id`    int(2) NOT NULL,
    `advisor_id` int(2) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = MyISAM
  AUTO_INCREMENT = 4
  DEFAULT CHARSET = utf8
  COLLATE = utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `collection`
--

LOCK TABLES `collection` WRITE;
/*!40000 ALTER TABLE `collection`
    DISABLE KEYS */;
INSERT INTO `collection`
VALUES (1, 3, 30001),
       (2, 3, 30002),
       (3, 3, 30003);
/*!40000 ALTER TABLE `collection`
    ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `service`
--

DROP TABLE IF EXISTS `service`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `service`
(
    `id`              int(1) NOT NULL AUTO_INCREMENT,
    `advisor_id`      int(1) NOT NULL,
    `service_name_id` int(1) NOT NULL,
    `service_name`    varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `price`           int(2)                               DEFAULT NULL,
    `status`          int(1)                               DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 13
  DEFAULT CHARSET = utf8
  COLLATE = utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `service`
--

LOCK TABLES `service` WRITE;
/*!40000 ALTER TABLE `service`
    DISABLE KEYS */;
INSERT INTO `service`
VALUES (1, 30001, 1, '24h Delivered Video Reading', 330, 1),
       (2, 30001, 2, '24h Delivered Audio Reading', 5000, 1),
       (3, 30001, 3, '24h Delivered Text Reading', 150, 1),
       (4, 30001, 4, 'Live Text Chat', 300, 0),
       (5, 30002, 4, 'Live Text Chat', 120, 0),
       (6, 30002, 1, '24h Delivered Video Reading', 200, 1),
       (7, 30002, 2, '24h Delivered Audio Reading', 230, 1),
       (8, 30002, 3, '24h Delivered Text Reading', 150, 1),
       (9, 30003, 1, '24h Delivered Video Reading', 100, 0),
       (10, 30003, 2, '24h Delivered Audio Reading', 100, 0),
       (11, 30003, 3, '24h Delivered Text Reading', 100, 0),
       (12, 30003, 4, 'Live Text Chat', 100, 0);
/*!40000 ALTER TABLE `service`
    ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user`
(
    `id`       int(11)                              NOT NULL AUTO_INCREMENT,
    `name`     varchar(255) COLLATE utf8_unicode_ci DEFAULT '',
    `birth`    varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `gender`   int(1) unsigned zerofill             DEFAULT '0',
    `bio`      varchar(60) COLLATE utf8_unicode_ci  DEFAULT NULL,
    `about`    varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
    `coin`     int(2) unsigned zerofill             DEFAULT '00',
    `phone`    varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `password` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 6
  DEFAULT CHARSET = utf8
  COLLATE = utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user`
    DISABLE KEYS */;
INSERT INTO `user`
VALUES (1, '鍛靛懙1231', '09-01-2010', 1, '123afdsf', '123', 11700, '17600002222',
        '$2a$10$0aZEyKiEpBdVQ2II5DWZv.qQqgNORvQqdHoiJKbIrNCWZktgHWw8e'),
       (2, '', NULL, 0, NULL, NULL, 00, '17600004444', '$2a$10$ExwgFYXrh9x83dPT0AvL3OZzCmdLOz7tA5CU3VxgI.thf2wmFiHby'),
       (3, '鐢ㄦ埛鍚?', '01-11-1999', 1, 'this is bio', '1238912839 128931', 489205, '17607175592',
        '$2a$10$qwYLtsdXHN/NX9qhF3.Iye5oQ3p/GP0iNAEWuXT.zI1FXqcs49jcC'),
       (4, '', NULL, 0, NULL, NULL, 00, '17600000000', '$2a$10$lmKv43U0L1Mtej2rWvk8/uJPeu8.w1wrTAe2aL.7D.CeTgnlwrO1i'),
       (5, '', NULL, 0, NULL, NULL, 00, '17666661111', '$2a$10$yHjsMbLkFLGJXwqUMmeQm.tpuxOL98Ryvu2V3HGQSj2jJai.g/g7S');
/*!40000 ALTER TABLE `user`
    ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_order`
--

DROP TABLE IF EXISTS `user_order`;
/*!40101 SET @saved_cs_client = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_order`
(
    `id`              int(2) NOT NULL AUTO_INCREMENT,
    `user_id`         int(2) NOT NULL,
    `service_id`      int(2) NOT NULL,
    `advisor_id`      int(2) NOT NULL,
    `service_name_id` int(2)                            DEFAULT NULL,
    `rush_coin`       int(2) NOT NULL,
    `coin`            int(2) NOT NULL,
    `situation`       varchar(9000) CHARACTER SET utf8  DEFAULT NULL,
    `question`        varchar(600) CHARACTER SET utf8   DEFAULT NULL,
    `rush_time`       int(1) NOT NULL,
    `status`          int(1)                            DEFAULT NULL,
    `create_time`     int(1) NOT NULL,
    `reply`           varchar(11000) CHARACTER SET utf8 DEFAULT NULL,
    `rate`            int(1)                            DEFAULT NULL,
    `comment`         varchar(255) CHARACTER SET utf8   DEFAULT NULL,
    `comment_time`    int(1)                            DEFAULT NULL,
    `comment_status`  int(1)                            DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 8
  DEFAULT CHARSET = utf8
  COLLATE = utf8_unicode_ci
  ROW_FORMAT = DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_order`
--

LOCK TABLES `user_order` WRITE;
/*!40000 ALTER TABLE `user_order`
    DISABLE KEYS */;
INSERT INTO `user_order`
VALUES (1, 3, 2, 30001, 2, 2500, 5000, '闂鎻忚堪锛?,'璁㈠崟1',0,3,1654755050,'Lorem ipsum dolor sit amet,
        consectetur adipiscing elit. Pellentesque nec purus pharetra, elementum augue at,
        semper eros. Nam lacinia tincidunt turpis vel viverra. Aliquam euismod orci eu consectetur dapibus. Aenean at justo interdum,
        pulvinar odio id, vulputate est. Duis non ultrices neque. Curabitur porta egestas nunc,
        non eleifend turpis auctor in. Nullam accumsan diam non elit molestie laoreet. Quisque cursus pretium bibendum.Curabitur varius egestas lectus,
        et luctus turpis eleifend id. Lorem ipsum dolor sit amet,
        consectetur adipiscing elit. Proin sed nisi a lacus efficitur interdum. Morbi gravida consequat tempor. Integer volutpat tellus at elit elementum,
        ac mattis lorem porttitor. Ut quis ipsum non erat pharetra pulvinar tincidunt sit amet elit. Praesent gravida neque turpis,
        vel hendrerit tellus elementum nec. Morbi eros orci, aliquet a venenatis eu, ornare in neque. Etiam eleifend,
        ligula nec elementum vulputate, quam magna bibendum justo, vel consequat erat elit at quam. Curabitur sagittis,
        lectus nec cursus porta, purus erat auctor est,
        ut varius est odio ut dui. In mollis orci non luctus blandit. Donec maximus rhoncus metus, in feugiat urna
        consectetur vitae integer.',5,' good good ',1654767021,1),(2,3,3,30001,3,75,150,' 闂鎻忚堪锛 ?, '璁㈠崟2', 0, 2,
        1654755054, '', 0, '', 0, 0),
       (3, 3, 1, 30001, 1, 165, 330,
        '闂鎻忚堪锛?,'璁㈠崟3',1654755073,2,1654755060,'',0,'',0,0),(4,3,2,30001,2,2500,5000,'闂鎻忚堪锛?, '璁㈠崟3', 0, 3,
        1654767137,
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque nec purus pharetra, elementum augue at, semper eros. Nam lacinia tincidunt turpis vel viverra. Aliquam euismod orci eu consectetur dapibus. Aenean at justo interdum, pulvinar odio id, vulputate est. Duis non ultrices neque. Curabitur porta egestas nunc, non eleifend turpis auctor in. Nullam accumsan diam non elit molestie laoreet. Quisque cursus pretium bibendum.Curabitur varius egestas lectus, et luctus turpis eleifend id. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin sed nisi a lacus efficitur interdum. Morbi gravida consequat tempor. Integer volutpat tellus at elit elementum, ac mattis lorem porttitor. Ut quis ipsum non erat pharetra pulvinar tincidunt sit amet elit. Praesent gravida neque turpis, vel hendrerit tellus elementum nec. Morbi eros orci, aliquet a venenatis eu, ornare in neque. Etiam eleifend, ligula nec elementum vulputate, quam magna bibendum justo, vel consequat erat elit at quam. Curabitur sagittis, lectus nec cursus porta, purus erat auctor est, ut varius est odio ut dui. In mollis orci non luctus blandit. Donec maximus rhoncus metus, in feugiat urna consectetur vitae integer.',
        3, 'just so so', 1654767157, 1),
       (5, 3, 1, 30001, 1, 165, 330, '闂鎻忚堪锛?,'鏈変釜鏂伴棶棰?, 1654767278, 3, 1654767260,
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque nec purus pharetra, elementum augue at, semper eros. Nam lacinia tincidunt turpis vel viverra. Aliquam euismod orci eu consectetur dapibus. Aenean at justo interdum, pulvinar odio id, vulputate est. Duis non ultrices neque. Curabitur porta egestas nunc, non eleifend turpis auctor in. Nullam accumsan diam non elit molestie laoreet. Quisque cursus pretium bibendum.Curabitur varius egestas lectus, et luctus turpis eleifend id. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin sed nisi a lacus efficitur interdum. Morbi gravida consequat tempor. Integer volutpat tellus at elit elementum, ac mattis lorem porttitor. Ut quis ipsum non erat pharetra pulvinar tincidunt sit amet elit. Praesent gravida neque turpis, vel hendrerit tellus elementum nec. Morbi eros orci, aliquet a venenatis eu, ornare in neque. Etiam eleifend, ligula nec elementum vulputate, quam magna bibendum justo, vel consequat erat elit at quam. Curabitur sagittis, lectus nec cursus porta, purus erat auctor est, ut varius est odio ut dui. In mollis orci non luctus blandit. Donec maximus rhoncus metus, in feugiat urna consectetur vitae integer.',
        1, 'bad experience', 1654767301, 1),
       (6, 3, 3, 30001, 3, 75, 150, '闂鎻忚堪锛?,'鏈変釜鏂伴棶棰?, 0, 3, 1654767374,
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque nec purus pharetra, elementum augue at, semper eros. Nam lacinia tincidunt turpis vel viverra. Aliquam euismod orci eu consectetur dapibus. Aenean at justo interdum, pulvinar odio id, vulputate est. Duis non ultrices neque. Curabitur porta egestas nunc, non eleifend turpis auctor in. Nullam accumsan diam non elit molestie laoreet. Quisque cursus pretium bibendum.Curabitur varius egestas lectus, et luctus turpis eleifend id. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin sed nisi a lacus efficitur interdum. Morbi gravida consequat tempor. Integer volutpat tellus at elit elementum, ac mattis lorem porttitor. Ut quis ipsum non erat pharetra pulvinar tincidunt sit amet elit. Praesent gravida neque turpis, vel hendrerit tellus elementum nec. Morbi eros orci, aliquet a venenatis eu, ornare in neque. Etiam eleifend, ligula nec elementum vulputate, quam magna bibendum justo, vel consequat erat elit at quam. Curabitur sagittis, lectus nec cursus porta, purus erat auctor est, ut varius est odio ut dui. In mollis orci non luctus blandit. Donec maximus rhoncus metus, in feugiat urna consectetur vitae integer.',
        1, 'bad experience', 1654767385, 1),
       (7, 3, 3, 30001, 3, 75, 150, '闂鎻忚堪锛?,'鏈変釜鏂伴棶棰?, 0, 3, 1654767413,
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque nec purus pharetra, elementum augue at, semper eros. Nam lacinia tincidunt turpis vel viverra. Aliquam euismod orci eu consectetur dapibus. Aenean at justo interdum, pulvinar odio id, vulputate est. Duis non ultrices neque. Curabitur porta egestas nunc, non eleifend turpis auctor in. Nullam accumsan diam non elit molestie laoreet. Quisque cursus pretium bibendum.Curabitur varius egestas lectus, et luctus turpis eleifend id. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin sed nisi a lacus efficitur interdum. Morbi gravida consequat tempor. Integer volutpat tellus at elit elementum, ac mattis lorem porttitor. Ut quis ipsum non erat pharetra pulvinar tincidunt sit amet elit. Praesent gravida neque turpis, vel hendrerit tellus elementum nec. Morbi eros orci, aliquet a venenatis eu, ornare in neque. Etiam eleifend, ligula nec elementum vulputate, quam magna bibendum justo, vel consequat erat elit at quam. Curabitur sagittis, lectus nec cursus porta, purus erat auctor est, ut varius est odio ut dui. In mollis orci non luctus blandit. Donec maximus rhoncus metus, in feugiat urna consectetur vitae integer.',
        1, '1111', 1654767678, 1);
/*!40000 ALTER TABLE `user_order`
    ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE = @OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE = @OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS = @OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS = @OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT = @OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS = @OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION = @OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES = @OLD_SQL_NOTES */;

-- Dump completed on 2022-06-10 14:20:29
