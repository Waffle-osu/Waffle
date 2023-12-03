use std::sync::Arc;

use common::db;

use crate::{osu, irc};

pub trait WaffleClient: Send + Sync  {
    fn get_user(&self) -> db::User;
}