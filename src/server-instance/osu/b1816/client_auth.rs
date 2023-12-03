use std::{sync::Arc, net::SocketAddr, ops::Deref};

use chrono::Utc;
use common::{packets::{derived::{BanchoLoginReply, BanchoAnnounce, BanchoFriendsList}, BanchoPacket}, db};
use dashmap::DashMap;
use sqlx::MySqlPool;
use tokio::{net::TcpStream, sync::mpsc::{self, Receiver, Sender}, io::{BufReader, AsyncBufReadExt, AsyncWriteExt}};

use crate::{clients::{ClientManager, waffle_client::WaffleClient}, osu::OsuClient};

use super::client::{ClientInformation, OsuClient2011};

async fn send_and_close(connection: &mut TcpStream, queue: &mut Receiver<BanchoPacket>) {
    while let Some(packet) = queue.recv().await {
        let buffer = packet.send();
        let slice = buffer.as_slice();

        connection.write(slice).await.expect("Failed to write packets!");
    }

    connection.flush().await.expect("Failed to flush packets");
    connection.shutdown().await.expect("Shutdown of the stream failed!");
}

async fn send_wrong_version(connection: &mut TcpStream, queue_send: &mut Sender<BanchoPacket>, queue_receive: &mut Receiver<BanchoPacket>) {
    BanchoLoginReply::send_wrong_version(&queue_send).await;
        
    send_and_close(connection, queue_receive).await;
}

pub async fn handle_new_client(pool: Arc<MySqlPool>, connection: &mut TcpStream, address: SocketAddr) {
    let login_start = Utc::now();
    let connection_arc = Arc::new(connection);

    let (mut tx, mut rx) = mpsc::channel::<BanchoPacket>(128);

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
    
    //b1816 is supposed to send version, timezone, allow showing city
    //aswell as the MAC Address hash, client hash, etc.
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
    let version_parse = client_info_split[0].trim_start_matches('b').parse::<i32>();

    if version_parse.is_err() {
        send_wrong_version(recovered_conn, &mut tx, &mut rx).await;
        return;
    }

    let parsed_version = version_parse.unwrap();

    //Older than b1816 not supprted over regular bancho
    if parsed_version < 1816 {
        send_wrong_version(recovered_conn, &mut tx, &mut rx).await;
        return;
    }

    let timezone_parse = client_info_split[1].parse::<i32>();

    if timezone_parse.is_err() {
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
        timezone: timezone_parse.unwrap(),
        client_hash: client_hash.to_string(),
        mac_address_hash: mac_address.to_string(),
        allow_city: allow_city_parsed
    };

    let user_query = db::User::from_username(pool.clone(), username).await;

    if user_query.is_none() {
        BanchoLoginReply::send_wrong_login(&mut tx).await;

        send_and_close(recovered_conn, &mut rx).await;
        return;
    }

    let user = user_query.unwrap();

    let password_valid = bcrypt::verify(password, user.password.as_str());

    //I hope this never happens? why would bcrypt fail...
    if password_valid.is_err() {
        BanchoLoginReply::send_server_error(&mut tx).await;

        send_and_close(recovered_conn, &mut rx).await;
        return;
    }

    if user.banned {
        BanchoLoginReply::send_banned(&mut tx).await;

        send_and_close(recovered_conn, &mut rx).await;
        return;
    }

    let existing_client = ClientManager::get_client_by_id(user.user_id);

    if existing_client.is_some() {
        BanchoAnnounce::send(&mut tx, String::from("Another client is already logged in under your name on this server! Disconnecting.")).await;

        send_and_close(recovered_conn, &mut rx).await;
        return;
    }

    BanchoLoginReply::send_success(&mut tx, user.user_id as i32).await;

    let osu_stats_query = db::Stats::from_id(pool.clone(), user.user_id, 0).await;
    let taiko_stats_query = db::Stats::from_id(pool.clone(), user.user_id, 1).await;
    let catch_stats_query = db::Stats::from_id(pool.clone(), user.user_id, 2).await;
    let mania_stats_query = db::Stats::from_id(pool.clone(), user.user_id, 3).await;

    if osu_stats_query.is_none() || taiko_stats_query.is_none() || catch_stats_query.is_none() || mania_stats_query.is_none() {
        BanchoAnnounce::send(&mut tx, String::from("Contact the host of this server. Your user exists in osu_users but your stats don't exist in osu_stats.")).await;

        send_and_close(recovered_conn, &mut rx).await;
        return;
    }

    let friends = db::Friends::get_users_friends(pool.clone(), user.user_id).await;
    let to_i32_list: Vec<i32> = friends.iter().map(|e| e.user_2 as i32).collect();
    
    BanchoFriendsList::send(&mut tx, to_i32_list).await;

    let client = OsuClient2011 {
        continue_running: true,
        logon_time: Utc::now(),
        last_receive: Utc::now(),
        last_ping: Utc::now(),
        away_message: String::from(""),
        spectators: DashMap::new(),
        spectatingClient: None,
        packetQueueSend: Arc::new(tx),
        packetQueueRecv: Arc::new(rx),
    };

    let as_arc = Arc::new(client.to_osu_client());

    ClientManager::register_client(
        Arc::new(WaffleClient::Osu(as_arc))
    );
}