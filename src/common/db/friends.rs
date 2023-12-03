use std::{sync::Arc, ops::Deref};

use sqlx::MySqlPool;

#[derive(sqlx::FromRow)]
pub struct Friends {
    pub user_1: u64,
    pub user_2: u64,
}

impl Friends {
    pub async fn get_users_friends(pool: Arc<MySqlPool>, user_id: u64) -> Vec<Friends> {
        let rows = 
            sqlx::query_as("SELECT * FROM osu_friends WHERE user_1 = $1")
                .bind(user_id)
                .fetch_all(pool.deref())
                .await;

        match rows {
            Ok(friends) => return friends,
            Err(_) => return Vec::new(),
        }
    }
}