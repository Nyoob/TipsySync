import { Grow, Paper, Stack } from '@mui/material';
import { useEffect, useState } from 'react';

function Overview() {
  const [events, setEvents] = useState([]);

  useEffect(() => {
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

