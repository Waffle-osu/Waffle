use tokio::sync::mpsc::Sender;

use crate::packets::{BanchoRequestType, BanchoInt, BanchoPacket};

pub struct BanchoLoginReply {

}

impl BanchoLoginReply {
    pub async fn send_wrong_login(queue: &Sender<BanchoPacket>) {
        BanchoInt::send_queue(queue, BanchoRequestType::BanchoLoginReply, -1).await;
    }

    pub async fn send_wrong_version(queue: &Sender<BanchoPacket>) {
        BanchoInt::send_queue(queue, BanchoRequestType::BanchoLoginReply, -2).await;
    }

    pub async fn send_banned(queue: &Sender<BanchoPacket>) {
        BanchoInt::send_queue(queue, BanchoRequestType::BanchoLoginReply, -3).await;
    }

    pub async fn send_server_error(queue: &Sender<BanchoPacket>) {
        BanchoInt::send_queue(queue, BanchoRequestType::BanchoLoginReply, -5).await;
    }

    pub async fn send_success(queue: &Sender<BanchoPacket>, user_id: i32) {
        BanchoInt::send_queue(queue, BanchoRequestType::BanchoLoginReply, user_id).await;
    }
}