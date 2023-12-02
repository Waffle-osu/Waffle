use binary_rw::{BinaryWriter, Endian, MemoryStream};
use tokio::sync::mpsc::Sender;

use super::{BanchoRequestType, BanchoSerializable, BanchoPacketHeader};

pub struct BanchoInt {
    number: i32
}

impl BanchoSerializable for BanchoInt {
    fn read(&mut self, reader: &mut binary_rw::BinaryReader) {
        self.number = reader.read_i32().expect("Failed to read BanchoInt");
    }
    
    fn write(&self, writer: &mut binary_rw::BinaryWriter) {
        writer.write_i32(self.number).expect("Failed to write BanchoInt");
    }
}

impl BanchoInt {
    pub fn send(packet_id: BanchoRequestType, number: i32) -> Vec<u8> {
        let mut buffer = MemoryStream::new();
        let mut writer = BinaryWriter::new(&mut buffer, Endian::Little);

        let header = BanchoPacketHeader {
            packet_id: packet_id,
            compressed: false,
            size: 4
        };

        header.write(&mut writer);

        let _ = writer.write_i32(number);

        return buffer.into();
    }

    pub async fn send_queue(queue: &Sender<Vec<u8>>, packet_id: BanchoRequestType, number: i32) {
        let _ = queue.send(
            BanchoInt::send(packet_id, number)
        ).await;
    }

    pub async fn self_send(&self, packet_id: BanchoRequestType) -> Vec<u8> {
        BanchoInt::send(packet_id, self.number)
    }

    pub async fn self_send_queue(&self, queue: &Sender<Vec<u8>>, packet_id: BanchoRequestType) {
        let _ = BanchoInt::send_queue(queue, packet_id, self.number);
    }
}

