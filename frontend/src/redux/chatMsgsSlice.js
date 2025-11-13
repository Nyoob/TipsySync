import { createSlice } from '@reduxjs/toolkit'

const initialState = []

export const chatMsgsSlice = createSlice({
  name: 'chatMsgs',
  initialState,
  reducers: {
    setChatMsgs: (state, action) => {
      var { maxMsgs, data } = action.payload;
      var newE = [data, ...state];
      if (newE.length > maxMsgs) {
        newE.length = maxMsgs;
      }
      return newE;
    },
  },
})


export const { setChatMsgs } = chatMsgsSlice.actions
export default chatMsgsSlice.reducer
