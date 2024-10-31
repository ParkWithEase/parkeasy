export const BACKEND_SERVER = import.meta.env.VITE_API_SERVER_URL || 'http://localhost:8080';
export const DAY_IN_A_WEEK = 7;
export const TOTAL_SEGMENTS_NUMBER = 48;
export const ERROR_MESSAGE_TIME_OUT = 3_000;
export const WAIT_TIME_BEFORE_AUTO_COMPLETE = 1_000;

//Error Message
export const PASSWORD_NOT_MATCH = 'password and confirm password not match';
export const INTERNAL_SERVER_ERROR = 'Something went wrong with the server';
export const CREATE_WITH_EMPTY_AVAILABILITY_TABLE_ERROR =
    "Can't create a parking spot with an empty availability table";
