use tokio::sync::mpsc::Sender;

use crate::packets::{bancho_string::BanchoString, BanchoRequestType};

pub struct BanchoAnnounce {

}

impl BanchoAnnounce {
    pub async fn send(queue: &Sender<Vec<u8>>, message: String) {
        BanchoString::send_queue(queue, BanchoRequestType::BanchoAnnounce, &message).await;
    }
}