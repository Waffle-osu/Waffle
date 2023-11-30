use sqlx::MySqlPool;

use super::db_pool;

#[derive(sqlx::FromRow)]
pub struct User {
    user_id: u64,
    username: String,
    password: String,
    country: String,
    silenced_until: u64,
    banned: bool,
    banned_reason: String,
    privileges: i32,
    joined_at: sqlx::types::time::PrimitiveDateTime
}

impl User {
    async fn from_id(pool: &MySqlPool, id: u64) -> Option<User> {
        let row = 
            sqlx::query_as("SELECT * FROM osu_users WHERE user_id = $1")
                .bind(id)
                .fetch_one(pool)
                .await;

        match row {
            Ok(user) => return Some(user),
            Err(_) => return None,
        };
    }

    async fn from_username(pool: &MySqlPool, username: String) -> Option<User> {
        let row = 
            sqlx::query_as("SELECT * FROM osu_users WHERE username = $1")
                .bind(username)
                .fetch_one(pool)
                .await;

        match row {
            Ok(user) => return Some(user),
            Err(_) => return None,
        };
    }
}