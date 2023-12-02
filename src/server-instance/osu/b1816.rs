mod client_auth;
mod client;

use common::db::{self, User};
use sqlx::types::time::PrimitiveDateTime;

use super::OsuClient;

pub struct osub1816 {

}

impl OsuClient for osub1816 {
    fn get_user(&self) -> db::User {
        User { user_id: 0, username: "()".to_string(), password: "()".to_string(), country: "()".to_string(), silenced_until: 0, banned: false, banned_reason: "()".to_string(), privileges: 0, joined_at: PrimitiveDateTime::MIN }
    }
}