use std::sync::Arc;

use dashmap::DashMap;
use lazy_static::lazy_static;
use tokio::sync::Mutex;

use crate::osu;

use super::waffle_client::WaffleClient;

pub struct ClientManager {
    clients_by_id: DashMap<u64, Arc<dyn WaffleClient + Send + Sync>>,
    clients_by_name: DashMap<String, Arc<dyn WaffleClient + Send + Sync>>
}

lazy_static! {
    static ref manager: ClientManager = ClientManager { clients_by_id: DashMap::new(), clients_by_name: DashMap::new() };
}

impl ClientManager {
    pub fn register_client(client: Arc<dyn WaffleClient + Send + Sync>) {
        let user = client.get_user();

        manager.clients_by_id.insert(user.user_id, client);
        manager.clients_by_name.insert(user.username, client);
    }   

    pub fn get_client_by_id(user_id: u64) -> Option<Arc<dyn WaffleClient + Send + Sync>> {
        let what = manager.clients_by_id.get(&user_id);

        match what {
            None => return None,
            Some(x) => return Some(x.value().clone())
        }
    } 

    pub fn get_client_by_name(username: String) -> Option<Arc<dyn WaffleClient + Send + Sync>> {
        let what = manager.clients_by_name.get(&username);

        match what {
            None => return None,
            Some(x) => return Some(x.value().clone())
        }
    } 
}