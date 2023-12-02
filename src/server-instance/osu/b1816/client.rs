use std::sync::Arc;

use chrono::{DateTime, Utc};
use common::packets::BanchoPacket;
use dashmap::DashMap;
use tokio::{net::TcpStream, sync::mpsc::{Sender, Receiver}};

use crate::clients::{self, waffle_client::WaffleClient};

pub struct ClientInformation {
    pub version: i32,
    pub client_hash: String,
    pub allow_city: bool,
    pub mac_address_hash: String,
    pub timezone: i32
}

pub struct OsuClient {
    connection: TcpStream,
    continue_running: bool,

    logon_time: DateTime<Utc>,

    last_receive: DateTime<Utc>,
    last_ping: DateTime<Utc>,

    // joinedChannels: 
    away_message: String,

    spectators: DashMap<u64, Arc<WaffleClient>>,
    spectatingClient: Arc<WaffleClient>,

    packetQueueSend: Arc<Sender<BanchoPacket>>, 
    packetQueueRecv: Arc<Receiver<BanchoPacket>>
}

impl OsuClient {

}