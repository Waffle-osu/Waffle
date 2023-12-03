use binary_rw::BinaryReader;
use tokio::sync::mpsc::Sender;

use super::{BanchoSerializable, BanchoRequestType, BanchoPacket};

pub fn write_bancho_string(string: &Option<String>) -> Vec<u8> {
    let mut output = Vec::new();

    if string.is_none() {
        return vec![0];
    }

    let string = string.as_ref().unwrap();

    if string == "" {
        return vec![0];
    }

    let mut length = 0;
    let mut i = string.len();
    let mut uleb_bytes: Vec<u8> = Vec::new();

    while i > 0 {
        uleb_bytes.push(0);
        uleb_bytes[length] = (i & 0x7F) as u8;

        i >>= 7;

        if i != 0 {
            uleb_bytes[length] |= 0x80;
        }

        length += 1;
    }

    output.push(11);
    output.append(&mut uleb_bytes);
    
    let mut str_to_vec = Vec::from(string.as_bytes());

    output.append(&mut str_to_vec);

    return output;
}

pub fn read_bancho_string(reader: &mut BinaryReader) -> Option<String> {
    let initial_byte_result = reader.read_u8();

    if initial_byte_result.is_err() {
        return None;
    }

    let initial_byte = initial_byte_result.unwrap();

    if initial_byte != 11 {
        return None;
    }

    let mut shift: u32 = 0;
    let mut total: u32 = 0;

    loop {
        let read_result = reader.read_u8();

        if read_result.is_err() {
            return None;
        }

        let last_byte = read_result.unwrap();

        total |= ((last_byte & 0x7F) as u32) << shift;

        if last_byte & 0x80 == 0 {
            break;
        }

        shift += 7;
    }

    let bytes = reader.read_bytes(total as usize);

    if bytes.is_err() {
        return None;
    }

    return Some(String::from_utf8(bytes.unwrap()).unwrap());
}

pub struct BanchoString {
    string: Option<String>
}


impl BanchoSerializable for BanchoString {
    fn read(&mut self, reader: &mut binary_rw::BinaryReader) {
        self.string = read_bancho_string(reader);
    }
    
    fn write(&self, writer: &mut binary_rw::BinaryWriter) {
        let bancho_string = write_bancho_string(&self.string);

        let _ = writer.write_bytes(bancho_string);
    }
}

impl BanchoString {
    pub fn send(packet_id: BanchoRequestType, string: &Option<String>) -> BanchoPacket {
        return BanchoPacket::from_data(
            packet_id, 
            write_bancho_string(string)
        );
    }

    pub async fn send_queue(queue: &Sender<BanchoPacket>, packet_id: BanchoRequestType, string: &Option<String>) {
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
