// MisViáticos Login - Type Definitions

export interface LoginFormData {
  email: string;
  password: string;
}

export interface LoginFormProps {
  onSubmit: (data: LoginFormData) => Promise<void>;
  isLoading: boolean;
}

export interface LoginResponse {
  success: boolean;
  token?: string;
  refreshToken?: string;
  expiresIn?: number;
  user?: {
    id: string;
    email: string;
    username: string;
    first_name: string;
    last_name: string;
    full_name: string;
    phone: string;
    email_verified: boolean;
    is_active: boolean;
    last_login: string;
  };
  error?: string;
}

export interface AuthError {
  code: string;
  message: string;
  field?: string;
}

export interface SocialLoginProvider {
  name: string;
  id: 'google' | 'microsoft';
  icon: React.ComponentType<{ className?: string }>;
  callback: () => void;
}
