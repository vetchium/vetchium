import { useSelector, useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { RootState } from '@/store';
import { signIn, signUp, signOut } from '@/store/slices/authSlice';
import { SignInCredentials, SignUpRequest } from '@/types/auth';

export const useAuth = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { isAuthenticated, token, user, loading, error } = useSelector(
    (state: RootState) => state.auth
  );

  const login = async (credentials: SignInCredentials) => {
    try {
      await dispatch(signIn(credentials)).unwrap();
      navigate('/');
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  const register = async (request: SignUpRequest) => {
    try {
      await dispatch(signUp(request)).unwrap();
      navigate('/');
    } catch (error) {
      console.error('Registration failed:', error);
    }
  };

  const logout = () => {
    dispatch(signOut());
    navigate('/signin');
  };

  return {
    isAuthenticated,
    token,
    user,
    loading,
    error,
    login,
    register,
    logout,
  };
}; 