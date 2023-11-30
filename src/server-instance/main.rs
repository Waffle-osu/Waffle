use tokio::net::TcpListener;
use osu_listener::bancho_listener;
use irc_listener::irc_listener;

mod osu_listener;
mod irc_listener;
mod clients;


#[tokio::main]
async fn main() {
    tokio::spawn(async move {
        bancho_listener().await
    });

    tokio::spawn(async move {
        irc_listener().await
    });
}
