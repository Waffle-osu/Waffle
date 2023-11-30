use common::db;

pub trait OsuClient {
    fn get_user(&self) -> db::User;
}