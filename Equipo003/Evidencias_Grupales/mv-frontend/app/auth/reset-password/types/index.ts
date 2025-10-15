// MisViáticos Reset Password - Type Definitions

export interface ResetPasswordFormData {
  email: string;
  token: string;
  new_password: string;
  confirm_password: string;
}

export interface ResetPasswordFormProps {
  onSubmit: (data: ResetPasswordFormData) => Promise<void>;
  isLoading: boolean;
  mode?: 'request' | 'reset';
}

export interface ResetPasswordResponse {
  success: boolean;
  message?: string;
  error?: string;
}

export interface ResetPasswordError {
  code: string;
  message: string;
  field?: string;
}

export interface TokenValidationResponse {
  isValid: boolean;
  email?: string;
  expiresAt?: string;
  error?: string;
}
