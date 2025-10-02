// MisViáticos Register - Type Definitions

export interface RegisterFormData {
  full_name: string;
  email: string;
  phone: string;
  password: string;
  password_confirm: string;
}

export interface RegisterFormProps {
  onSubmit: (data: RegisterFormData) => Promise<void>;
  isLoading: boolean;
}

export interface RegisterResponse {
  success: boolean;
  token?: string;
  refreshToken?: string;
  expiresIn?: number;
  user?: {
    id: string;
    email: string;
    full_name: string;
    phone: string;
    email_token?: string;
    requires_email_verification?: boolean;
  };
  error?: string;
}

export interface RegisterError {
  code: string;
  message: string;
  field?: string;
}

export interface PasswordStrength {
  score: number;
  feedback: string[];
  isValid: boolean;
}
