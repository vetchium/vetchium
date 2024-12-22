import { useSelector, useDispatch } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { RootState } from '@/store';
import { signIn, signOut, setTFAToken } from '@/store/slices/authSlice';
import { SignInCredentials, TFARequest } from '@/types/auth';

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

  const submitTFA = async (tfaRequest: TFARequest) => {
    try {
      await dispatch(setTFAToken(tfaRequest)).unwrap();
      navigate('/');
    } catch (error) {
      console.error('TFA verification failed:', error);
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
    logout,
    submitTFA,
  };
}; 