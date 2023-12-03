use std::{sync::Arc, ops::Deref, time::Duration};

use binary_rw::{MemoryStream, BinaryReader, Endian};
use common::packets::{BanchoPacketHeader, BanchoPacket, derived::BanchoPing};
use tokio::{io::{AsyncWriteExt, AsyncReadExt}, time::sleep};

use super::client::OsuClient2011;

impl OsuClient2011 {
    pub async fn maintain_client(client: Arc<OsuClient2011>) {
        while client.continue_running {
            sleep(Duration::from_millis(15000)).await;

            BanchoPing::send(&client.packet_queue_send.lock().await.to_owned()).await;
        }
    }

    pub async fn handle_incoming(client: Arc<OsuClient2011>) {
        while client.continue_running {
            let _ = client.connection.readable().await;

            let mut buffer = [0; 32768];

            let read_result = client.connection.try_read(&mut buffer);

            //There do come useless errors which shouldnt be accounted for
            if read_result.is_err() {
                continue;
            }

            let size = read_result.unwrap();

            let read_buffer: Vec<u8> = buffer.iter().cloned().take(size).collect();
            let mut memory_stream = MemoryStream::from(read_buffer);
            let mut binary_reader = BinaryReader::new(&mut memory_stream, Endian::Little);

            let mut read_index = 0;

            while read_index < size {
                let (packet, read) = BanchoPacket::read(&mut binary_reader);

                read_index += read;

                println!("Received Packet {:?} size {}", packet.header.packet_id, packet.header.size);
            }
        }
    }

    pub async fn send_outgoing(client: Arc<OsuClient2011>) {
        while client.continue_running {
            let packet = 
                client.packet_queue_recv
                    .lock()
                    .await
                    .recv()
                    .await;
            
            match packet {
                None => {},
                Some(packet) => {
                    println!("Writing Packet {:?} with size {}", packet.header.packet_id, packet.header.size);

                    let result = 
                        client.connection.try_write(
                            packet.send()
                            .as_slice()
                        );

                    if result.is_err() {
                        //TODO: error handling.
                    }
                }
            }

            sleep(Duration::from_millis(1)).await;
        }
    }
}