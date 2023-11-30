use sqlx::{MySqlPool, mysql::MySqlPoolOptions, pool::PoolOptions, Pool, MySql};

mod user;

pub static mut db_pool: Option<Pool<MySql>> = None;

async fn initialize_db() {
    let db_read_result = 
        std::fs::read_to_string("database.txt")
            .expect("Failed to read connection data to database.");
    
    let split: Vec<&str> = db_read_result.split("\n").collect();

    let server = split[0];
    let username = split[1];
    let password = split[2];
    let database = split[3];

    let database_url = format!("mysql://{}:{}@{}/{}", username, password, server, database);

    let pool = MySqlPoolOptions::new()
        .max_connections(8)
        .connect(&database_url)
        .await
        .expect("Failed to connect to database!");

    unsafe { db_pool = Some(pool) };
}

