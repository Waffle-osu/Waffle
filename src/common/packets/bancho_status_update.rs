use tokio::sync::mpsc::Sender;

use super::{BanchoSerializable, BanchoRequestType, BanchoPacket, bancho_string::BanchoString, read_bancho_string, write_bancho_string};

pub struct BanchoStatusUpdate {
    pub status: u8,
    pub status_text: Option<String>,
    pub beatmap_checksum: Option<String>,
    pub current_mods: u16,
    pub play_mode: u8,
    pub beatmap_id: i32,
}

impl BanchoSerializable for BanchoStatusUpdate {
    fn read(&mut self, reader: &mut binary_rw::BinaryReader) {
        let msg = "Failed to read BanchoStatusUpdate";

        self.status = reader.read_u8().expect(msg);
        self.status_text = read_bancho_string(reader);
        self.beatmap_checksum = read_bancho_string(reader);
        self.current_mods = reader.read_u16().expect(msg);
        self.play_mode = reader.read_u8().expect(msg);
        self.beatmap_id = reader.read_i32().expect(msg);
    }

    fn write(&self, writer: &mut binary_rw::BinaryWriter) {
        let msg = "Failed to write BanchoStatusUpdate";

        writer.write_u8(self.status).expect(msg);
        writer.write_bytes(write_bancho_string(&self.status_text)).expect(msg);
        writer.write_bytes(write_bancho_string(&self.beatmap_checksum)).expect(msg);
        writer.write_u16(self.current_mods).expect(msg);
        writer.write_u8(self.play_mode).expect(msg);
        writer.write_i32(self.beatmap_id).expect(msg);
    }
}

impl BanchoStatusUpdate {
    pub fn send(packet_id: BanchoRequestType, presence: &BanchoStatusUpdate) -> BanchoPacket {
        return BanchoPacket::from_serializable(
            packet_id, 
            presence
        );
    }

    pub async fn send_queue(queue: &Sender<BanchoPacket>, packet_id: BanchoRequestType, presence: &BanchoStatusUpdate) {
        let _ = queue.send(
            BanchoStatusUpdate::send(packet_id, presence)
        ).await;
    }

    pub async fn self_send(&self, packet_id: BanchoRequestType) -> BanchoPacket {
        BanchoStatusUpdate::send(packet_id, self)
    }

    pub async fn self_send_queue(&self, queue: &Sender<BanchoPacket>, packet_id: BanchoRequestType) {
        let _ = BanchoStatusUpdate::send_queue(queue, packet_id, self);
    }
}