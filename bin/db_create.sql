CREATE TABLE `AuthorizedClientUsers` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(45) NOT NULL,
  `client_credentials_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

CREATE TABLE `CheckedProvider` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `provider_id` int(11) NOT NULL,
  `fb_page_url` varchar(45) DEFAULT NULL,
  `fb_matched_phone` varchar(1) DEFAULT NULL,
  `fb_likes` int(11) DEFAULT NULL,
  `yelp_rating` decimal(5,2) DEFAULT NULL,
  `yelp_review_count` int(11) DEFAULT NULL,
  `fc_user_website` varchar(45) DEFAULT NULL,
  `fc_twitter_id` varchar(45) DEFAULT NULL,
  `fc_twitter_followers` int(11) DEFAULT NULL,
  `fc_linkedin_id` varchar(1) NOT NULL DEFAULT 'N',
  `in_master_list` varchar(1) NOT NULL DEFAULT 'N',
  `ln_biz_match` varchar(1) NOT NULL DEFAULT 'N',
  `ln_biz_matched_phone` varchar(1) NOT NULL DEFAULT 'N',
  `ln_biz_raw_json` varchar(2000) DEFAULT NULL,
  `ln_indiv_match` varchar(1) NOT NULL DEFAULT 'N',
  `ln_indiv_matched_phone` varchar(1) NOT NULL DEFAULT 'N',
  `ln_indiv_matched_zip` varchar(1) NOT NULL DEFAULT 'N',
  `ln_indiv_raw_json` varchar(2000) DEFAULT NULL,
  `google_match` varchar(1) NOT NULL DEFAULT 'N',
  `google_phone_match` varchar(1) NOT NULL DEFAULT 'N',
  `ln_bizinstid_match` varchar(1) NOT NULL DEFAULT 'N',
  `ln_bizinstid_matched_phone` varchar(1) NOT NULL DEFAULT 'N',
  `ln_bizinstid_matched_name` varchar(1) NOT NULL DEFAULT 'N',
  `ln_bizinstid_matched_zip` varchar(1) NOT NULL DEFAULT 'N',
  `fc_biz_match` varchar(1) NOT NULL DEFAULT 'N',
  `fc_biz_matched_org` varchar(1) NOT NULL DEFAULT 'N',
  `fc_biz_fb_page` varchar(45) DEFAULT NULL,
  `fc_biz_twitter_id` varchar(45) DEFAULT NULL,
  `fc_biz_twitter_followers` int(11) NOT NULL DEFAULT '0',
  `fc_biz_linkedin_match` varchar(1) NOT NULL DEFAULT 'N',
  `experian_match` varchar(1) NOT NULL DEFAULT 'N',
  `experian_biz_id` varchar(45) DEFAULT NULL,
  `bgc_status` varchar(45) NOT NULL DEFAULT 'I',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21081 DEFAULT CHARSET=latin1;

CREATE TABLE `Checker` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `CheckerLog` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `provider_id` int(11) NOT NULL,
  `check_id` int(11) NOT NULL,
  `check_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `ClientCredentials` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `client_platform_id` int(11) NOT NULL,
  `api_key` varchar(45) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

CREATE TABLE `ClientPlatform` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `type` varchar(45) NOT NULL,
  `do_id_verification` int(11) NOT NULL DEFAULT '0',
  `do_bkg_check` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1;

CREATE TABLE `Job` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` varchar(45) NOT NULL,
  `cost` int(11) NOT NULL,
  `description` varchar(200) DEFAULT NULL,
  `startdate` datetime DEFAULT NULL,
  `status` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=511 DEFAULT CHARSET=latin1;

CREATE TABLE `Note` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `date` datetime NOT NULL,
  `value` varchar(200) NOT NULL,
  `project_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=62 DEFAULT CHARSET=latin1;

CREATE TABLE `Project` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `job_id` int(11) DEFAULT NULL,
  `provider_id` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `status` varchar(45) NOT NULL,
  `create_date` datetime DEFAULT NULL,
  `encoded_id` varchar(45) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=527 DEFAULT CHARSET=latin1;

CREATE TABLE `Provider` (
  `id` int(32) NOT NULL AUTO_INCREMENT,
  `name` varchar(120) NOT NULL,
  `phone` bigint(64) NOT NULL,
  `type` varchar(45) NOT NULL,
  `email` varchar(45) DEFAULT NULL,
  `address1` varchar(45) NOT NULL,
  `address2` varchar(45) DEFAULT NULL,
  `city` varchar(45) NOT NULL,
  `state` varchar(45) NOT NULL,
  `zip` varchar(45) NOT NULL,
  `ownername` varchar(45) DEFAULT NULL,
  `website` varchar(45) DEFAULT NULL,
  `category` varchar(45) DEFAULT NULL,
  `source` varchar(15) NOT NULL,
  `is_business` varchar(1) NOT NULL DEFAULT 'N',
  `client_platform_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3460 DEFAULT CHARSET=latin1;

CREATE TABLE `ProviderMaster` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `company_name` varchar(45) NOT NULL,
  `city` varchar(45) NOT NULL,
  `state` varchar(2) NOT NULL,
  `zip` varchar(10) NOT NULL,
  `email` varchar(45) DEFAULT NULL,
  `phone` varchar(10) DEFAULT NULL,
  `average_grade` decimal(6,2) DEFAULT NULL,
  `jobs_count` int(11) DEFAULT NULL,
  `source` varchar(15) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=193401 DEFAULT CHARSET=utf8;

CREATE TABLE `User` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(45) NOT NULL,
  `last_name` varchar(45) NOT NULL,
  `address1` varchar(45) NOT NULL,
  `address2` varchar(45) DEFAULT NULL,
  `city` varchar(45) NOT NULL,
  `state` varchar(45) NOT NULL,
  `zip` varchar(45) NOT NULL,
  `email` varchar(45) DEFAULT NULL,
  `phone` bigint(64) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=128 DEFAULT CHARSET=latin1;

CREATE TABLE `Payment` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `project_id` int(11) NOT NULL,
  `transaction_id` varchar(45) NOT NULL,
  `email` varchar(45) NOT NULL,
  `fee` decimal(5,2) NOT NULL,
  `currency_code` varchar(3) NOT NULL,
  `merchant_account_id` varchar(45) NOT NULL,
  `transaction_type` varchar(45) NOT NULL,
  `order_time` varchar(45) NOT NULL,
  `amount` decimal(5,2) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;


insert into ClientPlatform  (name, type, do_id_verification, do_bkg_check)
values ('USER', 'CONTRACTORS', 1, 1);

insert into ClientPlatform  (name, type, do_id_verification, do_bkg_check)
values ('THUMBTACK', 'CONTRACTORS', 1, 1);

insert into ClientPlatform  (name, type, do_id_verification, do_bkg_check)
values ('ANGIESLIST', 'CONTRACTORS', 1, 1);

-- Checkers

insert into Checker  (name, description)
values ('LN_BIZ_INST_ID', 'LN Business Instant Id');

insert into Checker  (name, description)
values ('LN_BIZ_ID', 'LN Business Id Check');

insert into Checker  (name, description)
values ('EXPERIAN_ID', 'Experian Id Check');

insert into Checker  (name, description)
values ('FACEBOOK', 'Facebook Check');

insert into Checker  (name, description)
values ('FULL_CONTACT_IND', 'Full Contact Individual Check');

insert into Checker  (name, description)
values ('FULL_CONTACT_BIZ', 'Full Contact Business Check');

insert into Checker  (name, description)
values ('LN_IND_ID', 'LN Individual Id Check');

insert into Checker  (name, description)
values ('TC_PROVIDER', 'Trustcloud Provider List Check');

insert into Checker  (name, description)
values ('YELP', 'Yelp Check');