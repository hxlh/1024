CREATE TABLE video1024.video_info (
	vid BIGINT auto_increment NOT NULL,
	uploader BIGINT NOT NULL,
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
