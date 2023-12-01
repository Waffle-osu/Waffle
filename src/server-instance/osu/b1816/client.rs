use chrono::{DateTime, Utc};
use tokio::net::TcpStream;

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

    // spectators:
    // spectatingClient:

    // packetQueue: 
}

impl OsuClient {

}