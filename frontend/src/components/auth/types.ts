export type FormType = "login" | "register";

export interface LoginFormData {
	msisdn: string;
	password: string;
}

export interface RegisterFormData {
	name: string;
	msisdn: string;
	password: string;
	confirmPassword: string;
}

export interface AuthFormProps {
	isLoading: boolean;
	error: string;
	onError: (error: string) => void;
	onLogin: (data: LoginFormData) => Promise<void>;
	onRegister: (data: RegisterFormData) => Promise<void>;
}
