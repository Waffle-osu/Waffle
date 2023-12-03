CREATE TABLE osu_users (
    user_id        BIGINT      UNSIGNED NOT NULL AUTO_INCREMENT,
    username       VARCHAR(32)          NOT NULL,
    password       VARCHAR(64)          NOT NULL,
    country        SMALLINT    UNSIGNED NOT NULL DEFAULT '0',
    banned         TINYINT              NOT NULL DEFAULT '0',
    banned_reason  VARCHAR(512)         NOT NULL DEFAULT 'no reason',
    privileges     INT                  NOT NULL DEFAULT '1',
    joined_at      DATETIME                      DEFAULT CURRENT_TIMESTAMP,
    silenced_until BIGINT               NOT NULL DEFAULT 0,

    PRIMARY KEY (user_id, username),

    -- No duplicate user_id's or usernames.
    UNIQUE KEY id_UNIQUE (user_id),
    UNIQUE KEY username_UNIQUE (username),

    -- Quick lookups via username or user_id
    KEY user_INDEX (username, user_id)
) DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

CREATE TABLE osu_stats (
    user_id          BIGINT  UNSIGNED NOT NULL,
    mode             TINYINT          NOT NULL,
    ranked_score     BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    total_score      BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    user_level       DOUBLE           NOT NULL DEFAULT '0',
    accuracy         FLOAT            NOT NULL DEFAULT '0',
    playcount        BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    count_ssh        BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    count_ss         BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    count_sh         BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    count_s          BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    count_a          BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    count_b          BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    count_c          BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    count_d          BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    hit_300          BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    hit_100          BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    hit_50           BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    hit_miss         BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    hit_geki         BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    hit_katu         BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    replays_watched  BIGINT  UNSIGNED NOT NULL DEFAULT '0',
    playtime         BIGINT  UNSIGNED NOT NULL DEFAULT '0',

    PRIMARY KEY (user_id, mode)
) WITH SYSTEM VERSIONING DEFAULT CHARSET=utf8mb4;

INSERT INTO osu_users (username, password) VALUES ("WaffleBot", "no!");
INSERT INTO osu_stats (user_id, mode) VALUES (1, 0), (1, 1), (1, 2), (1, 3);

CREATE TABLE osu_friends (
    user_1 BIGINT UNSIGNED NOT NULL,
    user_2 BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (user_1, user_2),

    KEY index_user1 (user_1),
    KEY index_user2 (user_2)
) WITH SYSTEM VERSIONING;

CREATE TABLE osu_screenshots (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    filename VARCHAR(256) NOT NULL,
    creation_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id)
) WITH SYSTEM VERSIONING;

CREATE TABLE osu_beatmaps (
    beatmap_id     INT          NOT NULL AUTO_INCREMENT,
    beatmapset_id  INT          NOT NULL,
    creator_id     BIGINT       NOT NULL,
    filename       VARCHAR(256) NOT NULL,
    beatmap_md5    VARCHAR(32)  NOT NULL,
    version        VARCHAR(128) NOT NULL,
    total_length   INT          NOT NULL,
    drain_time     INT          NOT NULL,
    count_objects  INT          NOT NULL,
    count_normal   INT          NOT NULL,
    count_slider   INT          NOT NULL,
    count_spinner  INT          NOT NULL,
    diff_hp        TINYINT      NOT NULL,
    diff_cs        TINYINT      NOT NULL,
    diff_od        TINYINT      NOT NULL,
    diff_stars     FLOAT        NOT NULL,
    playmode       TINYINT      NOT NULL,
    ranking_status TINYINT      NOT NULL,
    last_update    DATETIME     NOT NULL,
    submit_date    DATETIME     NOT NULL,
    approve_date   DATETIME     NOT NULL,
    beatmap_source TINYINT      NOT NULL COMMENT 'Where the beatmap is from, used for when i may or may not add BSS and might be used for oldsu maps',
    
    PRIMARY KEY (beatmap_id),
    
    KEY checksum_index (beatmap_md5),
    KEY filename_index (filename),
    KEY setid_index    (beatmap_id),
    KEY source_index   (beatmap_source)
) DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

CREATE TABLE osu_beatmapsets (
		beatmapset_id  BIGINT        NOT NULL AUTO_INCREMENT,
		creator_id     BIGINT        NOT NULL,
		artist         VARCHAR(256)  NOT NULL,
		title          VARCHAR(256)  NOT NULL,
		creator        VARCHAR(256)  NOT NULL,
		source         VARCHAR(256)  NOT NULL,
		tags           VARCHAR(1024) NOT NULL,
		has_video      TINYINT       NOT NULL,
		has_storyboard TINYINT       NOT NULL,
		bpm            FLOAT         NOT NULL,
		
		PRIMARY KEY (beatmapset_id),
		
		KEY title_index (title),
		KEY artist_index (artist),
		KEY creator_index (creator),
		KEY source_index (source),
		KEY tags_index (tags)
) DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

-- There are 2 seperate tables because we need one to prevent multiple entries from a single user.

CREATE TABLE osu_beatmap_ratings (
    beatmapset_id BIGINT NOT NULL,
    rating_sum    BIGINT NOT NULL DEFAULT '0',
    votes         BIGINT NOT NULL DEFAULT '0',

    PRIMARY KEY (beatmapset_id)
) WITH SYSTEM VERSIONING;

CREATE TABLE osu_beatmap_rating_submissions (
    submission_id BIGINT          NOT NULL AUTO_INCREMENT,
    user_id       BIGINT UNSIGNED NOT NULL,
    beatmapset_id BIGINT          NOT NULL,

    PRIMARY KEY (submission_id, user_id, beatmapset_id)
) WITH SYSTEM VERSIONING;

CREATE TABLE osu_beatmap_offsets (
    offset_id    BIGINT NOT NULL AUTO_INCREMENT,
    beatmap_id   BIGINT NOT NULL,
    audio_offset INT    NOT NULL DEFAULT '0',

    PRIMARY KEY (offset_id, beatmap_id)
) WITH SYSTEM VERSIONING;

CREATE TABLE osu_beatmap_favourites (
    favourite_id  BIGINT          NOT NULL AUTO_INCREMENT,
    beatmapset_id BIGINT          NOT NULL,
    user_id       BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY(favourite_id, beatmapset_id, user_id),

    KEY userid (user_id)
) WITH SYSTEM VERSIONING;

CREATE TABLE osu_beatmap_comments (
    comment_id    BIGINT          NOT NULL AUTO_INCREMENT,
    user_id       BIGINT UNSIGNED NOT NULL,
    beatmap_id    BIGINT          NOT NULL,
    beatmapset_id BIGINT          NOT NULL,
    score_id      BIGINT          NOT NULL,
    time          BIGINT          NOT NULL,
    target        TINYINT         NOT NULL,
    comment       TEXT            NOT NULL,
    format_string VARCHAR(16)     NOT NULL,

    PRIMARY KEY(comment_id),

    KEY beatmapid (beatmap_id, target),
    KEY beatmapsetid (beatmapset_id, target),
    KEY scoreid (score_id, target),
    KEY userid (user_id)
) WITH SYSTEM VERSIONING;

CREATE TABLE osu_scores (
    score_id         BIGINT     UNSIGNED NOT NULL AUTO_INCREMENT,
    beatmap_id       BIGINT              NOT NULL,
    beatmapset_id    BIGINT              NOT NULL,
    user_id          BIGINT              NOT NULL,
    playmode         TINYINT             NOT NULL,
    score            INT                 NOT NULL DEFAULT '0',
    max_combo        INT                 NOT NULL DEFAULT '0',
    ranking          VARCHAR(2)          NOT NULL,
    hit300           INT                 NOT NULL DEFAULT '0',
    hit100           INT                 NOT NULL DEFAULT '0',
    hit50            INT                 NOT NULL DEFAULT '0',
    hitMiss          INT                 NOT NULL DEFAULT '0',
    hitGeki          INT                 NOT NULL DEFAULT '0',
    hitKatu          INT                 NOT NULL DEFAULT '0',
    enabled_mods     INT                 NOT NULL DEFAULT '0',
    perfect          TINYINT             NOT NULL DEFAULT '0',
    passed           TINYINT             NOT NULL DEFAULT '0',
    date             DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    leaderboard_best TINYINT             NOT NULL DEFAULT '0',
    mapset_best      TINYINT             NOT NULL DEFAULT '0',
    score_hash       VARCHAR(64)         NOT NULL DEFAULT '0',

    PRIMARY KEY (score_id),

    KEY scoreid (score_id),
    KEY leaderboardbest (leaderboard_best),
    KEY mapsetbest (mapset_best),
    KEY scorehash (score_hash),
    KEY beatmapid (beatmap_id),
    KEY userid (user_id)
) WITH SYSTEM VERSIONING;

CREATE TABLE osu_failtimes (
    failtime_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    failtime BIGINT NOT NULL,
    beatmap_id BIGINT NOT NULL,
    score_id BIGINT UNSIGNED NOT NULL,
    was_exit TINYINT NOT NULL,

    PRIMARY KEY (failtime_id),

    KEY beatmapid (beatmap_id)
) WITH SYSTEM VERSIONING;

CREATE TABLE irc_log (
    message_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    sender_id BIGINT UNSIGNED NOT NULL,
    target VARCHAR(64) NOT NULL,
    message TEXT NOT NULL,
    date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (message_id)
) WITH SYSTEM VERSIONING;

CREATE TABLE osu_updater_items (
    item_id BIGINT NOT NULL AUTO_INCREMENT,
    server_filename VARCHAR(256) NOT NULL,
    client_filename VARCHAR(256) NOT NULL,
    file_hash VARCHAR(64) NOT NULL,
    item_name VARCHAR(256) NOT NULL DEFAULT '',
    item_action VARCHAR(8) NOT NULL DEFAULt 'none',

    PRIMARY KEY (item_id, server_filename, client_filename)
) DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

INSERT INTO osu_updater_items (server_filename, client_filename, file_hash, item_name, item_action) VALUES
    ('osu!.exe'                   , 'osu!.exe'                   , '3623d9f7c693b786564e2d61b1c43af9', 'client_debug', 'none'),
    ('avcodec-51.dll'             , 'avcodec-51.dll'             , 'b22bf1e4ecd4be3d909dc68ccab74eec', 'client_ddls',  'none'),
    ('avformat-52.dll'            , 'avformat-52.dll'            , '2e7a800133625f827cf46aa0bb1af800', 'client_ddls',  'none'),
    ('avutil-49.dll'              , 'avutil-49.dll'              , 'c870147dff89c95c81f8fbdfbc6344ac', 'client_ddls',  'none'),
    ('bass.dll'                   , 'bass.dll'                   , 'bbfc7d855252b0211875769bbf667bcd', 'client_ddls',  'none'),
    ('bass_fx.dll'                , 'bass_fx.dll'                , 'f9ffe0a23a32b79653e31330764c4231', 'client_ddls',  'none'),
    ('d3dx9_31.dll'               , 'd3dx9_31.dll'               , '797e24743937d67d69f28f2cf5052ee8', 'client_ddls',  'none'),
    ('Microsoft.Ink.dll'          , 'Microsoft.Ink.dll'          , 'a02ee61542caae25f8a44c9428d30247', 'client_ddls',  'none'),
    ('Microsoft.Xna.Framework.dll', 'Microsoft.Xna.Framework.dll', '45a786658d3f69717652fed471d03ee0', 'client_ddls',  'none'),
    ('osu!common.dll'             , 'osu!common.dll'             , '820817b776374a0adcbb7fa2a7ca74f2', 'client_ddls',  'none'),
    ('osu.dll.zip'                , 'osu.dll'                    , '599c14bcfc9c43b88d70d1a9388b33b7', 'client_ddls',  'zip' ),
    ('OsuP2P.dll'                 , 'OsuP2P.dll'                 , '2342bfd835e2e487d040d8ba62eb1a72', 'client_ddls',  'none'),
    ('pthreadGC2.dll'             , 'pthreadGC2.dll'             , 'ce931021e18f385f519e945a8a10548e', 'client_ddls',  'none'),
    ('x3daudio1_1.dll'            , 'x3daudio1_1.dll'            , '121b131eaa369d8f58dacc5c39a77d80', 'client_ddls',  'none');

CREATE TABLE osu_achievements (
    achivement_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    image VARCHAR(64) NOT NULL,

    PRIMARY KEY (achivement_id, name, image)
) DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;

CREATE TABLE osu_achieved_achievement (
    user_achievement_id BIGINT NOT NULL AUTO_INCREMENT,
    achievement_id BIGINT NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,

    PRIMARY KEY (user_achievement_id)
) WITH SYSTEM VERSIONING;

INSERT INTO osu_achievements (name, image) VALUES
    ("500 Combo",                         "combo500.png"),
    ("750 Combo",                         "combo750.png"),
    ("1000 Combo",                        "combo1000.png"),
    ("2000 Combo",                        "combo2000.png"),
    ("Video Game Pack vol.1",             "gamer1.png"),
    ("Video Game Pack vol.2",             "gamer2.png"),
    ("Video Game Pack vol.3",             "gamer3.png"),
    ("Video Game Pack vol.4",             "gamer4.png"),
    ("Anime Pack vol.1",                  "anime1.png"),
    ("Anime Pack vol.2",                  "anime2.png"),
    ("Anime Pack vol.3",                  "anime3.png"),
    ("Anime Pack vol.4",                  "anime4.png"),
    ("Internet! Pack vol.1",              ".png"),
    ("Internet! Pack vol.2",              ".png"),
    ("Internet! Pack vol.3",              ".png"),
    ("Internet! Pack vol.4",              ".png"),
    ("Rhythm Game Pack vol.1",            "rhythm1.png"),
    ("Rhythm Game Pack vol.2",            "rhythm2.png"),
    ("Rhythm Game Pack vol.3",            "rhythm3.png"),
    ("Rhythm Game Pack vol.4",            "rhythm4.png"),
    ("Catch 20000 Fruits",                "fruitsalad.png"),
    ("Catch 200000 Fruits",               "fruitplatter.png"),
    ("Catch 2000000 Fruits",              "fruitod.png"),
    ("5000 Plays",                        "plays1.png"),
    ("15000 Plays",                       "plays2.png"),
    ("25000 Plays",                       "plays3.png"),
    ("50000 Plays",                       "plays4.png"),
    ("30000 Drum Hits",                   "taiko1.png"),
    ("300000 Drum Hits",                  "taiko2.png"),
    ("3000000 Drum Hits",                 "taiko3.png"),
    ("Don't let the bunny distract you!", "bunny.png"),
    ("S-Ranker",                          "s-ranker.png"),
    ("Most Improved",                     "improved.png"),
    ("Non-stop Dancer",                   "dancer.png");

CREATE TABLE web_tokens (
    token_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    token_hash VARCHAR(128) NOT NULL,
    creation_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (token_id, token_hash)
) DEFAULT CHARSET=utf8mb4 WITH SYSTEM VERSIONING;