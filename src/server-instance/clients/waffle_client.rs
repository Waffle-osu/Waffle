use std::sync::Arc;

use common::db;

use crate::{osu, irc};

pub enum WaffleClient {
    Osu(Arc<dyn osu::OsuClient + Sync + Send>),
    Irc(irc::IrcClient)
}