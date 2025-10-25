import { Grow, Paper, Stack } from '@mui/material';
import { useEffect, useState } from 'react';

import { EventsOn } from '../../wailsjs/runtime/runtime';
import { eventTypes } from '../helper';

function Overview() {
  const [events, setEvents] = useState([]);

  useEffect(() => {
    eventTypes.forEach(type => {
      EventsOn(type, (data) => {
        SetEvents(e => [data, ...e]);
      })
    })
  }, [])

  return (
    <div>
      <h1>Events</h1>
      <Stack spacing={2}>
        {events.map((event) => {
          <Paper elevation={2}>

          </Paper>
        })}
      </Stack>
    </div>
  )
}

export default Overview

