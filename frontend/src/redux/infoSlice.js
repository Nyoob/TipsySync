import { createSlice } from '@reduxjs/toolkit'

const initialState = {}

export const infoSlice = createSlice({
  name: 'info',
  initialState,
  reducers: {
    setInfo: (state, action) => {
      console.log("Info loaded:", action.payload)
      return action.payload;
    },
  },
})


export const { setInfo } = infoSlice.actions
export default infoSlice.reducer
