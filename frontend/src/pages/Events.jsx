import { useState, useEffect } from 'react';
import { Box, Grow, IconButton, Paper, Stack, Typography } from '@mui/material';

import { capitalizeFirstLetter, currencyName, eventTypes, genderLookup, getEventItemText, getSubscriptionName, imageLookup } from '../helper';
import CloseIcon from '@mui/icons-material/Close';

function Events({ events }) {
  const [selectedEvent, setSelectedEvent] = useState(null);
  const [highlightedEventIds, setHighlightedEventIds] = useState([]);

  useEffect(() => { // highlight new events
    console.log("events", events)
    if (events.length == 0) { return; }
    const newEventId = events[0].Event.Id;
    setHighlightedEventIds(p => [...p, newEventId]);
    setTimeout(() => {
      setHighlightedEventIds(prev => prev.filter(id => id !== newEventId));
    }, 10000);
  }, [events])

  return (
    <div>
      <h1>Events</h1>
      <Box sx={{ display: "flex", gap: 4 }}>
        <Stack spacing={2} sx={{ flexGrow: 4 }}>
          {events.map((event) => {
            return <EventItem event={event} setSelectedEvent={setSelectedEvent} isHighlighted={highlightedEventIds.includes(event.Event.Id)} />
          })}
        </Stack>
        {selectedEvent && <EventDetails event={selectedEvent} closeDetails={() => setSelectedEvent(null)} />}
      </Box>
    </div>
  )
}

function EventDetails({ event, closeDetails }) {
  console.log(event)
  const texts = getEventItemText(event);
  const subNames = getSubscriptionName(event.Provider);
  return <Paper elevation={2} sx={{ padding: 4, flexBasis: 100, minWidth: 300, position: 'relative' }}>
    <IconButton
      size="large"
      color="inherit"
      aria-label="menu"
      sx={{ position: 'absolute', right: 8, top: 8 }}
    >
      <CloseIcon />
    </IconButton>

    <Typography variant="h4" sx={{ fontWeight: 'bold', color: event.Event.User.SubscribedTierColor != "" ? event.Event.User.SubscribedTierColor : "unset" }}>
      {event.Event.User.Username}
    </Typography>
    <Typography variant="subtitle1" color="text.secondary" sx={{ marginTop: 0.5 }}>
      {texts.longText}
    </Typography>

    {event.Type == "tip" && <Box sx={{ marginTop: 4 }}>
      <Typography variant="h6">Tip</Typography>
      <Typography variant="body1"><b>Message:</b> {event.Event.TipMessage}</Typography>
      <Typography variant="body1"><b>Tipped:</b> {event.Event.TipValue + currencyName?.[event?.Provider]}</Typography>
      <Typography variant="body1"><b>Dollar Value:</b> {event.Event.TipValueInDollars}</Typography>
    </Box>}

    {event.Type == "subscribe" && <Box sx={{ marginTop: 4 }}>
      <Typography variant="h6">New Subscription</Typography>
      <Typography variant="body1"><b>Subscription Tier ID:</b> {event.Event.TierId}</Typography>
      <Typography variant="body1"><b>Subscription Tier:</b> {event.Event.TierName}</Typography>
      <Typography variant="body1"><b>Subscription Streak:</b> {event.Event.User.Streak}</Typography>
    </Box>}

    <Box sx={{ marginTop: 4 }}>
      <Typography variant="h6">User</Typography>
      <Typography variant="body1"><b>Username:</b> {event.Event.User.Username}</Typography>
      <Typography variant="body1"><b>Gender:</b> {genderLookup[event.Event.User.Gender]}</Typography>
      <Typography variant="body1"><b>Has Tokens:</b> {event.Event.User.HasTks}</Typography>
      <Typography variant="body1"><b>Subscribed:</b> {event.Event.User.Subscribed}</Typography>
      <Typography variant="body1"><b>Subscription Tier:</b> {event.Event.User.SubscribedTiername}</Typography>
    </Box>

    <Box sx={{ marginTop: 4 }}>
      <Typography variant="h6">Additional</Typography>
      <Typography variant="body1"><b>Event-ID:</b> {event.Event.Id}</Typography>
      <Typography variant="body1"><b>Timestamp:</b> {event.Event.Timestamp}</Typography>
    </Box>

  </Paper>
}

function EventItem({ event, setSelectedEvent, isHighlighted }) {
  if (!event) return null;

  const texts = getEventItemText(event);

  const highlightStyle = isHighlighted
    ? {
      background: "radial-gradient(ellipse at 70% 30%, #263557 0%, #151a2d 90%)",
      boxShadow: "0 0 16px 2px #3f73da55, 0 0 28px 6px #6379cb33",
      border: "1px solid",
      borderImage: "linear-gradient(90deg, #26e6ff 10%, #405fff 50%, #9740ff 90%) 1",
      filter: "brightness(1.18)",
      color: "#fff"
    }
    : {};

  return <Paper elevation={2}
    sx={{ display: 'flex', padding: 2, alignItems: "center", transition: "background 1.5s, box-shadow 1.4s, border 1.3s", ...highlightStyle }}
    onClick={() => setSelectedEvent(event)}>
    <Box
      component="img"
      src={imageLookup?.[event?.Provider] ?? imageLookup["chaturbate"]}
      alt="Platform"
      sx={{
        width: 50,
        marginRight: 2,
        flexShrink: 0
      }}
    />

    <Box sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column', justifyContent: 'center', textAlign: "left" }}>
      <Typography variant="h5" component="div" sx={{ fontWeight: 'bold' }}>
        {event.Event.User.Username}
      </Typography>
      <Typography variant="body2" color="text.secondary" sx={{ marginTop: 0.5 }}>
        {event.Event.TipMessage}
      </Typography>
    </Box>

    <Box sx={{ textAlign: 'right', minWidth: 80, alignSelf: "end" }}>
      {event.Type == "tip" && <Typography variant="h6" component="div" sx={{ fontWeight: 'bold' }}>
        ${event.Event.TipValueInDollars}
      </Typography>}
      <Typography variant="body2" color="text.secondary" sx={{ marginTop: 0.5 }}>
        {texts.shortText}
      </Typography>
    </Box>
  </Paper>
}

export default Events

