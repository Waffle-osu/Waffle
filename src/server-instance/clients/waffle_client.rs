use std::sync::Arc;

use common::db;

use crate::{osu, irc};

pub trait WaffleClient {
    fn get_user(&self) -> db::User;
}