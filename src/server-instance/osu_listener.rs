use std::sync::Arc;

use sqlx::MySqlPool;
use tokio::net::TcpListener;

use crate::osu;

pub async fn bancho_listener(db_conn: Arc<MySqlPool>) {
    let listener = TcpListener::bind("127.0.0.1:13381")
        .await
        .expect("Failed to create Bancho Server...");

    loop {
        let (socket, addr) = 
            listener
                .accept()
                .await
                .expect("Failed to accept socket!");

        tokio::spawn(osu::handle_new_client(db_conn.clone(), socket, addr));
    }
}