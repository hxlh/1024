CREATE DATABASE video1024 CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

CREATE TABLE video1024.video_info (
	vid BIGINT UNSIGNED auto_increment NOT NULL,
	uploader BIGINT UNSIGNED NOT NULL,
	cdn TEXT NOT NULL,
	subtitled varchar(100) NULL,
	likes BIGINT NOT NULL,
	tags TEXT NULL,
	CONSTRAINT video_info_PK PRIMARY KEY (vid)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_general_ci;
CREATE INDEX video_info_uploader_IDX USING BTREE ON video1024.video_info (uploader);

CREATE TABLE video1024.account (
	uid BIGINT UNSIGNED auto_increment NOT NULL,
	username varchar(20) NOT NULL,
	nickname varchar(20) NOT NULL,
	pwd TEXT NOT NULL,
	register_time DATETIME DEFAULT NOW() NOT NULL,
	avatar TEXT NULL,
	CONSTRAINT account_PK PRIMARY KEY (uid)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_general_ci;
CREATE UNIQUE INDEX account_username_IDX USING BTREE ON video1024.account (username);
