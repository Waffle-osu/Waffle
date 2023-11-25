use tokio::{io::AsyncReadExt, net::TcpListener};

use crate::server_handling::handle_socket;

mod server_handling;

#[tokio::main]
async fn main() {
    let server_file_read_result = std::fs::read_to_string("servers.txt");

    let cwd = std::env::current_dir().unwrap();

    println!("Running from: {}", cwd.to_str().unwrap());

    let servers: Vec<String> = match server_file_read_result {
        Err(e) => {
            eprintln!("Failed to read in server instance file! Server Statuses may be reported incorrectly; Error: {:?}", e);

            return;
        }
        Ok(str) => str.lines().map(String::from).collect(),
    };

    for server in servers {
        println!("{}", server);
    }

    let listener = TcpListener::bind("127.0.0.1:7419")
        .await
        .expect("Failed to create Head Server...");

    loop {
        let (mut socket, _) = listener
            .accept()
            .await
            .expect("Accepting TCP Connection failed");

        tokio::spawn(handle_socket(socket));
    }
}
