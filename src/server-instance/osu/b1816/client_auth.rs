use std::{sync::Arc, net::SocketAddr, ops::Sub};

use chrono::Utc;
use common::{packets::{derived::{BanchoLoginReply, BanchoAnnounce, BanchoFriendsList, BanchoProtocolNegotiation, BanchoLoginPermissions, BanchoUserPresence, BanchoStatsUpdate}, BanchoPacket, BanchoPresence, BanchoUserStats, BanchoStatusUpdate}, db};
use dashmap::DashMap;
use sqlx::MySqlPool;
use tokio::{net::TcpStream, sync::{mpsc::{self, Receiver, Sender}, Mutex}, io::{BufReader, AsyncBufReadExt}};

use crate::{clients::{ClientManager, waffle_client::WaffleClient}, osu};

use super::client::{ClientInformation, OsuClient2011};

async fn send_and_close(connection: TcpStream, queue: &mut Receiver<BanchoPacket>) {
    while let Some(packet) = queue.recv().await {
        let buffer = packet.send();
        let slice = buffer.as_slice();

        connection.try_write(slice).expect("Failed to write packets!");
    }
}

async fn send_wrong_version(connection: TcpStream, queue_send: &mut Sender<BanchoPacket>, queue_receive: &mut Receiver<BanchoPacket>) {
    BanchoLoginReply::send_wrong_version(&queue_send).await;
        
    send_and_close(connection, queue_receive).await;
}

pub async fn handle_new_client(pool: Arc<MySqlPool>, connection: TcpStream, address: SocketAddr) {
    let login_start = Utc::now();
    
    let (mut tx, mut rx) = mpsc::channel::<BanchoPacket>(128);
    
    let _ = connection.set_nodelay(true);
    let _ = connection.readable().await;
    
    let mut username = String::new();
    let mut password = String::new();
    let mut client_info = String::new();

    let mut line_reader = BufReader::new(connection);

    //Read everything
    let username_err = line_reader.read_line(&mut username).await;
    let password_err = line_reader.read_line(&mut password).await;
    let client_info_err = line_reader.read_line(&mut client_info).await;

    let truncate_rn = |str: &mut String| {
        if str.ends_with('\n') {
            str.pop();

            if str.ends_with('\r') {
                str.pop();
            }
        }
    };

    //Remove \r and \n
    truncate_rn(&mut username);
    truncate_rn(&mut password);
    truncate_rn(&mut client_info);

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

    if security_parts_split.len() != 3 {
        send_wrong_version(recovered_conn, &mut tx, &mut rx).await;
        return;
    }

    let client_hash = security_parts_split[0];
    let mac_address = security_parts_split[1];

    let comparable_version_string = 
        client_info_split[0]
            .trim_start_matches('b')
            .trim_end_matches(".peppy")
            .trim_end_matches(".test")
            .trim_end_matches(".ctbtest")
            .trim_end_matches(".arcade");

    //Parse version as int, so it's easier to compare
    let version_parse = comparable_version_string.parse::<i32>();

    if version_parse.is_err() {
        send_wrong_version(recovered_conn, &mut tx, &mut rx).await;
        return;
    }

    let parsed_version = version_parse.unwrap();

    //Older than b1816 not supprted over regular bancho
    if parsed_version < 1815 {
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

    let verify_begin = Utc::now();

    let password_valid = bcrypt::verify(password, user.password.as_str());

    let verify_end = Utc::now();

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

    let presence = BanchoPresence {
        user_id: user.user_id as i32,
        avatar_extension: 1,
        username: Some(user.username.clone()),
        timezone: client_info.timezone as u8,
        country: 0,
        city: None,
        permissions: 0,
        longitude: 0.0f32,
        latitude: 0.0f32,
        rank: 1
    };

    let osu_stats = BanchoUserStats {
        user_id: user.user_id as i32,
        status: BanchoStatusUpdate { 
            status: 0, 
            status_text: Some(String::from("Waffle!!!!")), 
            beatmap_checksum: Some(String::from("kek")), 
            current_mods: 0, 
            play_mode: 0, 
            beatmap_id: 0 
        },
        ranked_score: 0,
        accuracy: 0.0f32,
        playcount: 0,
        total_score: 0,
        rank: 1
    };
    
    BanchoFriendsList::send(&mut tx, to_i32_list).await;
    BanchoAnnounce::send(&mut tx, String::from("Welcome to Waffle!")).await;
    BanchoProtocolNegotiation::send(&mut tx, 7).await;
    BanchoLoginPermissions::send(&mut tx, user.privileges).await;
    BanchoUserPresence::send(&mut tx, &presence).await;
    BanchoStatsUpdate::send(&mut tx, &osu_stats).await;

    //TODO: protocol negotiation
    //TODO: permissions
    //TODO: user presence
    //TODO: osu update
    //TODO: other users presence
    //TODO: other users osu update
    //TODO: channel available
    //TODO: channel autojoin

    let client = OsuClient2011 {
        connection: recovered_conn,
        continue_running: true,
        logon_time: Utc::now(),
        last_receive: Utc::now(),
        last_ping: Utc::now(),
        away_message: String::from(""),
        spectators: DashMap::new(),
        spectating_client: None,
        packet_queue_send: Mutex::new(tx),
        packet_queue_recv: Mutex::new(rx),

        user: user,
        osu_stats: osu_stats_query.unwrap(),
        taiko_stats: taiko_stats_query.unwrap(),
        catch_stats: catch_stats_query.unwrap(),
        mania_stats: mania_stats_query.unwrap(),
    };

    let as_arc = Arc::new(client);

    ClientManager::register_client(
        Arc::new(WaffleClient::Osu(as_arc.clone()))
    );

    tokio::spawn(OsuClient2011::maintain_client(as_arc.clone()));
    tokio::spawn(OsuClient2011::handle_incoming(as_arc.clone()));
    tokio::spawn(OsuClient2011::send_outgoing(as_arc.clone()));

    let time = Utc::now().sub(login_start).num_milliseconds();

    println!("Login for {} took {}ms; bcrypt took {}ms", as_arc.clone().user.username, time, verify_end.sub(verify_begin).num_milliseconds());
}