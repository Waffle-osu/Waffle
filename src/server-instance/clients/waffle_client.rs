use std::sync::Arc;

use crate::{osu, irc};

pub enum WaffleClient {
    Osu(Arc<dyn osu::OsuClient>),
    Irc(irc::IrcClient),
}