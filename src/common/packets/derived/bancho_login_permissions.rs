
use tokio::sync::mpsc::Sender;

use crate::packets::{BanchoPacket, BanchoInt, BanchoRequestType};

pub struct BanchoLoginPermissions {

}

impl BanchoLoginPermissions {
    pub async fn send(queue: &Sender<BanchoPacket>, permissions: i32) {
        BanchoInt::send_queue(queue, BanchoRequestType::BanchoLoginPermissions, permissions).await;
    }
}