import { configureStore } from '@reduxjs/toolkit'

import configReducer from './configSlice'
import infoReducer from './infoSlice'
import chatMsgsReducer from './chatMsgsSlice'
import eventsReducer from './eventsSlice'

export const store = configureStore({
  reducer: {
    config: configReducer,
    info: infoReducer,
    chatMsgs: chatMsgsReducer,
    events: eventsReducer,
  },
})
