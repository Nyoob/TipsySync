import { useState, useEffect } from 'react';
import { Box, Stack, Tooltip, Typography } from '@mui/material';
import AdminPanelSettingsIcon from "@mui/icons-material/AdminPanelSettings"
import BroadcastOnPersonalIcon from '@mui/icons-material/BroadcastOnPersonal';
import FavoriteIcon from "@mui/icons-material/Favorite"
import VerifiedUserIcon from "@mui/icons-material/VerifiedUser"
import ShieldIcon from "@mui/icons-material/Shield"
import StarsIcon from "@mui/icons-material/Stars"
import { imageLookup } from '../helper';
import { useSelector } from 'react-redux';

export default function Chat({ isWidget }) {
  const chatMsgs = useSelector(s => s.chatMsgs);

  return (
    <div>
      <h1>Chat</h1>
      <Box sx={{ display: "flex", gap: 4 }}>
        <Stack spacing={2} sx={{ flexGrow: 4 }}>
          {chatMsgs.map((chatMsg) => {
            return <ChatMessage chatMsg={chatMsg} />
          })}
        </Stack>
      </Box>
    </div>
  )
}

function ChatMessage({ chatMsg }) {
  const { Event } = chatMsg
  const { Timestamp, ChatMessage, User } = Event
  const { Username, SubscribedTierColor, IsBroadcaster, IsMod } = User
  const Provider = chatMsg.Provider

  return (
    <Box
      sx={{
        display: "flex",
        alignItems: "flex-start",
        gap: 2,
        padding: 1,
        borderRadius: 1,
        backgroundColor: "background.paper",
      }}
    >
      <Box
        sx={{
          width: 48,
          height: 48,
          flexShrink: 0,
          borderRadius: "50%",
          display: "flex",
          alignItems: "center",
          justifyContent: "left",
          fontWeight: "bold",
          fontSize: 18,
          userSelect: "none",
        }}
      >
        <img src={imageLookup[Provider]} style={{ width: "100%", height: "100%" }} />
      </Box>

      {/* Message content */}
      <Stack spacing={0.5} sx={{ flexGrow: 1 }}>
        <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
          <Typography
            sx={{
              color: IsBroadcaster ? "gold" : IsMod ? "red" : SubscribedTierColor || "text.primary",
              fontWeight: "bold",
            }}
            component="span"
          >
            {Username}
          </Typography>
          {IsBroadcaster && (
            <Tooltip title="Broadcaster / Streamer">
              <BroadcastOnPersonalIcon sx={{ color: "gold", fontSize: 18 }} />
            </Tooltip>
          )}
          {IsMod && (
            <Tooltip title='This user is a moderator or admin'>
              <AdminPanelSettingsIcon sx={{ color: "red", fontSize: 18 }} />
            </Tooltip>
          )}
          {User.IsSubscribed && (
            <Tooltip title={`Subscribed (Tier:${User.SubscribedTierName})`}>
              <FavoriteIcon sx={{ color: "#FF4081", fontSize: 18 }} />
            </Tooltip>
          )}
          {User.StripchatIsKing && (
            <Tooltip title="Stripchat King">
              <VerifiedUserIcon sx={{ color: "#FFD700", fontSize: 18 }} />
            </Tooltip>
          )}
          {User.StripchatIsKnight && (
            <Tooltip title="Stripchat Knight">
              <ShieldIcon sx={{ color: "#C0C0C0", fontSize: 18 }} />
            </Tooltip>
          )}
          {User.StripchatIsUltimate && (
            <Tooltip title="Stripchat Ultimate">
              <StarsIcon sx={{ color: "#00BFFF", fontSize: 18 }} />
            </Tooltip>
          )}
          <Typography
            variant="caption"
            color="text.secondary"
            sx={{ marginLeft: "auto" }}
          >
            {new Date(Timestamp).toLocaleTimeString()}
          </Typography>
        </Box>
        <Typography sx={{ textAlign: "left" }}>{ChatMessage}</Typography>
      </Stack>
    </Box>
  )
}
