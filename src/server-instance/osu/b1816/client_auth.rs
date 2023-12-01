use std::sync::Arc;

use chrono::Utc;
use common::{packets::{BanchoInt, BanchoRequestType}, db};
use sqlx::MySqlPool;
use tokio::{net::{TcpStream, unix::SocketAddr}, sync::mpsc::{self, Receiver, Sender}, io::{BufReader, AsyncBufReadExt, AsyncWriteExt}};

use super::client::ClientInformation;

async fn send_and_close(connection: &mut TcpStream, queue: &mut Receiver<Vec<u8>>) {
    while let Some(message) = queue.recv().await {
        let buffer = message.as_slice();

        connection.write(buffer).await.expect("Failed to write packets!");
    }

    connection.flush().await.expect("Failed to flush packets");
    connection.shutdown().await.expect("Shutdown of the stream failed!");
}

async fn send_wrong_version(connection: &mut TcpStream, queue_send: &mut Sender<Vec<u8>>, queue_receive: &mut Receiver<Vec<u8>>) {
    BanchoInt::send_queue(queue_send, BanchoRequestType::BanchoLoginReply, -2).await;
        
    send_and_close(connection, queue_receive).await;
}

pub async fn handle_new_client(pool: Arc<MySqlPool>, connection: &mut TcpStream, address: SocketAddr) {
    let login_start = Utc::now();

    let (mut tx, mut rx) = mpsc::channel::<Vec<u8>>(128);

    let _ = connection.set_nodelay(true);
    
    let mut username = String::new();
    let mut password = String::new();
    let mut client_info = String::new();

    let mut line_reader = BufReader::new(connection);

    //Read everything
    let username_err = line_reader.read_line(&mut username).await;
    let password_err = line_reader.read_line(&mut password).await;
    let client_info_err = line_reader.read_line(&mut client_info).await;

    //Recover connection, as we moved `connection` into BufReader
    let recovered_conn = line_reader.into_inner();

    if username_err.is_err() || password_err.is_err() || client_info_err.is_err() {
        send_wrong_version(recovered_conn, &mut tx, &mut rx).await;
        return
    }
    
    let client_info_split: Vec<&str> = client_info.split('|').collect();

    if client_info_split.len() != 4 {
        send_wrong_version(recovered_conn, &mut tx, &mut rx).await;
        return;
    }

    let security_parts_split: Vec<&str> = client_info_split[3].split(':').collect();

    if security_parts_split.len() != 2 {
        send_wrong_version(recovered_conn, &mut tx, &mut rx).await;
        return;
    }

    let client_hash = security_parts_split[0];
    let mac_address = security_parts_split[1];

    //Parse version as int, so it's easier to compare
    let version_err = client_info_split[0].trim_start_matches('b').parse::<i32>();

    if version_err.is_err() {
        send_wrong_version(recovered_conn, &mut tx, &mut rx).await;
        return;
    }

    let parsed_version = version_err.unwrap();

    if parsed_version < 1816 {
        send_wrong_version(recovered_conn, &mut tx, &mut rx).await;
        return;
    }

    let timezone_err = client_info_split[1].parse::<i32>();

    if timezone_err.is_err() {
        send_wrong_version(recovered_conn, &mut tx, &mut rx).await;
        return;
    }

    let allow_city_parsed = {
        if client_info_split[2] == "1" {
            true
        } else {
            false
        }
    };

    let client_info = ClientInformation {
        version: parsed_version,
        timezone: timezone_err.unwrap(),
        client_hash: client_hash.to_string(),
        mac_address_hash: mac_address.to_string(),
        allow_city: allow_city_parsed
    };

    let user = db::User::from_username(pool, username).await;


}