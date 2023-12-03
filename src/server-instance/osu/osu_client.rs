use common::db;

pub trait OsuClient: Send + Sync {
    fn get_user(&self) -> db::User;
}