use tokio::net::TcpListener;

pub async fn irc_listener(db_conn: Arc<MySqlPool>) {
    let listener = TcpListener::bind("127.0.0.1:6667")
        .await
        .expect("Failed to create IRC Server...");

    loop {
        let (mut socket, addr) = 
            listener
                .accept()
                .await
                .expect("Failed to accept socket!");

        
    }
}