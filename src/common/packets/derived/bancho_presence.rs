use tokio::sync::mpsc::Sender;

use crate::packets::{bancho_presence::{self, BanchoPresence}, BanchoPacket, BanchoRequestType};

pub struct BanchoUserPresence {

}

impl BanchoUserPresence {
    pub async fn send(queue: &Sender<BanchoPacket>, presence: &BanchoPresence) {
        BanchoPresence::send_queue(queue, BanchoRequestType::BanchoUserPresence, &presence).await;
    }
}