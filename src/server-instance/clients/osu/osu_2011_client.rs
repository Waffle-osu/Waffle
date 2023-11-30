use tokio::net::TcpStream;

use super::generic_osu_client::GenericOsuClient;

struct Osu2011Client {
    net_stream: TcpStream
}

impl Osu2011Client {
    fn from_generic_client(client: GenericOsuClient) {
        
    }
}