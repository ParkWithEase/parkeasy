# ParkEasy Manual Tests
These manual tests should be run before every release of the Calculator app.

## ParkEasy Smoke Tests

### Android

### Log In, Sign Up, Log Out

**Test 1**
Steps:
1. From the Login Screen, press the "Register Now" link
2. Select the "Name" text field, input "John Appleseed"
3. Select the "Email" text field, input "user@example.com"
4. Select the "Password" text field, input "password"
5. Select the "Confirm Password" text field, input "password"
6. Press the "Register" button
*Expected: "Registered successfully" is displayed and app navigates away from the Login Screen*

**Test 2**
1. From the Profile Screen, press the "Logout" button
*Expected: App navigates to the Login Screen*

**Test 3**
1. From the Login Screen, select the "email" text field, input "user@example.com"
2. Select the "password" text field, input "password"
3. Press the "Login" button
*Expected: "Logged in successfully" is displayed and app navigates away from the Login Screen*

### Webapp

### Log In, Sign Up, Log Out

**Test 1**
1. Open Firefox, navigate to "http://localhost:5173/auth/login"
2. Press the "Create one" link
3. Select the "First Name" text field, input "Robert"
4. Select the "Last Name" text field, input "Guderian"
5. Select the "Email" text field, input "robg@cs.umanitoba.ca"
6. Select the "Password" text field, input "password"
7. Select the "Confirm Password" text field, input "password"
8. Press the "Sign Up" button
*Expected: "Account created successfully" is displayed*

**Test 2**
1. Open Firefox, navigate to "http://localhost:5173/auth/login"
2. Select the "Email" text field, input "robg@cs.umanitoba.ca"
3. Select the "Password" text field, input "password"
4. Press the "Login" button
*Expected: Website redirects to a different page*
5. Press the "my profile" button
*Expected: "Robert Guderian" and "robg@cs.umanitoba.ca" are displayed*
6. Press the "Logout" button
*Expected: Website redirects to a the login page"
