use tokio::sync::mpsc::Sender;

use crate::packets::{bancho_string::BanchoString, BanchoRequestType, BanchoPacket};

pub struct BanchoAnnounce {

}

impl BanchoAnnounce {
    pub async fn send(queue: &Sender<BanchoPacket>, message: String) {
        BanchoString::send_queue(queue, BanchoRequestType::BanchoAnnounce, &Some(message)).await;
    }
}