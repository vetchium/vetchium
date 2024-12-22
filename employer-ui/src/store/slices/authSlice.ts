import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import axios from 'axios';
import { AuthState, SignInCredentials, SignInResponse, TFARequest, TFAResponse } from '@/types/auth';

const initialState: AuthState = {
  isAuthenticated: false,
  token: null,
  user: null,
  loading: false,
  error: null,
};

export const signIn = createAsyncThunk(
  'auth/signIn',
  async (credentials: SignInCredentials) => {
    const response = await axios.post<SignInResponse>('/api/employer/signin', credentials);
    return response.data;
  }
);

export const setTFAToken = createAsyncThunk(
  'auth/setTFAToken',
  async (tfaRequest: TFARequest) => {
    const response = await axios.post<TFAResponse>('/api/employer/tfa', tfaRequest);
    return response.data;
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    signOut: (state) => {
      state.isAuthenticated = false;
      state.token = null;
      state.user = null;
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(signIn.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(signIn.fulfilled, (state, action) => {
        state.loading = false;
        state.token = action.payload.token;
        state.isAuthenticated = true;
      })
      .addCase(signIn.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Sign in failed';
      })
      .addCase(setTFAToken.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(setTFAToken.fulfilled, (state, action) => {
        state.loading = false;
        state.token = action.payload.session_token;
        state.isAuthenticated = true;
      })
      .addCase(setTFAToken.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'TFA verification failed';
      });
  },
});

export const { signOut } = authSlice.actions;
export default authSlice.reducer; 