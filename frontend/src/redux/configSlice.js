import { createSlice } from '@reduxjs/toolkit'

const initialState = {}

export const configSlice = createSlice({
  name: 'config',
  initialState,
  reducers: {
    setConfig: (state, action) => {
      console.log("Config loaded:", action.payload)
      return action.payload;
    },
  },
})


export const { setConfig } = configSlice.actions
export default configSlice.reducer
