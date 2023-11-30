use std::sync::Arc;

use sqlx::MySqlPool;
use tokio::net::TcpListener;

pub async fn bancho_listener(db_conn: Arc<MySqlPool>) {
    let listener = TcpListener::bind("127.0.0.1:13381")
        .await
        .expect("Failed to create Bancho Server...");

    loop {
        let (mut socket, addr) = 
            listener
                .accept()
                .await
                .expect("Failed to accept socket!");

        
    }
}