import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"

export interface AuthSliceState {
    token: string
}

const initialState: AuthSliceState = {
    token: ""
}

export const attemptLogin = createAsyncThunk("auth/login", async () => {

})

export const authSlice = createSlice({
    name: "authSlice",
    initialState,
    reducers: {
        getUser: (state) => {

        }
    },
    extraReducers: (builder) => {
        builder.addCase(attemptLogin.pending, (state) => {

        })
    }
})

export const { getUser } = authSlice.actions;

export const selectToken = (state: AuthSliceState) => state.token;

export default authSlice.reducer;