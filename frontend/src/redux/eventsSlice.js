import { createSlice } from '@reduxjs/toolkit'

const initialState = []

export const eventsSlice = createSlice({
  name: 'events',
  initialState,
  reducers: {
    setEvents: (state, action) => {
      var { maxEvents, data } = action.payload;
      var newE = [data, ...state];
      if (newE.length > maxEvents) {
        newE.length = maxEvents;
      }
      return newE;
    },
  },
})


export const { setEvents } = eventsSlice.actions
export default eventsSlice.reducer
