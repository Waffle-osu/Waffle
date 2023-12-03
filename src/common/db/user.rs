use std::{sync::Arc, ops::Deref};

use sqlx::{MySqlPool, mysql::MySqlQueryResult, Error};

#[derive(sqlx::FromRow, Clone)]
pub struct User {
    pub user_id: u64,
    pub username: String,
    pub password: String,
    pub country: u16,
    pub silenced_until: i64,
    pub banned: bool,
    pub banned_reason: String,
    pub privileges: i32,
    pub joined_at: sqlx::types::time::PrimitiveDateTime
}

impl User {
    pub async fn from_id(pool: Arc<MySqlPool>, id: u64) -> Option<User> {
        let row = 
            sqlx::query_as("SELECT * FROM osu_users WHERE user_id = ?")
                .bind(id)
                .fetch_one(pool.deref())
                .await;

        match row {
            Ok(user) => return Some(user),
            Err(err) => {
                println!("{}", err.to_string());
                
                return None;
            },
        };
    }

    pub async fn from_username(pool: Arc<MySqlPool>, username: String) -> Option<User> {
        let row = 
            sqlx::query_as("SELECT * FROM osu_users WHERE username = ?")
                .bind(username)
                .fetch_one(pool.deref())
                .await;

        match row {
            Ok(user) => return Some(user),
            Err(err) => {
                println!("{}", err.to_string());
                
                return None;
            },
        };
    }

    pub async fn create_user(pool: Arc<MySqlPool>, username: String, password: String) -> Result<(),()> {
        let bcrypt_pw = bcrypt::hash(password, bcrypt::DEFAULT_COST);

        if bcrypt_pw.is_err() {
            return Err(());
        }

        let bcrypt_string = bcrypt_pw.unwrap();

        let user_insert_result = sqlx::query("INSERT INTO osu_users (username, password) VALUES (?, ?)")
            .bind(username)
            .bind(bcrypt_string)
            .execute(pool.deref())
            .await;

        if user_insert_result.is_err() {
            let errstr = user_insert_result.err().unwrap().to_string();
            println!("{}", errstr);
            return Err(());
        }

        let user_id = user_insert_result.unwrap().last_insert_id();

        let stat_insert_result = sqlx::query("INSERT INTO osu_stats (user_id, mode) VALUES (?, 0), (?, 1), (?, 2), (?, 3)")
            .bind(user_id)
            .bind(user_id)
            .bind(user_id)
            .bind(user_id)
            .execute(pool.deref())
            .await;

        if stat_insert_result.is_err() {
            return Err(())
        }

        return Ok(())
    }
}