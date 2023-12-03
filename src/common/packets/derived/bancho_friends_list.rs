use tokio::sync::mpsc::Sender;

use crate::packets::{BanchoPacket, bancho_int_list::BanchoIntList, BanchoRequestType};

pub struct BanchoFriendsList {

}

impl BanchoFriendsList {
    pub async fn send(queue: &Sender<BanchoPacket>, list: Vec<i32>) {
        BanchoIntList::send_queue(queue, BanchoRequestType::BanchoFriendsList, list).await;
    }
}