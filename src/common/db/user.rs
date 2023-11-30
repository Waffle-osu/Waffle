#[derive(sqlx::FromRow)]
pub struct User {
    user_id: u64,
    username: String,
    password: String,
    country: String,
    silenced_until: u64,
    banned: bool,
}

impl User {
    fn from_id(id: u64) -> User {
        
    }

    fn from_username(username: String) -> User {

    }

    fn from_where_condition(sql_where: String) -> User {
        let mut 
    }
}