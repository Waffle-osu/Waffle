use tokio::sync::mpsc::Sender;

use super::{BanchoSerializable, BanchoRequestType, BanchoPacket, bancho_string::BanchoString, read_bancho_string, write_bancho_string};

pub struct BanchoPresence {
    pub user_id: i32,
    pub username: Option<String>,
    pub avatar_extension: u8,
    pub timezone: u8,
    pub country: u8,
    pub city: Option<String>,
    pub permissions: u8,
    pub longitude: f32,
    pub latitude: f32,
    pub rank: i32,
}

impl BanchoSerializable for BanchoPresence {
    fn read(&mut self, reader: &mut binary_rw::BinaryReader) {
        let msg = "Failed to read BanchoPresence";

        self.user_id = reader.read_i32().expect(msg);
        self.username = read_bancho_string(reader);
        self.avatar_extension = reader.read_u8().expect(msg);
        self.timezone = reader.read_u8().expect(msg);
        self.country = reader.read_u8().expect(msg);
        self.city = read_bancho_string(reader);
        self.permissions = reader.read_u8().expect(msg);
        self.longitude = reader.read_f32().expect(msg);
        self.latitude = reader.read_f32().expect(msg);
        self.rank = reader.read_i32().expect(msg);
    }

    fn write(&self, writer: &mut binary_rw::BinaryWriter) {
        let msg = "Failed to write BanchoPresence";

        writer.write_i32(self.user_id).expect(msg);
        writer.write_bytes(write_bancho_string(&self.username)).expect(msg);
        writer.write_u8(self.avatar_extension).expect(msg);
        writer.write_u8(self.timezone).expect(msg);
        writer.write_u8(self.country).expect(msg);
        writer.write_bytes(write_bancho_string(&self.city)).expect(msg);
        writer.write_u8(self.permissions).expect(msg);
        writer.write_f32(self.longitude).expect(msg);
        writer.write_f32(self.latitude).expect(msg);
        writer.write_i32(self.rank).expect(msg);
    }
}

impl BanchoPresence {
    pub fn send(packet_id: BanchoRequestType, presence: &BanchoPresence) -> BanchoPacket {
        return BanchoPacket::from_serializable(
            packet_id, 
            presence
        );
    }

    pub async fn send_queue(queue: &Sender<BanchoPacket>, packet_id: BanchoRequestType, presence: &BanchoPresence) {
        let _ = queue.send(
            BanchoPresence::send(packet_id, presence)
        ).await;
    }

    pub async fn self_send(&self, packet_id: BanchoRequestType) -> BanchoPacket {
        BanchoPresence::send(packet_id, self)
    }

    pub async fn self_send_queue(&self, queue: &Sender<BanchoPacket>, packet_id: BanchoRequestType) {
        let _ = BanchoPresence::send_queue(queue, packet_id, self);
    }
}