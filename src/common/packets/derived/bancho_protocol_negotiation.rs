use tokio::sync::mpsc::Sender;

use crate::packets::{BanchoPacket, BanchoInt, BanchoRequestType};

pub struct BanchoProtocolNegotiation {

}

impl BanchoProtocolNegotiation {
    pub async fn send(queue: &Sender<BanchoPacket>, version: i32) {
        BanchoInt::send_queue(queue, BanchoRequestType::BanchoProtocolNegotiation, version).await;
    }
}