use std::sync::{Arc};

use actix::{Actor, Context};
use chrono::{DateTime, Utc};
use common::{packets::BanchoPacket, db, send_box::SendBox};
use dashmap::DashMap;
use tokio::{net::TcpStream, sync::{mpsc::{Sender, Receiver}, Mutex}};

use crate::{clients::{waffle_client::WaffleClient}, osu::OsuClient};

pub struct ClientInformation {
    pub version: i32,
    pub client_hash: String,
    pub allow_city: bool,
    pub mac_address_hash: String,
    pub timezone: i32
}

pub struct OsuClient2011 {
    pub connection: TcpStream,
    pub continue_running: bool,

    pub logon_time: DateTime<Utc>,

    pub last_receive: DateTime<Utc>,
    pub last_ping: DateTime<Utc>,

    // joinedChannels: 
    pub away_message: String,

    pub spectators: DashMap<u64, WaffleClient>,
    pub spectating_client: Option<WaffleClient>,
    
    pub packet_queue_send: Mutex<Sender<BanchoPacket>>, 
    pub packet_queue_recv: Mutex<Receiver<BanchoPacket>>,

    pub user: db::User,
    pub osu_stats: db::Stats,
    pub taiko_stats: db::Stats,
    pub catch_stats: db::Stats,
    pub mania_stats: db::Stats,
}

impl OsuClient for OsuClient2011 {
    fn get_user(&self) -> db::User {
        self.user.clone()
    }
}

impl Actor for OsuClient2011 {
    type Context = Context<Self>;
}

