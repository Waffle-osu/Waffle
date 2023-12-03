use binary_rw::{MemoryStream, BinaryWriter, Endian};
use tokio::sync::mpsc::Sender;

use super::{BanchoSerializable, BanchoRequestType, BanchoPacketHeader, BanchoPacket};

pub fn write_bancho_string(string: &String) -> Vec<u8> {
    let mut output = Vec::new();

    if string == "" {
        return output;
    }

    let mut length = 0;
    let mut i = string.len();
    let mut ulebBytes: Vec<u8> = Vec::new();

    while i > 0 {
        ulebBytes.push(0);
        ulebBytes[length] = (i & 0x7F) as u8;

        i >>= 7;

        if i != 0 {
            ulebBytes[length] |= 0x80;
        }

        length += 1;
    }

    output.push(11);
    output.append(&mut ulebBytes);
    
    let mut str_to_vec = Vec::from(string.as_bytes());

    output.append(&mut str_to_vec);

    return output;
}

pub struct BanchoString {
    string: String
}

impl BanchoSerializable for BanchoString {
    fn read(&mut self, reader: &mut binary_rw::BinaryReader) {
        let initial_byte_result = reader.read_u8();

        if initial_byte_result.is_err() {
            return;
        }

        let initial_byte = initial_byte_result.unwrap();

        if initial_byte != 11 {
            return;
        }

        let mut shift: u32 = 0;
        let mut last_byte: u8 = 0;
        let mut total: u32 = 0;

        loop {
            let read_result = reader.read_u8();

            if read_result.is_err() {
                return;
            }

            last_byte = read_result.unwrap();

            total |= ((last_byte & 0x7F) as u32) << shift;

            if last_byte & 0x80 == 0 {
                break;
            }

            shift += 7;
        }

        let bytes = reader.read_bytes(total as usize);

        if bytes.is_err() {
            return;
        }

        self.string = String::from_utf8(bytes.unwrap()).unwrap();
    }
    
    fn write(&self, writer: &mut binary_rw::BinaryWriter) {
        let bancho_string = write_bancho_string(&self.string);

        let _ = writer.write_bytes(bancho_string);
    }
}

impl BanchoString {
    pub fn send(packet_id: BanchoRequestType, string: &String) -> BanchoPacket {
        return BanchoPacket::from_data(
            packet_id, 
            write_bancho_string(string)
        );
    }

    pub async fn send_queue(queue: &Sender<BanchoPacket>, packet_id: BanchoRequestType, string: &String) {
        let _ = queue.send(
            BanchoString::send(packet_id, string)
        ).await;
    }

    pub async fn self_send(&self, packet_id: BanchoRequestType) -> BanchoPacket {
        BanchoString::send(packet_id, &self.string)
    }

    pub async fn self_send_queue(&self, queue: &Sender<BanchoPacket>, packet_id: BanchoRequestType) {
        let _ = BanchoString::send_queue(queue, packet_id, &self.string);
    }
}
