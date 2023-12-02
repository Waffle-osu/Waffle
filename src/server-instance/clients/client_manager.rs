use std::sync::Arc;

use dashmap::DashMap;
use lazy_static::lazy_static;
use tokio::sync::Mutex;

use crate::osu;

use super::waffle_client::WaffleClient;

pub struct ClientManager {
    clients_by_id: DashMap<u64, Arc<WaffleClient>>,
    clients_by_name: DashMap<String, Arc<WaffleClient>>
}

lazy_static! {
    static ref manager: ClientManager = ClientManager { clients_by_id: DashMap::new(), clients_by_name: DashMap::new() };
}

impl ClientManager {
    pub fn register_client(client: Arc<WaffleClient>) {
        match *client {
            WaffleClient::Irc(ref irc_client) => {
                manager.clients_by_id.insert(irc_client.user.user_id, client.clone());
                manager.clients_by_name.insert(irc_client.user.username.clone(), client.clone());
            },
            WaffleClient::Osu(ref osu_client) => {
                let user = osu_client.get_user();
    
                manager.clients_by_id.insert(user.user_id, client.clone());
                manager.clients_by_name.insert(user.username, client.clone());
            }
        }
    }   

    pub fn get_client_by_id(user_id: u64) -> Option<Arc<WaffleClient>> {
        let what = manager.clients_by_id.get(&user_id);

        match what {
            None => return None,
            Some(x) => return Some(x.value().clone())
        }
    } 

    pub fn get_client_by_name(username: String) -> Option<Arc<WaffleClient>> {
        let what = manager.clients_by_name.get(&username);

        match what {
            None => return None,
            Some(x) => return Some(x.value().clone())
        }
    } 
}