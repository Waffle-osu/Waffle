use std::{sync::Arc, ops::Deref};

use sqlx::MySqlPool;

#[derive(sqlx::FromRow)]
pub struct Stats {
    pub user_id: u64,
    pub mode: u8,
    pub ranked_score: u64,
    pub total_score: u64,
    pub user_level: f64,
    pub accuracy: f32,
    pub playcount: u64,
    pub count_ssh: u64,
    pub count_ss: u64,
    pub count_sh: u64,
    pub count_s: u64,
    pub count_a: u64,
    pub count_b: u64,
    pub count_c: u64,
    pub count_d: u64,
    pub hit_300: u64,
    pub hit_100: u64,
    pub hit_50: u64,
    pub hit_miss: u64,
    pub hit_geki: u64,
    pub hit_katu: u64,
    pub replays_watched: u64,
    pub playtime: u64
}

impl Stats {
    pub async fn from_id(pool: Arc<MySqlPool>, user_id: u64, mode: u8) -> Option<Stats> {
        let row = 
            sqlx::query_as("SELECT * FROM osu_stats WHERE user_id = $1 AND mode = $2")
                .bind(user_id)
                .bind(mode)
                .fetch_one(pool.deref())
                .await;
        
        match row {
            Ok(stat) => return Some(stat),
            Err(_) => return None,
        };
    }
}