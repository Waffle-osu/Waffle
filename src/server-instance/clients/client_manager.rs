use std::sync::Arc;

use dashmap::DashMap;
use lazy_static::lazy_static;

use super::waffle_client::WaffleClient;

struct ClientManager {
    clients_by_id: DashMap<u64, Arc<dyn WaffleClient + Sync + Send>>,
    clients_by_name: DashMap<String, Arc<dyn WaffleClient + Sync + Send>>
}

lazy_static! {
    static ref manager: ClientManager = ClientManager { clients_by_id: DashMap::new(), clients_by_name: DashMap::new() };
}

fn register_client(client: Arc<dyn WaffleClient + Sync + Send>) {
    let user = client.get_user();

    manager.clients_by_id.insert(user.user_id, client.clone());
    manager.clients_by_name.insert(user.username, client.clone());
}
