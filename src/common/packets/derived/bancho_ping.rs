use tokio::sync::mpsc::Sender;

use crate::packets::{BanchoPacket, BanchoPacketHeader, BanchoRequestType};

pub struct BanchoPing {

}

impl BanchoPing {
    pub async fn send(queue: &Sender<BanchoPacket>) {
        let _ = queue.send(
            BanchoPacket { 
                header: BanchoPacketHeader { 
                    packet_id: BanchoRequestType::BanchoPing, 
                    compressed: false, 
                    size: 0 
                }, 
                data: Vec::new() 
            }).await;
    }
}