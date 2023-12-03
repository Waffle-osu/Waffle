use tokio::sync::mpsc::Sender;

use crate::packets::{BanchoPacket, bancho_user_stats::BanchoUserStats, BanchoRequestType};

pub struct BanchoStatsUpdate {

}

impl BanchoStatsUpdate {
    pub async fn send(queue: &Sender<BanchoPacket>, stats: &BanchoUserStats) {
        BanchoUserStats::send_queue(queue, BanchoRequestType::BanchoHandleOsuUpdate, &stats).await;
    }
}