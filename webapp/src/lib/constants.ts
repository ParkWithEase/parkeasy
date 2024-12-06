export const BACKEND_SERVER = import.meta.env.VITE_API_SERVER_URL || 'http://localhost:8080';
export const DAY_IN_A_WEEK = 7;
export const TOTAL_SEGMENTS_NUMBER = 48;
export const ERROR_MESSAGE_TIME_OUT = 3_000;
export const WAIT_TIME_BEFORE_AUTO_COMPLETE = 1_000;
export const RESPONSE_WIDTH = 600;
export const SPOT_PREFFERED_ICON_SIZE = 32;

//Error Message
export const PASSWORD_NOT_MATCH = 'password and confirm password not match';
export const INTERNAL_SERVER_ERROR = 'Something went wrong with the server';
export const CREATE_WITH_EMPTY_AVAILABILITY_TABLE_ERROR =
    "Can't create a parking spot with an empty availability table";
export const ACCOUNT_CREATE_SUCCESS =
    'Account created successfully. Normally we would ask for email verification but for demo... nah';
export const DEFAULT_ACCOUNT_CREATION_ERROR = 'Wrong invalid email or invalid password';
export const PASSWORD_RESET_TOKEN_GET_ERROR =
    "We shouldn't be doing this but for demo sake, your email doesn't exist";
export const BOOK_WITHOUT_SLOTS_ERROR = 'Cannot book without any slots selected';
export const BOOK_WITHOUT_CAR = 'Cannot book without any car selected';

export const PASSWORD_RESET_SUCCESS_MESSAGE = 'Password changed successfully';

export const LATITUDE = 49.88887;
export const LONGITUDE = -97.13449;
export const DISTANCE = 5000;

//Constants for app/+page.svelte
export const INIT_ZOOM = 3;
export const MAX_ZOOM = 12;
export const SELECTED_ZOOM = 11;
export const DEFAULT_DISTANCE = 2000;
export const MIN_DISTANCE_RADIUS = 100;
export const MAX_DISTANCE_RADIUS = 5000;
export const DISTANCE_RADIUS_STEP = 100;
