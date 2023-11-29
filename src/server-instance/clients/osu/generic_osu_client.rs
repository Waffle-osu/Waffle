use tokio::{net::{unix::SocketAddr, TcpStream}, io::AsyncReadExt};

pub struct GenericOsuClient {
    net_stream: TcpStream
}

pub async fn complete_initial_login(stream: TcpStream, addr: SocketAddr) {
    let mut buffer = [0; 256];

    stream
        .take(256)
        .read(&mut buffer)
        .await
        .expect("Failed to read initial string");

    let login_string = 
        std::str::from_utf8(&buffer)
            .expect("Invalid string received!")
            .to_string();

    let split: Vec<&str> = login_string.split('\n').collect();

    if split.len() != 3 {
        // disconnect?
    }

    let username = split[0];
    let password = split[1];
    let client_details = split[2];

    let client_details_split: Vec<&str> = client_details.split('|').collect();

    let version = client_details_split[0];
}