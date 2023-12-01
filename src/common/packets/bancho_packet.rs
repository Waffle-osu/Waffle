use std::sync::mpsc::Sender;

use binary_rw::{BinaryReader, BinaryWriter};

#[derive(Debug, PartialEq, Eq, PartialOrd, Ord)]
pub enum InternalRequestType {
    GetServerStatus,
    AnnounceUserJoin,
    AnnounceUserLeft,
    RetrievePresencePacket,
    SendChatMessage,
}

#[derive(Debug, Eq, PartialEq, Copy, Clone)]
#[repr(u16)]
pub enum BanchoRequestType {
    OsuSendUserStatus = 0,
    OsuSendIrcMessage = 1,
    OsuExit = 2,
    OsuRequestStatusUpdate = 3,
    OsuPong = 4,
    BanchoLoginReply = 5,
    BanchoSendMessage = 7,
    BanchoPing = 8,
    BanchoHandleIrcChangeUsername = 9,
    BanchoHandleIrcQuit = 10,
    BanchoHandleOsuUpdate = 11,
    BanchoHandleOsuQuit = 12,
    BanchoSpectatorJoined = 13,
    BanchoSpectatorLeft = 14,
    BanchoSpectateFrames = 15,
    OsuStartSpectating = 16,
    OsuStopSpectating = 17,
    OsuSpectateFrames = 18,
    OsuErrorReport = 20,
    OsuCantSpectate = 21,
    BanchoSpectatorCantSpectate = 22,
    BanchoGetAttention = 23, //TODO: maybe once there's an admin panel or something? or maybe as a chat command for admins
    BanchoAnnounce = 24,
    OsuSendIrcMessagePrivate = 25,
    BanchoMatchUpdate = 26,
    BanchoMatchNew = 27,
    BanchoMatchDisband = 28,
    OsuLobbyPart = 29,
    OsuLobbyJoin = 30,
    OsuMatchCreate = 31,
    OsuMatchJoin = 32,
    OsuMatchPart = 33,
    BanchoLobbyJoin = 34,
    BanchoLobbyPart = 35,
    BanchoMatchJoinSuccess = 36,
    BanchoMatchJoinFail = 37,
    OsuMatchChangeSlot = 38,
    OsuMatchReady = 39,
    OsuMatchLock = 40,
    OsuMatchChangeSettings = 41,
    BanchoFellowSpectatorJoined = 42,
    BanchoFellowSpectatorLeft = 43,
    OsuMatchStart = 44,
    BanchoMatchStart = 46,
    OsuMatchScoreUpdate = 47,
    BanchoMatchScoreUpdate = 48,
    OsuMatchComplete = 49,
    BanchoMatchTransferHost = 50,
    OsuMatchChangeMods = 51,
    OsuMatchLoadComplete = 52,
    BanchoMatchAllPlayersLoaded = 53,
    OsuMatchNoBeatmap = 54,
    OsuMatchNotReady = 55,
    OsuMatchFailed = 56,
    BanchoMatchPlayerFailed = 57,
    BanchoMatchComplete = 58,
    OsuMatchHasBeatmap = 59,
    OsuMatchSkipRequest = 60,
    BanchoMatchSkip = 61,
    BanchoUnauthorized = 62, //Unused
    OsuChannelJoin = 63,
    BanchoChannelJoinSuccess = 64,
    BanchoChannelAvailable = 65,
    BanchoChannelRevoked = 66,
    BanchoChannelAvailableAutojoin = 67,
    OsuBeatmapInfoRequest = 68,
    BanchoBeatmapInfoReply = 69,
    OsuMatchTransferHost = 70,
    BanchoLoginPermissions = 71,
    BanchoFriendsList = 72,
    OsuFriendsAdd = 73,
    OsuFriendsRemove = 74,
    BanchoProtocolNegotiation = 75,
    BanchoTitleUpdate = 76, //TODO: once site's a thing this could be used
    OsuMatchChangeTeam = 77,
    OsuChannelLeave = 78,
    OsuReceiveUpdates = 79, //Unused
    BanchoMonitor = 80,
    BanchoMatchPlayerSkipped = 81,
    OsuSetIrcAwayMessage = 82,
    BanchoUserPresence = 83,
    OsuUserStatsRequest = 85,
    BanchoRestart = 86,
}

impl From<u16> for BanchoRequestType {
    fn from(value: u16) -> Self {
        unsafe { std::mem::transmute(value) }
    }
}

#[repr(C)]
#[derive(Debug, PartialEq, Eq)]
pub struct BanchoPacketHeader {
    pub packet_id: BanchoRequestType,
    pub compressed: bool,
    pub size: i32,
}

impl BanchoRequestType {
    fn to_primitive(&self) -> u16 {
        *self as u16
    }
}

pub trait BanchoSerializable {
    fn read(&mut self, reader: &mut BinaryReader);
    fn write(&self, writer: &mut BinaryWriter);
}

impl BanchoSerializable for BanchoPacketHeader {
    fn read(&mut self, reader: &mut BinaryReader) {
        self.packet_id = BanchoRequestType::from(reader.read_u16().unwrap());
        self.compressed = reader.read_bool().unwrap();
        self.size = reader.read_i32().unwrap();
    }

    fn write(&self, writer: &mut BinaryWriter) {
        let msg = "Failed to write packet header!";

        let packet_id = self.packet_id.to_primitive();

        writer.write_u16(packet_id).expect(msg);
        writer.write_bool(self.compressed).expect(msg);
        writer.write_i32(self.size).expect(msg);
    }
}

pub struct BanchoPacket {
    header: BanchoPacketHeader,
    data: Vec<u8>
}

impl BanchoPacket {
    fn from(packet_id: BanchoRequestType, serializable: &dyn BanchoSerializable) {

    }
}
