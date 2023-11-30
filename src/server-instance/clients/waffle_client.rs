use common::db;

pub trait WaffleClient {
    fn get_user(&self) -> db::User;
}