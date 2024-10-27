export const BACKEND_SERVER = import.meta.env.VITE_API_SERVER_URL || 'http://localhost:8080';

//Error Message
export const PASSWORD_NOT_MATCH = 'password and confirm password not match';
export const INTERNAL_SERVER_ERROR = 'Something went wrong with the server';
export const CREATE_WITH_EMPTY_AVAILABILITY_TABLE_ERROR = "Can't create a parking spot with an empty availability table";
