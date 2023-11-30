use std::sync::Arc;

use sqlx::{MySqlPool, mysql::MySqlPoolOptions};
use tokio::net::TcpListener;
use osu_listener::bancho_listener;
use irc_listener::irc_listener;

mod osu_listener;
mod irc_listener;
mod clients;
mod chat;
mod bot;
mod lobby;
mod osu;
mod irc;

#[tokio::main]
async fn main() {
    let db_read_result = 
        std::fs::read_to_string("database.txt")
            .expect("Failed to read connection data to database.");

    let split: Vec<&str> = db_read_result.split("\n").collect();

    let server = split[0];
    let username = split[1];
    let password = split[2];
    let database = split[3];

    let database_url = format!("mysql://{}:{}@{}/{}", username, password, server, database);

    let pool: MySqlPool = MySqlPoolOptions::new()
        .max_connections(8)
        .connect(&database_url)
        .await
        .expect("Failed to connect to the database");

    let arc_pool = Arc::new(pool);

    let bancho_pool = arc_pool.clone();
    let irc_pool = arc_pool.clone();

    tokio::spawn(async move {
        bancho_listener(bancho_pool).await
    });

    irc_listener(irc_pool).await
}
