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
  user?: {
    id: string;
    email: string;
    firstname: string;
    lastname: string;
    phone: string;
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
