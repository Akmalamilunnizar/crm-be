-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               8.4.3 - MySQL Community Server - GPL
-- Server OS:                    Win64
-- HeidiSQL Version:             12.8.0.6908
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- Dumping database structure for iqgncnzy_skripsi
CREATE DATABASE IF NOT EXISTS `iqgncnzy_skripsi` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `iqgncnzy_skripsi`;

-- Dumping structure for table iqgncnzy_skripsi.accounts
CREATE TABLE IF NOT EXISTS `accounts` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `saldo` int NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.accounts: ~4 rows (approximately)
INSERT INTO `accounts` (`id`, `name`, `saldo`, `createdAt`, `updatedAt`) VALUES
	('b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'Lilly', 1345000, '2025-05-17 00:11:18.045', '2025-08-14 15:25:46.559'),
	('bank', 'Bank Account', 0, '2025-08-30 10:33:01.849', '2025-08-30 10:33:01.000'),
	('cash', 'Cash Account', 0, '2025-08-30 10:33:01.849', '2025-08-30 10:33:01.000'),
	('receivables', 'Accounts Receivable', 0, '2025-08-30 10:33:01.849', '2025-08-30 10:33:01.000');

-- Dumping structure for table iqgncnzy_skripsi.areas
CREATE TABLE IF NOT EXISTS `areas` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name_city` longtext COLLATE utf8mb4_unicode_ci,
  `name_subdistrict` longtext COLLATE utf8mb4_unicode_ci,
  `name_village` longtext COLLATE utf8mb4_unicode_ci,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.areas: ~11 rows (approximately)
INSERT INTO `areas` (`id`, `name_city`, `name_subdistrict`, `name_village`, `createdAt`, `updatedAt`) VALUES
	('2f4e724d-39ab-4c8e-a5f9-437e100f7b12', 'Jakarta Selatan', 'Kebayoran Baru', 'Senayan', '2025-08-22 10:43:34.390', '2025-08-22 10:43:34.390'),
	('61362c06-7042-4ed4-a391-237b742912c3', 'MLG', 'TURN', 'TLOK', '2025-05-10 12:27:01.992', '2025-06-02 13:22:50.602'),
	('acd9f1ea-6c68-4522-916a-3586d368ec70', 'asdd', 'asd2', 'asd2', '2025-06-03 05:49:28.396', '2025-06-02 22:49:29.025'),
	('b4d33f98-7ef4-43f6-8c97-9fa7854e5988', 'MLG', 'KDKG', 'SWJJR', '2025-06-17 15:32:08.958', '2025-06-17 08:32:09.161'),
	('b53e9371-036f-4dcb-88b8-0a3679789e76', 'MLG', 'TURN', 'SDYU', '2025-05-13 10:52:51.471', '2025-05-13 03:52:51.839'),
	('ee888c0f-c3d7-4cc1-aacb-c153d52d3e1b', 'Jakarta Pusat', 'Menteng', 'Menteng', '2025-08-22 10:43:34.385', '2025-08-22 10:43:34.384'),
	('jakarta-barat', 'Jakarta Barat', 'Grogol', 'Grogol', '2025-08-30 10:33:01.844', '2025-08-30 10:33:01.000'),
	('jakarta-pusat', 'Jakarta Pusat', 'Menteng', 'Menteng', '2025-08-30 10:33:01.844', '2025-08-30 10:33:01.000'),
	('jakarta-selatan', 'Jakarta Selatan', 'Kebayoran Baru', 'Kebayoran Baru', '2025-08-30 10:33:01.844', '2025-08-30 10:33:01.000'),
	('jakarta-timur', 'Jakarta Timur', 'Jatinegara', 'Jatinegara', '2025-08-30 10:33:01.844', '2025-08-30 10:33:01.000'),
	('jakarta-utara', 'Jakarta Utara', 'Penjaringan', 'Penjaringan', '2025-08-30 10:33:01.844', '2025-08-30 10:33:01.000');

-- Dumping structure for table iqgncnzy_skripsi.assets
CREATE TABLE IF NOT EXISTS `assets` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL,
  `brand` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `date` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `mac_address` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `model` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `price` double NOT NULL,
  `quantity` double NOT NULL,
  `serial_number` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `site` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `status_in_out` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.assets: ~4 rows (approximately)
INSERT INTO `assets` (`id`, `createdAt`, `updatedAt`, `brand`, `type`, `date`, `description`, `mac_address`, `model`, `price`, `quantity`, `serial_number`, `site`, `status`, `status_in_out`) VALUES
	('5ca1606b-66b3-4958-b7af-f48d4cda800a', '2025-05-12 20:09:33.822', '2025-06-02 14:47:46.478', 'TP-Link', 'Router Edit', '2025-06-02T00:00:00.000Z', '', '12:12:12:12:12:12', 'RBG128', 10000, 10, '123456789', 'asd', 'broken', 'in'),
	('842a48a4-6380-4280-96e1-a734f61a7d5b', '2025-06-03 05:48:35.019', '2025-06-02 22:48:35.649', 'ads', 'asdzxc', '2025-06-02T22:47:59.641Z', '123123', 'asd', 'asd', 33, 1, 'ads', 'asd', 'good', 'out'),
	('ac9e147a-4a75-4d39-9068-83e28ea0288b', '2025-05-10 17:17:03.502', '2025-05-10 17:18:18.355', 'TP-Link', 'Router Update2', '2025-05-05', '', '12:12:12:12:12:12', 'RBG128', 10000, 10, '123456789', '', '', 'in'),
	('f2061760-b41b-420b-9f70-4e96e46f2f57', '2025-06-03 05:48:26.647', '2025-06-02 22:48:27.277', 'ads', 'asd', '2025-06-02T22:47:59.641Z', '123123', 'asd', 'asd', 33, 1, 'ads', 'asd', 'good', 'out');

-- Dumping structure for table iqgncnzy_skripsi.cable
CREATE TABLE IF NOT EXISTS `cable` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Cable type (UTP, Fiber, etc.)',
  `length` decimal(10,2) DEFAULT NULL COMMENT 'Total cable length in meters',
  `status` enum('available','in_use','damaged','retired') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'available',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.cable: ~3 rows (approximately)
INSERT INTO `cable` (`id`, `name`, `type`, `length`, `status`, `created_at`, `updated_at`) VALUES
	('cable-1', 'UTP Cable 100m', 'UTP Cat6', 100.00, 'available', '2025-08-30 03:49:53', '2025-08-30 03:49:53'),
	('cable-2', 'Fiber Cable 500m', 'Single Mode Fiber', 500.00, 'available', '2025-08-30 03:49:53', '2025-08-30 03:49:53'),
	('cable-3', 'UTP Cable 50m', 'UTP Cat5e', 50.00, 'available', '2025-08-30 03:49:53', '2025-08-30 03:49:53');

-- Dumping structure for table iqgncnzy_skripsi.company
CREATE TABLE IF NOT EXISTS `company` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `logo_url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `npwp` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `address` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.company: ~4 rows (approximately)
INSERT INTO `company` (`id`, `name`, `url`, `email`, `phone`, `logo_url`, `npwp`, `address`, `createdAt`, `updatedAt`, `description`) VALUES
	('23e49db5-239c-4434-8848-a5f6de7d6d06', 'asdasd', 'https://www.google.com/url?sa=i&url=https%3A%2F%2Fid.pinterest.com%2Fcahyagaluh92%2Fgambar-simpel%2F&psig=AOvVaw1ylchxhjvGkY6QjeEHXTX-&ust=1748991019102000&source=images&cd=vfe&opi=89978449&v', 'asdasdasdsa@email.com', '0823125426123', 'https://www.google.com/url?sa=i&url=https%3A%2F%2Fid.pinterest.com%2Fcahyagaluh92%2Fgambar-simpel%2F&psig=AOvVaw1ylchxhjvGkY6QjeEHXTX-&ust=1748991019102000&source=images&cd=vfe&opi=89978449&v', '6543234t12', '213123123', '2025-06-03 06:04:31.289', '2025-06-03 06:04:31.738', ''),
	('53140661-9c49-423d-a00e-c8fe21680213', 'Lilly Networks', 'https://www.lilly.net.id', 'info@lilly.net.id', '03418222099', 'https://rtrw.net.id/img/freeze/logo.png', '1112', '1112', '2025-05-12 13:03:51.179', '2025-06-02 14:55:58.219', ''),
	('786a516f-11de-4eac-88cf-55f42e7e1227', 'CV Sukses Mandiri', 'https://suksesmandiri.co.id', 'contact@suksesmandiri.co.id', '021-5550456', '', '', 'Jl. Thamrin No. 45, Jakarta Pusat', '2025-08-22 10:42:23.737', '2025-08-22 10:42:23.737', 'Innovative business solutions'),
	('d5a7d66e-c313-4520-a482-03d58ba2764b', 'PT Maju Bersama', 'https://majubersama.com', 'info@majubersama.com', '021-5550123', '', '', 'Jl. Sudirman No. 123, Jakarta Pusat', '2025-08-22 10:42:23.732', '2025-08-22 10:42:23.732', 'Leading technology solutions provider'),
	('ead2060c-37d7-41e7-9e73-c9a419340e3e', 'Company 1', 'http://127.0.0.1', 'company1@email.com', '9764020363', '', '8782728399387222', 'Jalan Pisang Kipas, Dinoyo, Malang, Kota Malang, East Java, Jawa, 65113, Indonesia', '2025-06-17 15:34:16.948', '2025-06-17 08:34:17.157', ''),
	('ee6ccd53-a43e-40d2-a017-5f458ec754da', 'asdasd', 'https://www.google.com/url?sa=i&url=https%3A%2F%2Fid.pinterest.com%2Fcahyagaluh92%2Fgambar-simpel%2F&psig=AOvVaw1ylchxhjvGkY6QjeEHXTX-&ust=1748991019102000&source=images&cd=vfe&opi=89978449&v', 'asdasdasdsa@email.com', '0823125426123', 'https://www.google.com/url?sa=i&url=https%3A%2F%2Fid.pinterest.com%2Fcahyagaluh92%2Fgambar-simpel%2F&psig=AOvVaw1ylchxhjvGkY6QjeEHXTX-&ust=1748991019102000&source=images&cd=vfe&opi=89978449&v', '6543234t12', '213123123', '2025-06-03 06:04:31.440', '2025-06-03 06:04:31.889', ''),
	('fc9f6efe-9914-49a1-b175-35ee425a8d66', 'UD Makmur Jaya', 'https://makmurjaya.com', 'admin@makmurjaya.com', '021-5550789', '', '', 'Jl. Gatot Subroto No. 67, Jakarta Selatan', '2025-08-22 10:42:23.742', '2025-08-22 10:42:23.742', 'Reliable service provider'),
	('lilly-isp', 'Lilly ISP', 'https://lillyisp.com', 'info@lillyisp.com', '+62-123-456-789', NULL, '123456789012345', 'Jl. Example No. 123, Jakarta', '2025-08-30 10:33:01.841', '2025-08-30 10:33:01.000', NULL);

-- Dumping structure for table iqgncnzy_skripsi.customer
CREATE TABLE IF NOT EXISTS `customer` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `address` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `area_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `card_identition` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `company_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `gender` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `status_user` enum('Active','Nonactive') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'Active',
  `product_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `job` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `latitude` double NOT NULL,
  `longitude` double NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `alias` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `no_identition` int NOT NULL,
  `password` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type_of_service` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL,
  `installation_date` date NOT NULL,
  `next_payment_date` date NOT NULL,
  PRIMARY KEY (`id`),
  KEY `company_id` (`company_id`),
  KEY `customer_ibfk_1` (`area_id`),
  KEY `customer_ibfk_3` (`product_id`),
  KEY `idx_customer_company_id` (`company_id`),
  KEY `idx_customer_area_id` (`area_id`),
  KEY `idx_customer_product_id` (`product_id`),
  CONSTRAINT `customer_ibfk_1` FOREIGN KEY (`area_id`) REFERENCES `areas` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `customer_ibfk_2` FOREIGN KEY (`company_id`) REFERENCES `company` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `customer_ibfk_3` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.customer: ~7 rows (approximately)
INSERT INTO `customer` (`id`, `address`, `area_id`, `card_identition`, `company_id`, `email`, `gender`, `status_user`, `product_id`, `job`, `latitude`, `longitude`, `name`, `alias`, `no_identition`, `password`, `phone`, `type_of_service`, `createdAt`, `updatedAt`, `installation_date`, `next_payment_date`) VALUES
	('24225552-b0c3-43d7-9a67-20e04f36fa5f', 'Jatirenggo, Turen, Talok, Kabupaten Malang, Jawa Timur, Jawa, Indonesia', '61362c06-7042-4ed4-a391-237b742912c3', 'paspor', '53140661-9c49-423d-a00e-c8fe21680213', 'nunu@email.com', 'male', 'Active', '196df683-76ae-4a60-8357-87ce6c7a74a5', '', -8.188722878690669, 112.7005007657605, 'Cha Eunwoo', '', 2147483647, '', '081217706557', 'internet', '2025-06-29 04:33:02.614', '2025-06-07 15:28:58.780', '2025-06-02', '2025-07-02'),
	('55f03555-973d-425b-92bd-a0825c20b8ce', 'Mataraman, Turen, Talok, Kabupaten Malang, East Java, Java, Indonesia', '61362c06-7042-4ed4-a391-237b742912c3', 'ktp', '53140661-9c49-423d-a00e-c8fe21680213', 'rita.darmayanti@gmail.com', 'female', 'Active', '5c96b70e-4d7f-4cfa-8a57-779c0e00e3f4', '', -8.1947818, 112.6998068, 'Rita Darmayanti', '', 2147483647, '$2a$12$JfHseQPJBKPNSFNZ4HRtbepWGlHRzpFH/eYd2R7gPgZ0sOeHi39kS', '081333207500', 'internet', '2025-06-28 06:47:46.653', '2025-06-28 14:24:31.595', '2025-06-28', '2025-07-28'),
	('7770c03c-0e8c-4a1e-8b04-5f9183b4cc82', 'Jl. Gatot Subroto No. 67, Jakarta Selatan', '2f4e724d-39ab-4c8e-a5f9-437e100f7b12', '', 'fc9f6efe-9914-49a1-b175-35ee425a8d66', 'admin@makmurjaya.com', '', 'Active', '7052ad85-bc6b-4dfa-852c-9fa1b2e50f58', '', -6.2088, 106.8233, 'UD Makmur Jaya', '', 0, '', '021-5550789', '', '2025-08-22 10:45:04.126', '2025-08-22 10:45:04.126', '2025-08-22', '2025-09-22'),
	('94cb331d-e180-4259-a99a-7f24c24bcc64', 'Panunggangan Utara, Pinang, Tangerang, Banten, Jawa, 15143, Indonesia', '61362c06-7042-4ed4-a391-237b742912c3', 'ktp', '53140661-9c49-423d-a00e-c8fe21680213', 'customer@email.com', 'male', 'Active', 'ad24e250-fa6a-4685-b0cd-bda226ec11fd', '', -6.219311792067713, 106.64605227736672, 'customer', '', 2147483647, '$2a$12$JfHseQPJBKPNSFNZ4HRtbepWGlHRzpFH/eYd2R7gPgZ0sOeHi39kS', '082231900226', 'internet', '2025-07-04 16:18:46.079', '2025-07-04 16:18:46.079', '2025-07-04', '2025-08-04'),
	('9b8c7e59-b6ef-4097-8ad9-6b980ce89505', 'Jalan Pisang Kipas, Dinoyo, Malang, Kota Malang, East Java, Jawa, 65113, Indonesia', 'b53e9371-036f-4dcb-88b8-0a3679789e76', 'ktp', '53140661-9c49-423d-a00e-c8fe21680213', 'mufidah@email.com', 'female', 'Active', '196df683-76ae-4a60-8357-87ce6c7a74a5', 'job', -7.9429632, 112.6170624, 'Mufidah', '', 2147483647, '$2a$12$JfHseQPJBKPNSFNZ4HRtbepWGlHRzpFH/eYd2R7gPgZ0sOeHi39kS', '082234469663', 'internet', '2025-05-13 04:23:37.510', '2025-06-05 19:32:00.429', '2025-05-13', '2025-06-13'),
	('bace4349-5f70-419d-93dc-fc54efbcb839', 'Jl. Thamrin No. 45, Jakarta Pusat', 'ee888c0f-c3d7-4cc1-aacb-c153d52d3e1b', '', '786a516f-11de-4eac-88cf-55f42e7e1227', 'info@suksesmandiri.co.id', '', 'Active', '6e5c9a56-707e-4486-b86c-9bfa3cfeba4d', '', -6.1865, 106.8222, 'CV Sukses Mandiri', '', 0, '', '021-5550456', '', '2025-08-22 10:45:04.122', '2025-08-22 10:45:04.122', '2025-08-22', '2025-09-22'),
	('f0fdbb22-7c7b-4584-8e4e-0f775baefc3c', 'Jl. Sudirman No. 123, Jakarta Pusat', 'ee888c0f-c3d7-4cc1-aacb-c153d52d3e1b', '', 'd5a7d66e-c313-4520-a482-03d58ba2764b', 'contact@majubersama.com', '', 'Active', '378c4d01-6248-4c74-a4c0-f9815526f4fe', '', -6.2088, 106.8456, 'PT Maju Bersama', '', 0, '', '021-5550123', '', '2025-08-22 10:45:04.120', '2025-08-22 10:45:04.120', '2025-08-22', '2025-09-22');

-- Dumping structure for table iqgncnzy_skripsi.customer_details
CREATE TABLE IF NOT EXISTS `customer_details` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `technician_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT (now(3)),
  `updatedAt` datetime(3) NOT NULL,
  KEY `customer_installations_ibfk_2` (`technician_id`) USING BTREE,
  KEY `idx_customer_installations_technician_id` (`technician_id`) USING BTREE,
  KEY `Index 3` (`id`) USING BTREE,
  CONSTRAINT `customer_details_ibfk_2` FOREIGN KEY (`technician_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_customer_details_customer_installations` FOREIGN KEY (`id`) REFERENCES `customer_installations` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC;

-- Dumping data for table iqgncnzy_skripsi.customer_details: ~0 rows (approximately)

-- Dumping structure for table iqgncnzy_skripsi.customer_installations
CREATE TABLE IF NOT EXISTS `customer_installations` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `customer_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `details_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `date_trial` date NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT (now()),
  `updatedAt` datetime(3) NOT NULL DEFAULT (now()),
  PRIMARY KEY (`id`),
  KEY `archive_installation_ibfk_1` (`customer_id`),
  KEY `idx_customer_installations_customer_id` (`customer_id`),
  KEY `customer_installations_ibfk_2` (`details_id`) USING BTREE,
  KEY `idx_customer_installations_technician_id` (`details_id`) USING BTREE,
  CONSTRAINT `customer_installations_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.customer_installations: ~0 rows (approximately)

-- Dumping structure for table iqgncnzy_skripsi.customer_services
CREATE TABLE IF NOT EXISTS `customer_services` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `customer_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `device_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cable_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cable_length` decimal(10,2) DEFAULT NULL COMMENT 'Cable length in meters',
  `end_port_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Type of end port (RJ45, Fiber, etc.)',
  `user_login` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `service_activation_date` date DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `customer_id` (`customer_id`),
  KEY `device_id` (`device_id`),
  CONSTRAINT `customer_services_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`) ON DELETE CASCADE,
  CONSTRAINT `customer_services_ibfk_2` FOREIGN KEY (`device_id`) REFERENCES `network_devices` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.customer_services: ~0 rows (approximately)

-- Dumping structure for table iqgncnzy_skripsi.device
CREATE TABLE IF NOT EXISTS `device` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.device: ~0 rows (approximately)

-- Dumping structure for table iqgncnzy_skripsi.history_ip
CREATE TABLE IF NOT EXISTS `history_ip` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `customer_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `old_ip` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `changed_at` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`id`) USING BTREE,
  KEY `Index 2` (`customer_id`) USING BTREE,
  CONSTRAINT `FK_history_ip_customer` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='ip_static';

-- Dumping data for table iqgncnzy_skripsi.history_ip: ~0 rows (approximately)

-- Dumping structure for table iqgncnzy_skripsi.history_mac
CREATE TABLE IF NOT EXISTS `history_mac` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `customer_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `old_ip` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `changed_at` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`id`) USING BTREE,
  KEY `Index 2` (`customer_id`) USING BTREE,
  CONSTRAINT `FK_history_mac_customer` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.history_mac: ~0 rows (approximately)

-- Dumping structure for table iqgncnzy_skripsi.images
CREATE TABLE IF NOT EXISTS `images` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `file` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `full_path` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `archive_installation_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.images: ~11 rows (approximately)
INSERT INTO `images` (`id`, `file`, `full_path`, `archive_installation_id`, `createdAt`, `updatedAt`) VALUES
	('0f21b7c4-dfea-4683-bd8f-8720ff826e1d', 'imageInstallation1.png', 'uploads/c13a6c87-ec28-47ba-84c2-58b5ace2af57/94cb331d-e180-4259-a99a-7f24c24bcc64/imageInstallation1.png', '404ba44d-98e1-4645-966a-0e1e7f46dcdb', '2025-07-04 16:38:23.631', '2025-07-04 16:38:32.118'),
	('40b6b53d-4988-4044-be50-d40a8fcd3c0c', 'imageInstallation1.jpg', 'uploads/c13a6c87-ec28-47ba-84c2-58b5ace2af57/55f03555-973d-425b-92bd-a0825c20b8ce/imageInstallation1.jpg', '', '2025-06-30 23:07:32.961', NULL),
	('4a182a21-39b6-4d7b-a131-6a34016d7162', 'imageInstallation1.jpg', 'uploads/c13a6c87-ec28-47ba-84c2-58b5ace2af57/55f03555-973d-425b-92bd-a0825c20b8ce/imageInstallation1.jpg', 'f1f093a6-05c5-47e9-8f93-b8556dc26fa0', '2025-06-30 23:08:06.653', '2025-06-30 16:08:17.217'),
	('5a77081a-39ca-4fa0-b559-ec8e55f14c46', 'ticket_cs_1756280257033.jpg', 'uploads/tickets/cs/ticket_cs_1756280257033.jpg', '', '2025-08-27 14:37:37.116', NULL),
	('70e805f9-9ada-4e6d-8b4d-dca10c73bdf1', 'imageInstallation1.jpg', 'uploads/c13a6c87-ec28-47ba-84c2-58b5ace2af57/55f03555-973d-425b-92bd-a0825c20b8ce/imageInstallation1.jpg', '', '2025-06-30 23:06:58.192', NULL),
	('b3aea40c-be6b-4bad-bc4f-e091cb259fbf', 'ticket_cs_1756262705972.png', 'uploads/tickets/cs/ticket_cs_1756262705972.png', '', '2025-08-27 09:45:06.058', NULL),
	('c66f7ab5-47a1-4b88-8927-ad6a1d48e1f4', 'imageInstallation1.png', 'uploads/ba73859d-8957-418a-8311-28a6b31cad95/24225552-b0c3-43d7-9a67-20e04f36fa5f/imageInstallation1.png', 'dbea45d6-f094-49d6-b509-7375a23ff120', '2025-06-27 22:48:57.194', '2025-06-27 22:49:00.236'),
	('cc76d0af-9484-4711-a2fc-e31d5b885694', 'ticket_cs_1756263490179.jpg', 'uploads/tickets/cs/ticket_cs_1756263490179.jpg', '', '2025-08-27 09:58:10.270', NULL),
	('d4b49880-9cb7-4492-9bb7-edccc13451fa', 'ticket_cs_1756262757432.jpg', 'uploads/tickets/cs/ticket_cs_1756262757432.jpg', '', '2025-08-27 09:45:57.492', NULL),
	('dd63dffd-6d9f-4300-b319-f1780287732c', 'ticket_cs_1756262755214.webp', 'uploads/tickets/cs/ticket_cs_1756262755214.webp', '', '2025-08-27 09:45:55.277', NULL),
	('e7bae98b-d2ff-41cf-9c05-8fcf890d850e', 'imageInstallation1.png', 'uploads/c13a6c87-ec28-47ba-84c2-58b5ace2af57/9b4f7d06-fe35-47fd-8c69-6c7c2155894c/imageInstallation1.png', '2be54c4d-e4f6-4270-9d98-e311fa3d15fe', '2025-06-17 16:17:43.670', '2025-06-17 09:17:57.142');

-- Dumping structure for table iqgncnzy_skripsi.invoices
CREATE TABLE IF NOT EXISTS `invoices` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `amount` int NOT NULL,
  `link` text CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `status` enum('paid','unpaid','pending') CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `customer_id` varchar(191) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.invoices: ~5 rows (approximately)
INSERT INTO `invoices` (`id`, `amount`, `link`, `status`, `customer_id`, `createdAt`, `updatedAt`) VALUES
	('0b948ff4-5a8b-4576-8ac9-05430d5fb7a5', 100000, 'https://app.sandbox.midtrans.com/snap/v4/redirection/2623bfa5-30e3-448a-b4dc-c1a0acd85a7f', 'paid', '9b8c7e59-b6ef-4097-8ad9-6b980ce89505', '2025-08-13 22:04:21.262', '2025-08-14 15:25:45.111'),
	('5400f096-ba28-4582-9442-fe3b5f54858a', 250000, 'https://app.sandbox.midtrans.com/snap/v4/redirection/b8f5726d-c1d7-4aa3-8f1d-9dba4c552f86', 'paid', '55f03555-973d-425b-92bd-a0825c20b8ce', '2025-06-28 14:03:57.118', '2025-06-28 07:29:15.379'),
	('9841e8ab-5c41-494e-8c9e-8708d3301f2b', 5000000, '', 'paid', '24225552-b0c3-43d7-9a67-20e04f36fa5f', '2025-08-28 14:58:04.541', '2025-08-28 14:58:28.848'),
	('dd3f4811-ffc0-4c14-92dd-1cff17d145da', 100000, 'https://app.sandbox.midtrans.com/snap/v4/redirection/fe667b0c-246d-46b8-ac5f-f18f2b38652b', 'paid', '9b8c7e59-b6ef-4097-8ad9-6b980ce89505', '2025-06-28 10:02:25.991', '2025-06-28 04:56:31.002'),
	('e1f1e6b6-04a9-4bba-b243-da60f9c7ec7e', 150000, 'https://app.sandbox.midtrans.com/snap/v4/redirection/c4b77f1a-8019-471e-a5c9-0791b9e41489', 'unpaid', '94cb331d-e180-4259-a99a-7f24c24bcc64', '2025-07-04 16:34:24.505', '2025-08-29 09:59:15.300');

-- Dumping structure for table iqgncnzy_skripsi.invoice_items
CREATE TABLE IF NOT EXISTS `invoice_items` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `qty` int NOT NULL,
  `total` int NOT NULL,
  `price` int NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL,
  `invoices_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `invoice_items_ibfk_1` (`invoices_id`),
  KEY `idx_invoice_items_invoices_id` (`invoices_id`),
  CONSTRAINT `invoice_items_ibfk_1` FOREIGN KEY (`invoices_id`) REFERENCES `invoices` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.invoice_items: ~5 rows (approximately)
INSERT INTO `invoice_items` (`id`, `name`, `qty`, `total`, `price`, `createdAt`, `updatedAt`, `invoices_id`) VALUES
	('3874c673-2ab6-4c15-8ff8-7d5487fe6802', '100 Mbps', 1, 100000, 100000, '2025-06-28 10:02:26.027', '2025-06-28 03:02:26.357', 'dd3f4811-ffc0-4c14-92dd-1cff17d145da'),
	('3b14ec84-a1c1-466f-a059-8c5dc2dc8cbb', '100 Mbps', 0, 0, 100000, '2025-08-14 22:04:25.055', '2025-08-14 22:04:25.416', '0b948ff4-5a8b-4576-8ac9-05430d5fb7a5'),
	('5bf86f17-b0d6-4e97-940b-c157ffd5cc43', '35 Mbps', 1, 250000, 250000, '2025-06-28 14:03:57.155', '2025-06-28 07:03:57.440', '5400f096-ba28-4582-9442-fe3b5f54858a'),
	('81555fea-50e4-4419-9624-4904e9973ade', 'Internet 100 Mbps', 10, 5000000, 500000, '2025-08-28 14:58:04.555', '2025-08-28 14:58:04.552', '9841e8ab-5c41-494e-8c9e-8708d3301f2b'),
	('a5969e20-99cb-46ba-bece-16ef0833fa22', '20 Mbps', 1, 150000, 150000, '2025-07-04 16:34:24.570', '2025-07-04 16:34:25.360', 'e1f1e6b6-04a9-4bba-b243-da60f9c7ec7e');

-- Dumping structure for table iqgncnzy_skripsi.logs
CREATE TABLE IF NOT EXISTS `logs` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `action` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `idx_logs_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.logs: ~0 rows (approximately)

-- Dumping structure for table iqgncnzy_skripsi.network_devices
CREATE TABLE IF NOT EXISTS `network_devices` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `customer_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `assets_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `switch_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `port_number` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `remote_port` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `eth_port` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `kepemilikan_perangkat` enum('owned','leased','customer') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'owned' COMMENT 'Device ownership status',
  `status_perangkat` enum('active','inactive','maintenance','faulty') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT 'Device status',
  `last_ping_status` enum('up','down','unknown') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'unknown',
  `last_ping_timestamp` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `mac_address` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ip_static` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `network_devices_ibfk_1` (`customer_id`),
  KEY `network_devices_ibfk_2` (`assets_id`),
  CONSTRAINT `network_devices_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `network_devices_ibfk_2` FOREIGN KEY (`assets_id`) REFERENCES `assets` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.network_devices: ~1 rows (approximately)
INSERT INTO `network_devices` (`id`, `customer_id`, `assets_id`, `switch_id`, `port_number`, `remote_port`, `eth_port`, `kepemilikan_perangkat`, `status_perangkat`, `last_ping_status`, `last_ping_timestamp`, `created_at`, `updated_at`, `mac_address`, `ip_static`) VALUES
	('sadvwabafbsdgd', '24225552-b0c3-43d7-9a67-20e04f36fa5f', '842a48a4-6380-4280-96e1-a734f61a7d5b', NULL, NULL, NULL, NULL, 'owned', 'active', 'up', '2025-08-30 06:00:03', '2025-08-30 05:54:41', '2025-08-30 06:01:39', '14:E2:O2:19:22:90', '172.168.92.12');

-- Dumping structure for table iqgncnzy_skripsi.products
CREATE TABLE IF NOT EXISTS `products` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `price` bigint NOT NULL,
  `description` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.products: ~12 rows (approximately)
INSERT INTO `products` (`id`, `name`, `price`, `description`, `createdAt`, `updatedAt`) VALUES
	('196df683-76ae-4a60-8357-87ce6c7a74a5', '100 Mbps', 100000, '', '2025-05-10 12:26:51.331', '2025-05-10 12:26:51.574'),
	('378c4d01-6248-4c74-a4c0-f9815526f4fe', 'Internet 100 Mbps', 500000, 'High-speed internet connection with 100 Mbps download speed', '2025-08-22 10:45:04.107', '2025-08-22 10:45:04.107'),
	('5c96b70e-4d7f-4cfa-8a57-779c0e00e3f4', '35 Mbps', 250000, '', '2025-06-28 13:36:11.644', '2025-06-28 06:36:11.934'),
	('6e5c9a56-707e-4486-b86c-9bfa3cfeba4d', 'Internet 50 Mbps', 300000, 'Standard internet connection with 50 Mbps download speed', '2025-08-22 10:45:04.111', '2025-08-22 10:45:04.110'),
	('7052ad85-bc6b-4dfa-852c-9fa1b2e50f58', 'Internet 25 Mbps', 200000, 'Basic internet connection with 25 Mbps download speed', '2025-08-22 10:45:04.113', '2025-08-22 10:45:04.113'),
	('ad24e250-fa6a-4685-b0cd-bda226ec11fd', '20 Mbps', 150000, 'Kecepatan 20 Mbps', '2025-06-17 15:33:04.585', '2025-06-17 08:33:04.793'),
	('aefb5b2e-c640-4a49-8038-145f4fd122d0', '1 Mbps', 1000, '', '2025-05-10 12:26:41.200', '2025-05-10 12:26:41.423'),
	('b88d16cf-5ed0-4a5b-92fa-c56b00ce7234', '10 Mbps', 10000, '', '2025-05-10 12:26:26.671', '2025-05-10 12:26:26.913'),
	('basic-10', 'Basic 10 Mbps', 200000, 'Basic internet package with 10 Mbps speed', '2025-08-30 10:33:01.847', '2025-08-30 10:33:01.000'),
	('business-100', 'Business 100 Mbps', 800000, 'Business internet package with 100 Mbps speed', '2025-08-30 10:33:01.847', '2025-08-30 10:33:01.000'),
	('premium-50', 'Premium 50 Mbps', 500000, 'Premium internet package with 50 Mbps speed', '2025-08-30 10:33:01.847', '2025-08-30 10:33:01.000'),
	('standard-25', 'Standard 25 Mbps', 350000, 'Standard internet package with 25 Mbps speed', '2025-08-30 10:33:01.847', '2025-08-30 10:33:01.000');

-- Dumping structure for table iqgncnzy_skripsi.roles
CREATE TABLE IF NOT EXISTS `roles` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `roles_name_key` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.roles: ~6 rows (approximately)
INSERT INTO `roles` (`id`, `name`, `createdAt`, `updatedAt`) VALUES
	('11f0fba5-b49c-4237-9250-ee1b873a7c2b', 'TECHNICIAN', '2025-05-10 14:51:17.000', NULL),
	('1631f7d5-7d01-40af-8d24-1692cefa205a', 'ADMIN', '2025-05-10 14:50:42.000', NULL),
	('752e699f-c48c-44ed-8dfa-c8962b4be7ab', 'NOC', '2025-08-19 16:34:49.639', '2025-08-19 16:34:49.637'),
	('8b439b4a-c34c-4c2b-9114-2bc43121c2c2', 'CUSTOMER_SERVICE', '2025-08-23 16:30:54.232', '2025-08-23 16:30:54.231'),
	('a6eea5dc-ddd1-409c-bc3b-2985d6787b4f', 'CUSTOMER SERVICE', '2025-08-21 09:20:08.057', '2025-08-21 09:20:08.056'),
	('d5e4681b-9fe1-4cc1-80bb-a2b569458207', 'FINANCE', '2025-05-10 14:51:36.000', NULL);

-- Dumping structure for table iqgncnzy_skripsi.switch
CREATE TABLE IF NOT EXISTS `switch` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `model` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ip_address` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` enum('active','inactive','maintenance') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'active',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.switch: ~3 rows (approximately)
INSERT INTO `switch` (`id`, `name`, `model`, `ip_address`, `location`, `status`, `created_at`, `updated_at`) VALUES
	('switch-1', 'Core Switch 1', 'Cisco Catalyst 3850', '192.168.1.1', 'Data Center', 'active', '2025-08-30 03:49:53', '2025-08-30 03:49:53'),
	('switch-2', 'Access Switch 1', 'Cisco Catalyst 2960', '192.168.1.2', 'Building A', 'active', '2025-08-30 03:49:53', '2025-08-30 03:49:53'),
	('switch-3', 'Access Switch 2', 'Cisco Catalyst 2960', '192.168.1.3', 'Building B', 'active', '2025-08-30 03:49:53', '2025-08-30 03:49:53');

-- Dumping structure for table iqgncnzy_skripsi.transactions
CREATE TABLE IF NOT EXISTS `transactions` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `account_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type_cash` enum('internet','cash_flow') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type_in_out` enum('debit','credit') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `date` datetime(3) NOT NULL,
  `description` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `amount` bigint NOT NULL,
  `category` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `method` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `invoice_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `transactions_account_id_fkey` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.transactions: ~0 rows (approximately)
INSERT INTO `transactions` (`id`, `account_id`, `type_cash`, `type_in_out`, `date`, `description`, `amount`, `category`, `method`, `invoice_id`, `createdAt`, `updatedAt`) VALUES
	('02233804-f41e-4485-8608-f1aba7ad4068', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-07-04 10:10:00.000', '', 150000, '', '', '', '2025-07-04 17:09:59.313', '2025-07-04 10:10:00.172'),
	('176b3a76-ee72-4912-88c7-4afad026bee9', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'cash_flow', 'debit', '2025-05-17 06:12:18.000', 'Uang Masuk Perusahaan', 10000, '', '', NULL, '2025-05-17 06:12:18.061', '2025-05-17 06:12:18.784'),
	('4bb50986-1652-4922-a97a-c49b18171196', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-05-17 06:28:28.000', '', 5000, '', '', NULL, '2025-05-17 06:28:27.792', '2025-05-17 06:28:28.521'),
	('58912ebb-f3d9-4b51-9aca-a8dddcf4d116', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-05-18 05:28:41.000', '', 50000, '', '', NULL, '2025-05-18 12:28:41.009', '2025-05-18 05:28:41.225'),
	('702d7640-06e5-42c3-81f4-758a30f4681f', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-05-17 05:28:00.000', '', 50000, '', '', NULL, '2025-05-17 12:28:00.005', '2025-05-17 05:28:00.689'),
	('722d93b9-be79-415b-932f-94759aaf4f85', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-05-17 06:40:33.000', '', 5000, '', '', NULL, '2025-05-17 06:40:32.274', '2025-05-17 06:40:33.008'),
	('89955fa4-8b53-4821-8e46-4d1d7586f0eb', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-07-04 09:53:58.000', '', 150000, '', '', '', '2025-07-04 16:53:57.378', '2025-07-04 09:53:58.241'),
	('8fa1a010-2d7e-4ced-a8cf-6d55e1c39e53', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-08-14 15:25:46.000', '', 100000, '', '', '', '2025-08-14 22:25:46.055', '2025-08-14 15:25:46.441'),
	('9748c7fd-99f5-4c52-beb8-e897fea5e32a', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'cash_flow', 'debit', '2025-05-17 06:12:31.000', 'Uang Masuk Perusahaan', 5000, '', '', NULL, '2025-05-17 06:12:30.811', '2025-05-17 06:12:31.534'),
	('9e4f3743-800a-4a22-9b6a-baf6b19f9646', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'cash_flow', 'debit', '2025-05-17 06:12:20.000', 'Uang Masuk Perusahaan', 10000, '', '', NULL, '2025-05-17 06:12:19.805', '2025-05-17 06:12:20.517'),
	('abba65fe-bcf0-4b82-a792-f924f5b6eb7d', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'cash_flow', 'credit', '2025-06-06 13:36:12.000', 'Hutanng', 100000, 'pemberian pinjaman kasbon', 'transfer', '', '2025-06-06 13:36:11.465', '2025-06-06 13:36:12.056'),
	('ac431992-f022-49ce-82a9-dd7dd3ed56f9', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-05-17 06:30:08.000', '', 5000, '', '', NULL, '2025-05-17 06:30:07.440', '2025-05-17 06:30:08.167'),
	('ae039927-a7c2-4084-b496-3c0a5a225797', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'cash_flow', 'debit', '2025-06-06 13:18:46.000', 'Hutanng', 100000, 'pengembalian kasbon pegawai', 'transfer', '', '2025-06-06 13:18:45.594', '2025-06-06 13:18:46.228'),
	('cbacdb2c-b785-42be-bc0a-8d28dd06bc59', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-06-28 07:29:15.000', '', 250000, '', '', '', '2025-06-28 14:29:15.285', '2025-06-28 07:29:15.558'),
	('de10948b-2adc-4a32-b000-5d4dc1c2b8df', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'cash_flow', 'debit', '2025-06-06 13:16:53.000', 'Hutanng', 100000, '', '', '', '2025-06-06 13:16:53.342', '2025-06-06 13:16:53.985'),
	('e55171ed-ba5b-4010-a68e-d10268fa970d', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-06-28 07:13:14.000', '', 250000, '', '', '', '2025-06-28 14:13:13.825', '2025-06-28 07:13:14.110'),
	('e55dbc81-7e01-41a9-aa97-1af45f5bbef7', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-08-14 15:10:40.000', '', 100000, '', '', '', '2025-08-14 22:10:39.738', '2025-08-14 15:10:40.145'),
	('eac18e83-a958-4c62-8ad9-ce74dfa21f57', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-05-17 06:28:54.000', '', 5000, '', '', NULL, '2025-05-17 06:28:53.499', '2025-05-17 06:28:54.228'),
	('ebce5278-ec6b-47b8-8519-ed6164e03a65', 'b82074e7-7acb-40e2-ab33-014f4b09c1f8', 'internet', 'debit', '2025-08-14 15:09:53.000', '', 100000, '', '', '', '2025-08-14 22:09:54.548', '2025-08-14 15:09:53.199');

-- Dumping structure for table iqgncnzy_skripsi.trouble_tickets
CREATE TABLE IF NOT EXISTS `trouble_tickets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `customer_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `title` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `status` enum('finished','ongoing','unfinished') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `assigned_to` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `current_assignee_role` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '1631f7d5-7d01-40af-8d24-1692cefa205a',
  `customer_note` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `noc_note` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `technician_note` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `img_cs` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'Image uploaded by Customer Service',
  `img_noc` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'Image uploaded by NOC',
  `img_tech_bf` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'Image uploaded by Technician (Before)',
  `img_tech_af` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'Image uploaded by Technician (After)',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `Index 3` (`assigned_to`,`current_assignee_role`),
  KEY `Index 4` (`customer_id`),
  KEY `Index 2` (`type`) USING BTREE,
  KEY `idx_trouble_tickets_customer_id` (`customer_id`),
  KEY `idx_trouble_tickets_assigned_to` (`assigned_to`),
  KEY `idx_trouble_tickets_current_assignee_role` (`current_assignee_role`),
  KEY `idx_trouble_tickets_type` (`type`),
  KEY `idx_trouble_tickets_status` (`status`),
  CONSTRAINT `FK_trouble_tickets_customer` FOREIGN KEY (`customer_id`) REFERENCES `customer` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_trouble_tickets_roles` FOREIGN KEY (`current_assignee_role`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_trouble_tickets_trouble_type` FOREIGN KEY (`type`) REFERENCES `trouble_type` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_trouble_tickets_users` FOREIGN KEY (`assigned_to`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Dumping data for table iqgncnzy_skripsi.trouble_tickets: ~0 rows (approximately)
INSERT INTO `trouble_tickets` (`id`, `customer_id`, `type`, `title`, `description`, `status`, `assigned_to`, `current_assignee_role`, `customer_note`, `noc_note`, `technician_note`, `created_at`, `updated_at`, `img_cs`, `img_noc`, `img_tech_bf`, `img_tech_af`) VALUES
	(1, '24225552-b0c3-43d7-9a67-20e04f36fa5f', 'cdb1e92d-130a-42aa-bc5e-cd0c2bc34076', 'Gangguan setelah pemesangan 20 hari', 'Kabelnya mungkin ada gangguan', 'finished', 'b6cbf8b6-6a2e-45f3-a9ab-e30e5941bf5a', 'a6eea5dc-ddd1-409c-bc3b-2985d6787b4f', 'a', 'dw', NULL, '2025-08-22 12:34:58.000', '2025-08-25 19:28:32.960', '1-cs.jpg', '1-noc.png', NULL, NULL),
	(15, '24225552-b0c3-43d7-9a67-20e04f36fa5f', 'cdb1e92d-130a-42aa-bc5e-cd0c2bc34076', 'HVGGH', 'AWD', 'ongoing', '19226ccf-515e-44dd-b0bc-e62611f7a0af', '11f0fba5-b49c-4237-9250-ee1b873a7c2b', NULL, 'a', 'tes', '2025-08-23 10:45:24.688', '2025-08-25 17:44:28.114', NULL, NULL, '15-bf.png', '15-af.png'),
	(16, '24225552-b0c3-43d7-9a67-20e04f36fa5f', 'cdb1e92d-130a-42aa-bc5e-cd0c2bc34076', 'Internet down setelah 20 hari pemasangan di wajak', 'saya telah masang 20 hari tapi internet sering lemot saat jam 12 malam', 'unfinished', '19226ccf-515e-44dd-b0bc-e62611f7a0af', 'a6eea5dc-ddd1-409c-bc3b-2985d6787b4f', NULL, NULL, NULL, '2025-08-27 09:45:58.566', '2025-08-27 09:45:58.566', NULL, NULL, NULL, NULL),
	(17, '24225552-b0c3-43d7-9a67-20e04f36fa5f', 'cdb1e92d-130a-42aa-bc5e-cd0c2bc34076', 'Internet down setelah 20 hari pemasangan di wajak', 'saya telah masang 20 hari tapi internet sering lemot saat jam 12 malam', 'unfinished', '19226ccf-515e-44dd-b0bc-e62611f7a0af', 'a6eea5dc-ddd1-409c-bc3b-2985d6787b4f', NULL, NULL, NULL, '2025-08-27 09:46:02.213', '2025-08-27 09:46:02.213', NULL, NULL, NULL, NULL),
	(18, '24225552-b0c3-43d7-9a67-20e04f36fa5f', 'cdb1e92d-130a-42aa-bc5e-cd0c2bc34076', 'Internet down setelah 20 hari pemasangan di wajak', 'saya telah masang 20 hari tapi internet sering lemot saat jam 12 malam', 'unfinished', '19226ccf-515e-44dd-b0bc-e62611f7a0af', 'a6eea5dc-ddd1-409c-bc3b-2985d6787b4f', NULL, NULL, NULL, '2025-08-27 09:46:03.474', '2025-08-27 09:46:03.474', NULL, NULL, NULL, NULL),
	(19, '24225552-b0c3-43d7-9a67-20e04f36fa5f', 'cdb1e92d-130a-42aa-bc5e-cd0c2bc34076', 'Internet down setelah 20 hari pemasangan di wajak', 'saya telah masang 20 hari tapi internet sering lemot saat jam 12 malam', 'ongoing', NULL, '11f0fba5-b49c-4237-9250-ee1b873a7c2b', NULL, NULL, NULL, '2025-08-27 09:46:03.684', '2025-08-29 10:40:10.613', NULL, NULL, NULL, NULL),
	(20, '24225552-b0c3-43d7-9a67-20e04f36fa5f', 'cdb1e92d-130a-42aa-bc5e-cd0c2bc34076', 'Internet down setelah 20 hari pemasangan di wajak', 'saya telah masang 20 hari tapi internet sering lemot saat jam 12 malam', 'ongoing', NULL, '11f0fba5-b49c-4237-9250-ee1b873a7c2b', NULL, NULL, NULL, '2025-08-27 09:46:04.783', '2025-08-28 10:20:32.278', NULL, NULL, NULL, NULL),
	(21, '24225552-b0c3-43d7-9a67-20e04f36fa5f', '1965675f-fe09-4e5d-870d-88910ca37390', 'qwert', 'iuy', 'ongoing', NULL, '11f0fba5-b49c-4237-9250-ee1b873a7c2b', NULL, NULL, NULL, '2025-08-27 09:58:12.617', '2025-08-28 09:58:47.507', NULL, NULL, NULL, NULL),
	(22, '24225552-b0c3-43d7-9a67-20e04f36fa5f', '1965675f-fe09-4e5d-870d-88910ca37390', 'qwert', 'iuy', 'ongoing', '19226ccf-515e-44dd-b0bc-e62611f7a0af', '752e699f-c48c-44ed-8dfa-c8962b4be7ab', '', NULL, NULL, '2025-08-27 09:58:14.801', '2025-08-27 15:01:02.439', NULL, NULL, NULL, NULL),
	(23, '24225552-b0c3-43d7-9a67-20e04f36fa5f', '7e6bc903-c927-4aa6-a31c-42861e3f4f0e', 'qwert', 'iuy', 'finished', '19226ccf-515e-44dd-b0bc-e62611f7a0af', 'a6eea5dc-ddd1-409c-bc3b-2985d6787b4f', 'Kesalahan hanya ada di konfigurasi ', 'ini salah config tunggu jam 12 estimasi penyelesaian', NULL, '2025-08-27 09:58:16.998', '2025-08-27 10:30:21.385', '23-cs.jpg', NULL, NULL, NULL),
	(24, '55f03555-973d-425b-92bd-a0825c20b8ce', '77b909ff-bc77-412d-b9d5-fff4e4f25be9', 'Setelah pemasangan 5 hari ada lampu indikator mati di modem', 'kejadian jam 12 siang tanggal 17 agustus', 'ongoing', NULL, '11f0fba5-b49c-4237-9250-ee1b873a7c2b', NULL, NULL, NULL, '2025-08-27 14:37:39.483', '2025-08-28 09:58:31.056', NULL, NULL, NULL, NULL);

-- Dumping structure for table iqgncnzy_skripsi.trouble_type
CREATE TABLE IF NOT EXISTS `trouble_type` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.trouble_type: ~0 rows (approximately)
INSERT INTO `trouble_type` (`id`, `name`) VALUES
	('1965675f-fe09-4e5d-870d-88910ca37390', 'Modem'),
	('4508e57d-b623-46ec-a9b9-10ebfb7349a0', 'Power Outage'),
	('77b909ff-bc77-412d-b9d5-fff4e4f25be9', 'Modem Issue'),
	('7e6bc903-c927-4aa6-a31c-42861e3f4f0e', 'Configuration Error'),
	('86aa1703-6918-4b1c-8da2-6cd471897d43', 'Slow Connection'),
	('95b58430-2c55-494b-bfdb-ff3e7618b498', 'Internet Down'),
	('billing', 'Billing Problem'),
	('cdb1e92d-130a-42aa-bc5e-cd0c2bc34076', 'Cable Damage'),
	('hardware', 'Hardware Problem'),
	('installation', 'Installation Issue'),
	('maintenance', 'Maintenance Required'),
	('network', 'Network Issue'),
	('software', 'Software Issue');

-- Dumping structure for table iqgncnzy_skripsi.users
CREATE TABLE IF NOT EXISTS `users` (
  `id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAt` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updatedAt` datetime(3) DEFAULT NULL,
  `token` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `role_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `logo_url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `phone` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_email_key` (`email`),
  KEY `role_id` (`role_id`),
  KEY `idx_users_role_id` (`role_id`),
  CONSTRAINT `users_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table iqgncnzy_skripsi.users: ~17 rows (approximately)
INSERT INTO `users` (`id`, `email`, `name`, `password`, `createdAt`, `updatedAt`, `token`, `role_id`, `logo_url`, `phone`) VALUES
	('19226ccf-515e-44dd-b0bc-e62611f7a0af', 'cs@email.com', 'customerservice', '$2a$12$fH8zSNzzQCDzFJiX0XguQuP13SfOW/lBz2qYEMQG3dw6UXFP2sEhG', '2025-08-21 09:20:40.961', '2025-09-01 09:02:17.759', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJDVVNUT01FUiBTRVJWSUNFIiwiZXhwIjoxNzU2Njk1NzM3LCJpYXQiOjE3NTY2OTIxMzcsImlzcyI6ImxpbGx5YXBwcyIsInN1YiI6IjE5MjI2Y2NmLTUxNWUtNDRkZC1iMGJjLWU2MjYxMWY3YTBhZiJ9._0ZmS-4LBnSeEUj7efQnm8tTrz4cUUXYtxWjF5dVUtg', 'a6eea5dc-ddd1-409c-bc3b-2985d6787b4f', NULL, '0847377372'),
	('2f5f0cb8-9f65-46e9-9d1f-4499ab01201d', 'tech@lillyisp.com', 'Field Technician', '$2a$10$.640sjCovVn3buRu92ed0.DfDEXbfIGR2h87wfnu5g1FdBWQ33G2u', '2025-08-22 10:39:33.647', '2025-08-22 10:46:44.373', NULL, '11f0fba5-b49c-4237-9250-ee1b873a7c2b', NULL, '081234567893'),
	('32f6aca3-93dd-4178-ac90-eeb5f497807e', 'noc@email.com', 'NOC Engineer', '$2a$10$.640sjCovVn3buRu92ed0.DfDEXbfIGR2h87wfnu5g1FdBWQ33G2u', '2025-08-22 10:39:33.642', '2025-08-27 09:41:34.306', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJOT0MiLCJleHAiOjE3NTYyNjYwOTQsImlhdCI6MTc1NjI2MjQ5NCwiaXNzIjoibGlsbHlhcHBzIiwic3ViIjoiMzJmNmFjYTMtOTNkZC00MTc4LWFjOTAtZWViNWY0OTc4MDdlIn0.W5i2FRWYufEa6Yfe-Da4SSMF6bMhH-PRXo6rQPBIS3c', '752e699f-c48c-44ed-8dfa-c8962b4be7ab', NULL, '081234567892'),
	('4674ce2b-52d1-4bc5-8706-08c0b397f4ed', 'userq@email.com', 'user', '$2a$12$JfHseQPJBKPNSFNZ4HRtbepWGlHRzpFH/eYd2R7gPgZ0sOeHi39kS', '2025-05-10 21:01:32.904', '2025-06-03 23:40:44.371', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBRE1JTiIsImV4cCI6MTc0ODk3MjQ0NCwiaWF0IjoxNzQ4OTY4ODQ0LCJpc3MiOiJsaWxseWFwcHMiLCJzdWIiOiI0Njc0Y2UyYi01MmQxLTRiYzUtODcwNi0wOGMwYjM5N2Y0ZWQifQ.3AZ4yRfus5xWHYE7aaBoqa6P-BMnujxQRPExJ4v-aqg', '1631f7d5-7d01-40af-8d24-1692cefa205a', NULL, '123456789'),
	('5c82bc7b-4e96-4f9b-822e-00598d2ce389', 'finance@lillyisp.com', 'Finance Manager', '$2a$10$5q30dUEOR.67QsVBovCRmOnnL6He7qqLsIzDpNiInsvwTnZuRpM02', '2025-08-23 16:30:54.313', '2025-08-23 16:30:54.313', NULL, 'd5e4681b-9fe1-4cc1-80bb-a2b569458207', NULL, '081234567894'),
	('5f97d057-98dc-45c4-9094-dd4896baa6e1', 'angga@email.com', 'angga', '$2a$12$JfHseQPJBKPNSFNZ4HRtbepWGlHRzpFH/eYd2R7gPgZ0sOeHi39kS', '2025-05-10 12:24:36.000', '2025-06-06 12:38:29.903', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBRE1JTiIsImV4cCI6MTc0OTIxNzEwOSwiaWF0IjoxNzQ5MjEzNTA5LCJpc3MiOiJsaWxseWFwcHMiLCJzdWIiOiI1Zjk3ZDA1Ny05OGRjLTQ1YzQtOTA5NC1kZDQ4OTZiYWE2ZTEifQ.XrMirXoXxbhoba45TDChj_d9yxyKR6n2AgibNXLOaOI', '1631f7d5-7d01-40af-8d24-1692cefa205a', NULL, ''),
	('5f97d057-98dc-45c4-9094-dd4896baa6e3', 'mufidah@email.com', 'mufidah', '$2a$12$JfHseQPJBKPNSFNZ4HRtbepWGlHRzpFH/eYd2R7gPgZ0sOeHi39kS', '2025-05-10 12:24:36.000', '2025-06-30 16:01:03.913', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBRE1JTiIsImV4cCI6MTc1MTMwMjg2MywiaWF0IjoxNzUxMjk5MjYzLCJpc3MiOiJsaWxseWFwcHMiLCJzdWIiOiI1Zjk3ZDA1Ny05OGRjLTQ1YzQtOTA5NC1kZDQ4OTZiYWE2ZTMifQ.ARtsosjxeEPLSrJCqV9YpnV_3je2V96ZWL4QHuEel30', '1631f7d5-7d01-40af-8d24-1692cefa205a', NULL, ''),
	('5f97d057-98dc-45c4-9094-dd4896baa6e4', 'postman@email.com', 'postman', '$2a$12$JfHseQPJBKPNSFNZ4HRtbepWGlHRzpFH/eYd2R7gPgZ0sOeHi39kS', '2025-05-10 12:24:36.000', '2025-07-01 01:08:21.017', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBRE1JTiIsImV4cCI6MTc1MTMxMDUwMCwiaWF0IjoxNzUxMzA2OTAwLCJpc3MiOiJsaWxseWFwcHMiLCJzdWIiOiI1Zjk3ZDA1Ny05OGRjLTQ1YzQtOTA5NC1kZDQ4OTZiYWE2ZTQifQ.OwKog0yFCgeOIMr1cRRDnFXKyJn3u0RrygBcQQr9vpM', '1631f7d5-7d01-40af-8d24-1692cefa205a', NULL, ''),
	('5f97d057-98dc-45c4-9094-dd4896baa6ef', 'rama@email.com', 'rama', '$2a$12$JfHseQPJBKPNSFNZ4HRtbepWGlHRzpFH/eYd2R7gPgZ0sOeHi39kS', '2025-05-10 12:24:36.000', '2025-08-14 22:04:04.284', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBRE1JTiIsImV4cCI6MTc1NTE4NzQ0NCwiaWF0IjoxNzU1MTgzODQ0LCJpc3MiOiJsaWxseWFwcHMiLCJzdWIiOiI1Zjk3ZDA1Ny05OGRjLTQ1YzQtOTA5NC1kZDQ4OTZiYWE2ZWYifQ.Y9dgKIpa6WUJ0_2xjLvyaphq0PPb_8hv45Cz60gRAC0', '1631f7d5-7d01-40af-8d24-1692cefa205a', NULL, ''),
	('72aba9a4-616d-4a7e-b745-b3508052f767', 'user@email.com', 'user', '$2a$12$ZCRUSclA4HYRd2PJ/9tYq.NzvPnpOxlCY/4Nrs4rgmvRXZ42dLSGK', '2025-05-10 19:42:43.680', '2025-05-10 19:42:43.738', NULL, '1631f7d5-7d01-40af-8d24-1692cefa205a', NULL, '123456789'),
	('9ef48bbc-57bf-45cf-8f85-29b9b8950483', 'finance@email.com', 'finance', '$2a$12$W/keTbaM44fQgCbvf.jM/uHB.HhU.FRMnE2vvrzgY6tfPooO30t5K', '2025-05-17 11:19:37.773', '2025-06-28 04:39:02.985', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJGSU5BTkNFIiwiZXhwIjoxNzUxMDg5MTQyLCJpYXQiOjE3NTEwODU1NDIsImlzcyI6ImxpbGx5YXBwcyIsInN1YiI6IjllZjQ4YmJjLTU3YmYtNDVjZi04Zjg1LTI5YjliODk1MDQ4MyJ9.BdwJ56wQ9OCnbq9hIWcF43SnKLYCtozoz16Dhvorea8', 'd5e4681b-9fe1-4cc1-80bb-a2b569458207', NULL, '08672825172'),
	('b368ddcd-4e3e-4da6-8c6a-c5d6e5e56c36', 'noc@lillyisp.com', 'NOC Engineer', '$2a$10$5q30dUEOR.67QsVBovCRmOnnL6He7qqLsIzDpNiInsvwTnZuRpM02', '2025-08-23 16:30:54.308', '2025-08-23 16:30:54.308', NULL, '752e699f-c48c-44ed-8dfa-c8962b4be7ab', NULL, '081234567892'),
	('b6cbf8b6-6a2e-45f3-a9ab-e30e5941bf5a', 'admin@email.com', 'admin', '$2a$12$/Yx0a6X.rn7J.KLK3Yy4Oeimv3ASOqwn/UKnYXxCl52t4v/aq8sRa', '2025-06-28 09:10:33.933', '2025-08-23 12:40:20.289', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBRE1JTiIsImV4cCI6MTc1NTkzMTIyMCwiaWF0IjoxNzU1OTI3NjIwLCJpc3MiOiJsaWxseWFwcHMiLCJzdWIiOiJiNmNiZjhiNi02YTJlLTQ1ZjMtYTlhYi1lMzBlNTk0MWJmNWEifQ.xlnknmIg-vEN4pFDs65rfyrhEmtDsX4bayzEKy7ryB8', '1631f7d5-7d01-40af-8d24-1692cefa205a', NULL, '09864236772938'),
	('ba73859d-8957-418a-8311-28a6b31cad95', 'cantik@email.com', 'cantiks', '$2a$12$HNsi.5y1Kq/KbMeUdN8m/uNiWkjIUSDFs.ulxswdloHTAopP.5.9i', '2025-05-14 15:12:48.572', '2025-08-19 07:14:39.691', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJURUNITklDSUFOIiwiZXhwIjoxNzU1NTkxMjc5LCJpYXQiOjE3NTU1ODc2NzksImlzcyI6ImxpbGx5YXBwcyIsInN1YiI6ImJhNzM4NTlkLTg5NTctNDE4YS04MzExLTI4YTZiMzFjYWQ5NSJ9.JK8lduQXLQOPzfg1eaQLQWTy5YrRpZpIFrzRxShe-gk', '11f0fba5-b49c-4237-9250-ee1b873a7c2b', NULL, '086243292'),
	('c115a426-4db2-4bad-81f7-a50990b51406', 'cs@lillyisp.com', 'Customer Service Agent', '$2a$10$5q30dUEOR.67QsVBovCRmOnnL6He7qqLsIzDpNiInsvwTnZuRpM02', '2025-08-23 16:30:54.303', '2025-08-23 16:30:54.303', NULL, '8b439b4a-c34c-4c2b-9114-2bc43121c2c2', NULL, '081234567891'),
	('c13a6c87-ec28-47ba-84c2-58b5ace2af57', 'teknisi@email.com', 'teknisi', '$2a$12$HNsi.5y1Kq/KbMeUdN8m/uNiWkjIUSDFs.ulxswdloHTAopP.5.9i', '2025-06-03 15:37:29.178', '2025-08-25 17:42:16.703', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJURUNITklDSUFOIiwiZXhwIjoxNzU2MTIyMTM2LCJpYXQiOjE3NTYxMTg1MzYsImlzcyI6ImxpbGx5YXBwcyIsInN1YiI6ImMxM2E2Yzg3LWVjMjgtNDdiYS04NGMyLTU4YjVhY2UyYWY1NyJ9.13mWplKvt47DZFCgWRcHOxeB98NShkddWZ6TqwdZSLs', '11f0fba5-b49c-4237-9250-ee1b873a7c2b', NULL, '9764020363'),
	('c7d376ca-626c-4e33-bf47-342d2082c93d', 'admin@lillyisp.com', 'System Administrator', '$2a$10$5q30dUEOR.67QsVBovCRmOnnL6He7qqLsIzDpNiInsvwTnZuRpM02', '2025-08-23 16:30:54.299', '2025-09-01 15:30:40.514', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJBRE1JTiIsImV4cCI6MTc1NjcxOTA0MCwiaWF0IjoxNzU2NzE1NDQwLCJpc3MiOiJsaWxseWFwcHMiLCJzdWIiOiJjN2QzNzZjYS02MjZjLTRlMzMtYmY0Ny0zNDJkMjA4MmM5M2QifQ.sy8bM4aTDv9J0xKEBD0XSpGSMMrEcruNSc-vKV5NGHA', '1631f7d5-7d01-40af-8d24-1692cefa205a', NULL, '081234567890');

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
