use tokio::{io::AsyncReadExt, net::TcpStream};

pub async fn handle_socket(mut socket: TcpStream) {
    let mut buf = Vec::with_capacity(4096);

    loop {
        let read_result = socket.read_buf(&mut buf).await;

        let read = match read_result {
            Err(e) => {
                eprintln!("Read from client failed; Error: {:?}", e);
                return;
            }

            Ok(n) => n,
        };

        //Apperantly that means disconnect in tokio
        if read == 0 {
            return;
        }
    }
}
