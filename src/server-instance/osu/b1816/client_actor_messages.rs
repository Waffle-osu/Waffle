use actix::{Handler, Message, dev::MessageResponse};

use super::client::OsuClient2011;
struct SetContinueRunning(bool);

impl Message for SetContinueRunning {
    type Result = What;
}

struct What(bool);

impl MessageResponse<OsuClient2011, SetContinueRunning> for What {
    fn handle(self, ctx: &mut <OsuClient2011 as actix::prelude::Actor>::Context, tx: Option<actix::prelude::dev::OneshotSender<<SetContinueRunning as Message>::Result>>) {
        ctx.
    }
}

impl Handler<SetContinueRunning> for OsuClient2011 {
    fn handle(&mut self, msg: SetContinueRunning, ctx: &mut Self::Context) -> Self::Result {
        What(self.continue_running)
    }

    type Result = What;
}