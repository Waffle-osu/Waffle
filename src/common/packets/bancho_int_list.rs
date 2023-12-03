use tokio::sync::mpsc::Sender;

use super::{BanchoSerializable, BanchoRequestType, BanchoPacket};

pub struct BanchoIntList {
    pub length: i16,
    pub list: Vec<i32>
}

impl BanchoSerializable for BanchoIntList {
    fn read(&mut self, reader: &mut binary_rw::BinaryReader) {
        self.length = reader.read_i16().expect("Failed to read BanchoIntList");

        for _ in 0..self.length {
            self.list.push(
                reader.read_i32().expect("Failed to read BanchoIntList element")
            );
        }
    }

    fn write(&self, writer: &mut binary_rw::BinaryWriter) {
        let _ = writer.write_i16(self.length);

        for i in 0..self.length {
            let _ = writer.write_i32(self.list[i as usize]);
        }
    }
}

impl BanchoIntList {
    pub fn send(packet_id: BanchoRequestType, list: Vec<i32>) -> BanchoPacket {
        return BanchoPacket::from_serializable(
            packet_id, 
            &BanchoIntList {
                length: list.len() as i16,
                list
            }
        );
    }

    pub async fn send_queue(queue: &Sender<BanchoPacket>, packet_id: BanchoRequestType, list: Vec<i32>) {
        let _ = queue.send(
            BanchoIntList::send(packet_id, list)
        ).await;
    }

    pub async fn self_send(&self, packet_id: BanchoRequestType) -> BanchoPacket {
        BanchoIntList::send(packet_id, self.list.clone())
    }

    pub async fn self_send_queue(&self, queue: &Sender<BanchoPacket>, packet_id: BanchoRequestType) {
        let _ = BanchoIntList::send_queue(queue, packet_id, self.list.clone());
    }
}