use tokio::sync::mpsc::Sender;

use super::{BanchoSerializable, BanchoRequestType, BanchoPacket, bancho_string::BanchoString, read_bancho_string, write_bancho_string, bancho_status_update::BanchoStatusUpdate};

pub struct BanchoUserStats {
    pub user_id: i32,
    pub status: BanchoStatusUpdate,
    pub ranked_score: i64,
    pub accuracy: f32,
    pub playcount: i32,
    pub total_score: i64,
    pub rank: i32,
}

impl BanchoSerializable for BanchoUserStats {
    fn read(&mut self, reader: &mut binary_rw::BinaryReader) {
        let msg = "Failed to read BanchoUserStats";

        self.user_id = reader.read_i32().expect(msg);
        self.status.read(reader);
        self.ranked_score = reader.read_i64().expect(msg);
        self.accuracy = reader.read_f32().expect(msg);
        self.playcount = reader.read_i32().expect(msg);
        self.total_score = reader.read_i64().expect(msg);
        self.rank = reader.read_i32().expect(msg);
    }

    fn write(&self, writer: &mut binary_rw::BinaryWriter) {
        let msg = "Failed to write BanchoUserStats";

        writer.write_i32(self.user_id).expect(msg);
        self.status.write(writer);
        writer.write_i64(self.ranked_score).expect(msg);
        writer.write_f32(self.accuracy).expect(msg);
        writer.write_i32(self.playcount).expect(msg);
        writer.write_i64(self.total_score).expect(msg);
        writer.write_i32(self.rank).expect(msg);
    }
}

impl BanchoUserStats {
    pub fn send(packet_id: BanchoRequestType, presence: &BanchoUserStats) -> BanchoPacket {
        return BanchoPacket::from_serializable(
            packet_id, 
            presence
        );
    }

    pub async fn send_queue(queue: &Sender<BanchoPacket>, packet_id: BanchoRequestType, presence: &BanchoUserStats) {
        let _ = queue.send(
            BanchoUserStats::send(packet_id, presence)
        ).await;
    }

    pub async fn self_send(&self, packet_id: BanchoRequestType) -> BanchoPacket {
        BanchoUserStats::send(packet_id, self)
    }

    pub async fn self_send_queue(&self, queue: &Sender<BanchoPacket>, packet_id: BanchoRequestType) {
        let _ = BanchoUserStats::send_queue(queue, packet_id, self);
    }
}